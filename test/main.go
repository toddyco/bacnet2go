package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/toddyco/bacnet2go/bacnet"
	"github.com/toddyco/bacnet2go/bacnet_ip"
	"github.com/toddyco/bacnet2go/bacnet_ip/services"
	"time"
)

func main() {
	client, err := bacnet_ip.NewClient("en7", bacnet_ip.DefaultUDPPort)

	if err != nil {
		fmt.Println(err)
		return
	}

	j, _ := json.Marshal(client)

	fmt.Println(string(j))

	err = client.IAm()

	if err != nil {
		fmt.Println(err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 800*time.Second)
	// defer cancel()
	//
	// data, err := client.ReadProperty()

	addr := bacnet.Address{
		Mac: []byte{10, 1, 1, 64, 186, 192},
		Net: 0,
		Adr: nil,
	}

	//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
	//	Type:     bacnet.AnalogInput,
	//	Instance: bacnet.ObjectInstance(1),
	//})
	//
	//GetPointList(client, addr, 700900)

	GetPointDetails(client, addr, 700900, bacnet.ObjectID{
		Type:     bacnet.AnalogInput,
		Instance: bacnet.ObjectInstance(1),
	})

	// devices, err := client.WhoIs(bacip.WhoIs{
	//    Low:  nil,
	//    High: nil,
	//}, time.Second*10)
	//
	// j, err = json.Marshal(devices)
	//
	// fmt.Println(string(j))

}

func GetPresentValue(c *bacnet_ip.Client, addr bacnet.Address, instanceID int, objectID bacnet.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	// Addr := bacnet.AddressFromUDP(net.UDPAddr{
	//    IP:   net.ParseIP(ipAddr),
	//    Port: bacip.DefaultUDPPort,
	//})

	val, err := c.ReadProperty(ctx, makeDevice(addr, instanceID), services.ReadProperty{
		ObjectID: objectID,
		PropertyID: bacnet.PropertyIdentifier{
			Type: bacnet.PresentValue,
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)
}

func GetPointList(c *bacnet_ip.Client, addr bacnet.Address, instanceID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	// Addr := bacnet.AddressFromUDP(net.UDPAddr{
	//    IP:   net.ParseIP(ipAddr),
	//    Port: bacip.DefaultUDPPort,
	//})

	val, err := c.ReadProperty(ctx, makeDevice(addr, instanceID), services.ReadProperty{
		ObjectID: bacnet.ObjectID{
			Type:     bacnet.BacnetDevice,
			Instance: bacnet.ObjectInstance(instanceID),
		},
		PropertyID: bacnet.PropertyIdentifier{
			Type: bacnet.ObjectList,
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)

	time.Sleep(5)
}

func GetPointDetails(c *bacnet_ip.Client, addr bacnet.Address, instanceID int, objectID bacnet.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	val, err := c.ReadPropertyMultiple(ctx, makeDevice(addr, instanceID), services.ReadPropertyMultiple{
		ObjectID: bacnet.ObjectID{
			Type:     bacnet.BacnetDevice,
			Instance: bacnet.ObjectInstance(instanceID),
		},
		PropertyIDs: []bacnet.PropertyIdentifier{
			bacnet.PropertyIdentifier{
				Type: bacnet.ObjectName,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.PresentValue,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.StatusFlags,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.Description,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.Units,
			},
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)

	time.Sleep(5)
}

func makeDevice(addr bacnet.Address, instanceID int) bacnet.Device {
	return bacnet.Device{
		ID: bacnet.ObjectID{
			Type:     bacnet.BacnetDevice,
			Instance: bacnet.ObjectInstance(instanceID),
		},
		MaxApdu:      0,
		Segmentation: 0,
		Vendor:       0,
		Addr:         addr,
	}
}
