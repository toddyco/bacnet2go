// Package bacip implements a Bacnet/IP client
package client

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/toddyco/bacnet2go/bacnet_ip"
	"github.com/toddyco/bacnet2go/bacnet_ip/const"
	"github.com/toddyco/bacnet2go/bacnet_ip/network"
	"github.com/toddyco/bacnet2go/bacnet_ip/services"
	"net"
	"sync"
	"time"

	"github.com/toddyco/bacnet2go/specs"
)

type Client struct {
	// Maybe change to baetyl-bacnet address
	ipAddress        net.IP
	broadcastAddress net.IP
	udpPort          int
	udp              *net.UDPConn
	subscriptions    *Subscriptions
	transactions     *bacnet_ip.Transactions
	Logger           Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
}

type Subscriptions struct {
	sync.RWMutex
	f func(network.BVLC, net.UDPAddr)
}

const DefaultUDPPort = 47808

func broadcastAddr(n *net.IPNet) (net.IP, error) {
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}

// NewClient creates a new baetyl-bacnet client. It binds on the given port
// and network interface (eth0 for example). If Port if 0, the default
// baetyl-bacnet port is used
func NewClient(netInterface string, port int) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		subscriptions: &Subscriptions{},
		transactions:  bacnet_ip.NewTransactions(),
		Logger:        NoOpLogger{},
		ctx:           ctx,
		cancelFunc:    cancel,
	}

	i, err := net.InterfaceByName(netInterface)

	if err != nil {
		return nil, fmt.Errorf("interface %s: %w", netInterface, err)
	}

	if port == 0 {
		port = DefaultUDPPort
	}

	c.udpPort = port
	addrs, err := i.Addrs()

	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("interface %s has no addresses", netInterface)
	}

	for _, adr := range addrs {
		ip, ipnet, err := net.ParseCIDR(adr.String())
		if err != nil {
			return nil, err
		}
		// To4 is nil when type is ip6
		if ip.To4() != nil {
			broadcast, err := broadcastAddr(ipnet)
			if err != nil {
				return nil, err
			}
			c.ipAddress = ip.To4()
			c.broadcastAddress = broadcast
			break
		}
	}

	if c.ipAddress == nil {
		return nil, fmt.Errorf("no IPv4 address assigned to interface %s", netInterface)
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: c.udpPort,
	})

	if err != nil {
		return nil, err
	}

	c.udp = conn
	go c.listen()
	return c, nil
}

// NewClientByIP
func NewClientByIP(ip string, port int) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		subscriptions: &Subscriptions{},
		transactions:  bacnet_ip.NewTransactions(),
		Logger:        NoOpLogger{},
		ctx:           ctx,
		cancelFunc:    cancel,
	}

	if port == 0 {
		port = DefaultUDPPort
	}

	c.udpPort = port
	c.ipAddress = net.ParseIP(ip)
	addr, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}

	for _, ad := range addr {
		if ipNet, ok := ad.(*net.IPNet); ok {
			if ipNet.Contains(c.ipAddress) {
				broadcast, err := broadcastAddr(ipNet)
				if err != nil {
					return nil, err
				}
				c.broadcastAddress = broadcast
				break
			}
		}
	}

	if c.broadcastAddress == nil {
		return nil, errors.New("broadcast address not found")
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: c.udpPort,
	})

	if err != nil {
		return nil, err
	}

	c.udp = conn
	go c.listen()
	return c, nil
}

// listen for incoming baetyl-bacnet packets.
func (c *Client) listen() {
	var done bool

	for {
		b := make([]byte, 2048)
		i, addr, err := c.udp.ReadFromUDP(b)

		if err != nil {
			c.Logger.Error(err.Error())
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					c.Logger.Error("panic in handle message: ", r)
				}
			}()

			err := c.handleMessage(addr, b[:i])

			if err != nil {
				c.Logger.Error("handle msg: ", err)
			}
		}()

		select {
		case <-c.ctx.Done():
			done = true
		default:

		}

		if done {
			break
		}
	}
}

// Context
func (c *Client) Context() context.Context {
	return c.ctx
}

// Close
func (c *Client) Close() {
	c.cancelFunc()
	c.udp.Close()
}

// handleMessage
func (c *Client) handleMessage(src *net.UDPAddr, b []byte) error {
	var bvlc network.BVLC

	err := bvlc.UnmarshalBinary(b)

	if err != nil && errors.Is(err, _const.ErrNotBACnetIP) {
		return err
	}

	apdu := bvlc.NPDU.APDU

	if apdu == nil {
		c.Logger.Info(fmt.Sprintf("Received network packet %+v", bvlc.NPDU))
		return nil
	}

	c.subscriptions.RLock()

	if c.subscriptions.f != nil {
		// If f block, there is a deadlock here
		c.subscriptions.f(bvlc, *src)
	}

	c.subscriptions.RUnlock()

	if apdu.DataType.SupportsInvokeID() {
		invokeID := bvlc.NPDU.APDU.InvokeID
		tx, ok := c.transactions.GetTransaction(invokeID)

		if !ok {
			return fmt.Errorf("no transaction found for id %d", invokeID)
		}

		select {
		case tx.APDU <- *apdu:
			return nil
		case <-tx.Ctx.Done():
			return fmt.Errorf("handler for tx %d: %w", invokeID, tx.Ctx.Err())
		}
	}

	return nil
}

func (c *Client) IAm() error {
	npdu := network.NPDU{
		Version:               network.Version1,
		IsNetworkLayerMessage: false,
		ExpectingReply:        false,
		Priority:              network.Normal,
		Destination: &specs.Address{
			Net: uint16(0xffff),
		},
		Source: nil,
		APDU: &network.APDU{
			DataType:    network.UnconfirmedServiceRequest,
			ServiceType: network.ServiceUnconfirmedIAm,
			Payload: &services.IAm{
				ObjectID: specs.ObjectID{
					Type:     specs.BacnetDevice,
					Instance: 99999,
				},
				MaxApduLength:       0,
				SegmentationSupport: 0,
				VendorID:            0,
			},
		},
		HopCount: 255,
	}

	_, err := c.broadcast(npdu)

	return err
}

func (c *Client) WhoIs(data services.WhoIs, timeout time.Duration) ([]specs.Device, error) {
	npdu := network.NPDU{
		Version:               network.Version1,
		IsNetworkLayerMessage: false,
		ExpectingReply:        false,
		Priority:              network.Normal,
		Destination: &specs.Address{
			Net: uint16(0xffff),
		},
		Source: nil,
		APDU: &network.APDU{
			DataType:    network.UnconfirmedServiceRequest,
			ServiceType: network.ServiceUnconfirmedWhoIs,
			Payload:     &data,
		},
		HopCount: 255,
	}

	rChan := make(chan struct {
		bvlc network.BVLC
		src  net.UDPAddr
	})

	c.subscriptions.Lock()
	// TODO:  add errgroup ?, ensure all f are done and not blocked
	c.subscriptions.f = func(bvlc network.BVLC, src net.UDPAddr) {
		rChan <- struct {
			bvlc network.BVLC
			src  net.UDPAddr
		}{
			bvlc: bvlc,
			src:  src,
		}
	}
	c.subscriptions.Unlock()
	defer func() {
		c.subscriptions.f = nil
	}()
	_, err := c.broadcast(npdu)
	if err != nil {
		return nil, err
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	// Use a set to deduplicate results
	set := map[services.IAm]specs.Address{}

	for {
		select {
		case <-timer.C:
			result := []specs.Device{}
			for iam, addr := range set {
				result = append(result, specs.Device{
					ID:           iam.ObjectID,
					MaxApdu:      iam.MaxApduLength,
					Segmentation: iam.SegmentationSupport,
					Vendor:       iam.VendorID,
					Addr:         addr,
				})
			}
			return result, nil
		case r := <-rChan:
			// clean/filter  network answers here
			apdu := r.bvlc.NPDU.APDU
			if apdu != nil {
				if apdu.DataType.IsType(network.UnconfirmedServiceRequest) &&
					apdu.ServiceType == network.ServiceUnconfirmedIAm {
					iam, ok := apdu.Payload.(*services.IAm)
					if !ok {
						return nil, fmt.Errorf("unexpected payload type %T", apdu.Payload)
					}
					// Only add result that we are interested in. Well
					// behaved devices should not answer if their
					// InstanceID isn't in the given range. But because
					// the IAM response is in broadcast mode, we might
					// receive an answer triggered by an other whois
					if data.High != nil && data.Low != nil {
						if iam.ObjectID.Instance >= specs.ObjectInstance(*data.Low) &&
							iam.ObjectID.Instance <= specs.ObjectInstance(*data.High) {
							addr := specs.AddressFromUDP(r.src)
							if r.bvlc.NPDU.Source != nil {
								addr.Net = r.bvlc.NPDU.Source.Net
								addr.Adr = r.bvlc.NPDU.Source.Adr
							}
							set[*iam] = *addr
						}
					} else {
						addr := specs.AddressFromUDP(r.src)
						if r.bvlc.NPDU.Source != nil {
							addr.Net = r.bvlc.NPDU.Source.Net
							addr.Adr = r.bvlc.NPDU.Source.Adr
						}
						set[*iam] = *addr
					}

				}
			}
		}
	}
}

// ReadProperty reads a single property from an object
func (c *Client) ReadProperty(ctx context.Context, device specs.Device, readProp services.ReadProperty) (interface{}, error) {
	invokeID := c.transactions.GetID()
	defer c.transactions.FreeID(invokeID)

	npdu := network.NPDU{
		Version:               network.Version1,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              network.Normal,
		Destination:           &device.Addr,
		Source: specs.AddressFromUDP(net.UDPAddr{
			IP:   c.ipAddress,
			Port: c.udpPort,
		}),
		HopCount: 255,
		APDU: &network.APDU{
			DataType:    network.ConfirmedServiceRequest,
			ServiceType: network.ServiceConfirmedReadProperty,
			InvokeID:    invokeID,
			Payload:     &readProp,
		},
	}

	rChan := make(chan network.APDU)
	c.transactions.SetTransaction(invokeID, rChan, ctx)
	defer c.transactions.StopTransaction(invokeID)

	_, err := c.send(npdu)

	if err != nil {
		return nil, err
	}

	select {
	case apdu := <-rChan:
		// TODO: ensure response validity, ensure conversion cannot panic
		if apdu.DataType.IsType(network.Error) {
			return nil, *apdu.Payload.(*services.APDUError)
		}
		if apdu.DataType.IsType(network.ComplexAck) && apdu.ServiceType == network.ServiceConfirmedReadProperty {
			data := apdu.Payload.(*services.ReadProperty).Data
			return data, nil
		}
		return nil, errors.New("invalid answer")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ReadPropertyMultiple reads multiple properties from one or more objects
func (c *Client) ReadPropertyMultiple(ctx context.Context, device specs.Device, readProp services.ReadPropertyMultiple) (interface{}, error) {
	invokeID := c.transactions.GetID()
	defer c.transactions.FreeID(invokeID)

	npdu := network.NPDU{
		Version:               network.Version1,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              network.Normal,
		Destination:           &device.Addr,
		Source: specs.AddressFromUDP(net.UDPAddr{
			IP:   c.ipAddress,
			Port: c.udpPort,
		}),
		HopCount: 255,
		APDU: &network.APDU{
			DataType:    network.ConfirmedServiceRequest,
			ServiceType: network.ServiceConfirmedReadPropertyMultiple,
			InvokeID:    invokeID,
			Payload:     &readProp,
		},
	}

	rChan := make(chan network.APDU)
	c.transactions.SetTransaction(invokeID, rChan, ctx)
	defer c.transactions.StopTransaction(invokeID)

	_, err := c.send(npdu)

	if err != nil {
		return nil, err
	}

	select {
	case apdu := <-rChan:
		// TODO: ensure response validity, ensure conversion cannot panic
		if apdu.DataType.IsType(network.Error) {
			return nil, *apdu.Payload.(*services.APDUError)
		}

		if apdu.DataType.IsType(network.Abort) {
			if abort, ok := apdu.Payload.(*services.APDUAbort); ok {
				if abort.Reason == specs.SegmentationNotSupportedAbortReason {
					return nil, _const.ErrSegmentationNotSupported
				}
			}
		}

		if apdu.DataType.IsType(network.ComplexAck) && apdu.ServiceType == network.ServiceConfirmedReadPropertyMultiple {
			data := apdu.Payload.(*services.ReadPropertyMultiple).Data
			return data, nil
		}

		return nil, errors.New("invalid answer")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// WriteProperty writes a value to a property
func (c *Client) WriteProperty(ctx context.Context, device specs.Device, writeProp services.WriteProperty) error {
	invokeID := c.transactions.GetID()
	defer c.transactions.FreeID(invokeID)

	npdu := network.NPDU{
		Version:               network.Version1,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              network.Normal,
		Destination:           &device.Addr,
		Source: specs.AddressFromUDP(net.UDPAddr{
			IP:   c.ipAddress,
			Port: c.udpPort,
		}),
		HopCount: 255,
		APDU: &network.APDU{
			DataType:    network.ConfirmedServiceRequest,
			ServiceType: network.ServiceConfirmedWriteProperty,
			InvokeID:    invokeID,
			Payload:     &writeProp,
		},
	}

	wrChan := make(chan network.APDU)
	c.transactions.SetTransaction(invokeID, wrChan, ctx)
	defer c.transactions.StopTransaction(invokeID)

	_, err := c.send(npdu)

	if err != nil {
		return err
	}

	select {
	case apdu := <-wrChan:
		//Todo: ensure response validity, ensure conversion cannot panic
		if apdu.DataType.IsType(network.Error) {
			return *apdu.Payload.(*services.APDUError)
		}
		if apdu.DataType.IsType(network.SimpleAck) && apdu.ServiceType == network.ServiceConfirmedWriteProperty {
			return nil
		}
		return errors.New("invalid answer")
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (c *Client) send(npdu network.NPDU) (int, error) {
	bytes, err := network.BVLC{
		Type:     network.TypeBacnetIP,
		Function: network.BacFuncUnicast,
		NPDU:     npdu,
	}.MarshalBinary()

	if err != nil {
		return 0, err
	}

	if npdu.Destination == nil {
		return 0, fmt.Errorf("destination baetyl-bacnet address should be not nil to send unicast")
	}

	addr := specs.UDPFromAddress(*npdu.Destination)

	return c.udp.WriteToUDP(bytes, &addr)

}

func (c *Client) broadcast(npdu network.NPDU) (int, error) {
	bytes, err := network.BVLC{
		Type:     network.TypeBacnetIP,
		Function: network.BacFuncBroadcast,
		NPDU:     npdu,
	}.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return c.udp.WriteToUDP(bytes, &net.UDPAddr{
		IP:   c.broadcastAddress,
		Port: DefaultUDPPort,
	})
}
