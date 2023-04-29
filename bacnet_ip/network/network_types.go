package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/toddyco/bacnet2go/bacnet_ip/client"
	"github.com/toddyco/bacnet2go/bacnet_ip/services"

	"github.com/toddyco/bacnet2go/bacnet"
)

type Version byte

const Version1 Version = 1

//go:generate stringer -type=NPDUPriority
type NPDUPriority byte

const (
	LifeSafety        NPDUPriority = 3
	CriticalEquipment NPDUPriority = 2
	Urgent            NPDUPriority = 1
	Normal            NPDUPriority = 0
)

type NPDU struct {
	Version Version // Always one
	// This 3 fields are packed in the control byte
	IsNetworkLayerMessage bool // If true, there is no APDU
	ExpectingReply        bool
	Priority              NPDUPriority

	Destination *bacnet.Address
	Source      *bacnet.Address
	HopCount    byte
	// The two are only significant if IsNetworkLayerMessage is true
	NetworkMessageType byte
	VendorID           uint16

	APDU *APDU
}

func (npdu NPDU) MarshalBinary() ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteByte(byte(npdu.Version))
	var control byte
	var hasSrc, hasDest, isNetworkMessage bool
	if npdu.IsNetworkLayerMessage {
		control += 1 << 7
		isNetworkMessage = true
	}
	if npdu.ExpectingReply {
		control += 1 << 2
	}
	if npdu.Priority > 3 {
		return nil, fmt.Errorf("invalid Priority %d", npdu.Priority)
	}
	control += byte(npdu.Priority)
	if npdu.Destination != nil && npdu.Destination.Net != 0 {
		control += 1 << 5
		hasDest = true
	}
	if npdu.Source != nil && npdu.Source.Net != 0 {
		control += 1 << 3
		hasSrc = true
	}
	b.WriteByte(control)
	if hasDest {
		_ = binary.Write(b, binary.BigEndian, npdu.Destination.Net)
		_ = binary.Write(b, binary.BigEndian, byte(len(npdu.Destination.Adr)))
		_ = binary.Write(b, binary.BigEndian, npdu.Destination.Adr)
	}
	if hasSrc {
		_ = binary.Write(b, binary.BigEndian, npdu.Source.Net)
		_ = binary.Write(b, binary.BigEndian, byte(len(npdu.Source.Adr)))
		_ = binary.Write(b, binary.BigEndian, npdu.Source.Adr)
	}
	if hasDest {
		b.WriteByte(npdu.HopCount)
	}
	if isNetworkMessage {
		b.WriteByte(npdu.NetworkMessageType)
		if npdu.NetworkMessageType >= 0x80 {
			_ = binary.Write(b, binary.BigEndian, npdu.VendorID)
		}
	}
	bytes := b.Bytes()
	if npdu.APDU != nil {
		bytesapdu, err := npdu.APDU.MarshalBinary()
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, bytesapdu...)
	}
	return bytes, nil
}

func (npdu *NPDU) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &npdu.Version)
	if err != nil {
		return fmt.Errorf("read NPDU version: %w", err)
	}
	if npdu.Version != Version1 {
		return fmt.Errorf("invalid NPDU version %d", npdu.Version)
	}
	control, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("read NPDU control byte:  %w", err)
	}
	if control&(1<<7) > 0 {
		npdu.IsNetworkLayerMessage = true
	}
	if control&(1<<2) > 0 {
		npdu.ExpectingReply = true
	}
	npdu.Priority = NPDUPriority(control & 0x3)

	if control&(1<<5) > 0 {
		npdu.Destination = &bacnet.Address{}
		err := binary.Read(buf, binary.BigEndian, &npdu.Destination.Net)
		if err != nil {
			return fmt.Errorf("read NPDU dest Address.Net: %w", err)
		}
		var length byte
		err = binary.Read(buf, binary.BigEndian, &length)
		if err != nil {
			return fmt.Errorf("read NPDU dest Address.Len: %w", err)
		}
		npdu.Destination.Adr = make([]byte, int(length))
		err = binary.Read(buf, binary.BigEndian, &npdu.Destination.Adr)
		if err != nil {
			return fmt.Errorf("read NPDU dest Address.Net: %w", err)
		}
	}

	if control&(1<<3) > 0 {
		npdu.Source = &bacnet.Address{}
		err := binary.Read(buf, binary.BigEndian, &npdu.Source.Net)
		if err != nil {
			return fmt.Errorf("read NPDU src Address.Net: %w", err)
		}
		var length byte
		err = binary.Read(buf, binary.BigEndian, &length)
		if err != nil {
			return fmt.Errorf("read NPDU src Address.Len: %w", err)
		}
		npdu.Source.Adr = make([]byte, int(length))
		err = binary.Read(buf, binary.BigEndian, &npdu.Source.Adr)
		if err != nil {
			return fmt.Errorf("read NPDU src Address.Net: %w", err)
		}
	}

	if npdu.Destination != nil {
		err := binary.Read(buf, binary.BigEndian, &npdu.HopCount)
		if err != nil {
			return fmt.Errorf("read NPDU HopCount: %w", err)
		}
	}

	if npdu.IsNetworkLayerMessage {
		err := binary.Read(buf, binary.BigEndian, &npdu.NetworkMessageType)
		if err != nil {
			return fmt.Errorf("read NPDU NetworkMessageType: %w", err)
		}
		if npdu.NetworkMessageType > 0x80 {
			err := binary.Read(buf, binary.BigEndian, &npdu.VendorID)
			if err != nil {
				return fmt.Errorf("read NPDU VendorId: %w", err)
			}
		}
	} else {
		npdu.APDU = &APDU{} // wipes out APDU from original request!
		return npdu.APDU.UnmarshalBinary(buf.Bytes())
	}

	return nil
}

// Todo: support more complex APDU
type APDU struct {
	DataType    PDUType
	ServiceType ServiceType
	Payload     Payload
	//Only meaningfully for confirmed  and ack
	InvokeID byte
	// MaxSegs
	// Segmented message
	// MoreFollow
	// SegmentedResponseAccepted
	// MaxApdu int
	// Sequence                  uint8
	// WindowNumber              uint8
}

func (apdu APDU) MarshalBinary() ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteByte(byte(apdu.DataType))
	if apdu.DataType.IsType(ConfirmedServiceRequest) {
		b.WriteByte(5) //Todo: Write other  control flag here
		b.WriteByte(apdu.InvokeID)
	}

	b.WriteByte(byte(apdu.ServiceType))
	bytes, err := apdu.Payload.MarshalBinary()
	if err != nil {
		return nil, err
	}
	b.Write(bytes)
	return b.Bytes(), nil
}

func (apdu *APDU) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &apdu.DataType)

	if err != nil {
		return fmt.Errorf("read APDU DataType: %w", err)
	}

	if apdu.DataType.SupportsInvokeID() {
		apdu.InvokeID, err = buf.ReadByte()

		if err != nil {
			return err
		}
	}

	// TODO: refactor
	err = binary.Read(buf, binary.BigEndian, &apdu.ServiceType)

	if err != nil {
		return fmt.Errorf("read APDU ServiceType: %w", err)
	}

	if apdu.DataType.IsType(UnconfirmedServiceRequest) && apdu.ServiceType == ServiceUnconfirmedWhoIs {
		apdu.Payload = &services.WhoIs{}
	} else if apdu.DataType.IsType(UnconfirmedServiceRequest) && apdu.ServiceType == ServiceUnconfirmedIAm {
		apdu.Payload = &services.IAm{}
	} else if apdu.DataType.IsType(ComplexAck) && apdu.ServiceType == ServiceConfirmedReadProperty {
		apdu.Payload = &services.ReadProperty{}
	} else if apdu.DataType.IsType(ComplexAck) && apdu.ServiceType == ServiceConfirmedReadPropertyMultiple {
		apdu.Payload = &services.ReadPropertyMultiple{}
	} else if apdu.DataType.IsType(Error) {
		apdu.Payload = &services.APDUError{}
	} else if apdu.DataType.IsType(Abort) {
		apdu.Payload = &services.APDUAbort{}
	} else {
		apdu.Payload = &DataPayload{} // Just pass raw data, decoding is not yet ready
	}

	return apdu.Payload.UnmarshalBinary(buf.Bytes())
}

type Payload interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
}

type DataPayload struct {
	Bytes []byte
}

func (p DataPayload) MarshalBinary() ([]byte, error) {
	return p.Bytes, nil
}

func (p *DataPayload) UnmarshalBinary(data []byte) error {
	p.Bytes = make([]byte, len(data))
	copy(p.Bytes, data)
	return nil
}

type BVLCType byte

const TypeBacnetIP BVLCType = 0x81

//go:generate stringer -type=Function
type Function byte

const (
	BacFuncResult                          Function = 0
	BacFuncWriteBroadcastDistributionTable Function = 1
	BacFuncBroadcastDistributionTable      Function = 2
	BacFuncBroadcastDistributionTableAck   Function = 3
	BacFuncForwardedNPDU                   Function = 4
	BacFuncUnicast                         Function = 10
	BacFuncBroadcast                       Function = 11
)

type BVLC struct {
	Type     BVLCType
	Function Function
	NPDU     NPDU
}

func (bvlc BVLC) MarshalBinary() ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteByte(byte(bvlc.Type))
	b.WriteByte(byte(bvlc.Function))
	data, err := bvlc.NPDU.MarshalBinary()
	if err != nil {
		return nil, err
	}
	len := uint16(4 + len(data)) //len includes Type,Function and itself
	_ = binary.Write(b, binary.BigEndian, len)
	b.Write(data)
	return b.Bytes(), nil
}

func (bvlc *BVLC) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	bvlcType, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("read bvlc type: %w", err)
	}
	bvlc.Type = BVLCType(bvlcType)
	if bvlc.Type != TypeBacnetIP {
		return client.ErrNotBACnetIP
	}
	bvlcFunc, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("read bvlc func: %w", err)
	}
	var length uint16
	err = binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return fmt.Errorf("read bvlc length: %w", err)
	}
	remaining := buf.Bytes()

	bvlc.Function = Function(bvlcFunc)
	if len(remaining) != int(length)-4 {
		return fmt.Errorf("incoherent Length field in BVCL. Advertized payload size is %d, real size  %d", length-4, len(remaining))
	}
	return bvlc.NPDU.UnmarshalBinary(remaining)
}
