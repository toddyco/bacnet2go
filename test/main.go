package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/toddyco/bacnet2go/bacnet"
	"github.com/toddyco/bacnet2go/bacnet_ip/client"
	"github.com/toddyco/bacnet2go/bacnet_ip/services"
	"sync"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in handle message: ", r)
		}
	}()

	t := time.Now().UnixMilli()
	fmt.Println(t)

	client, err := client.NewClient("en0", client.DefaultUDPPort)

	if err != nil {
		fmt.Println(err)
		return
	}

	j, _ := json.Marshal(client)

	fmt.Println(string(j))

	//err = client.IAm()
	//
	//if err != nil {
	//	fmt.Println(err)
	//}

	// ctx, cancel := context.WithTimeout(context.Background(), 800*time.Second)
	// defer cancel()
	//
	// data, err := client.ReadProperty()

	addr := bacnet.Address{
		Mac: []byte{10, 1, 1, 64, 186, 192},
		Net: 0,
		Adr: nil,
	}

	for range make([]int, 1) {
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogInput,
		//	Instance: bacnet.ObjectInstance(1),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(175),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(176),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(177),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(178),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(179),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(180),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(181),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(182),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(183),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(184),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(185),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(186),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(187),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(188),
		//})
		//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(189),
		//})
	}

	wg := sync.WaitGroup{}

	//GetPointList(client, addr, 700900)

	//GetPresentValue(client, addr, 700900, bacnet.ObjectID{
	//	Type:     bacnet.AnalogInput,
	//	Instance: bacnet.ObjectInstance(1),
	//})

	for range make([]int, 1) {
		wg.Add(1)

		//go func() {
			pts := []bacnet.ObjectID{
				bacnet.ObjectID{
					Type:     bacnet.AnalogInput,
					Instance: bacnet.ObjectInstance(1),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(175),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(176),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(177),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(178),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(179),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(180),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(181),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(182),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogInput,
					Instance: bacnet.ObjectInstance(183),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(184),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(185),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(186),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(187),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(188),
				},
				bacnet.ObjectID{
					Type:     bacnet.AnalogValue,
					Instance: bacnet.ObjectInstance(189),
				},
			}

			GetPointDetails(client, addr, 700900, pts)

			wg.Done()
		//}()
	}

	wg.Wait()

	// devices, err := client.WhoIs(bacip.WhoIs{
	//    Low:  nil,
	//    High: nil,
	//}, time.Second*10)
	//
	// j, err = json.Marshal(devices)
	//
	// fmt.Println(string(j))

	fmt.Println(time.Now().UnixMilli() - t)

}

func GetPresentValue(c *client.Client, addr bacnet.Address, instanceID int, objectID bacnet.ObjectID) {
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

func GetPointList(c *client.Client, addr bacnet.Address, instanceID int) {
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

func GetPointDetails(c *client.Client, addr bacnet.Address, instanceID int, objectIDs []bacnet.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	val, err := c.ReadPropertyMultiple(ctx, makeDevice(addr, instanceID), services.ReadPropertyMultiple{
		ObjectIDs: objectIDs,
		PropertyIDs: [][]bacnet.PropertyIdentifier{{
			bacnet.PropertyIdentifier{
				Type: bacnet.ObjectName,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.PresentValue,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.Description,
			},
			bacnet.PropertyIdentifier{
				Type: bacnet.Units,
			},
		}},
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
