package main

import (
	"context"
	"fmt"
	"github.com/toddyco/bacnet2go/bac_ip/client"
	"github.com/toddyco/bacnet2go/bac_ip/services"
	"github.com/toddyco/bacnet2go/bac_specs"
	"net"
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

	c, err := client.NewClientByIP("10.1.1.147", client.DefaultUDPPort)

	if err != nil {
		fmt.Println(err)
		return
	}

	//err = c.IAm()
	//
	//if err != nil {
	//	fmt.Println(err)
	//}

	// ctx, cancel := context.WithTimeout(context.Background(), 800*time.Second)
	// defer cancel()
	//
	// data, err := c.ReadProperty()

	ip := net.ParseIP("10.1.1.64")[12:16]

	//a, _ := hex.DecodeString("00d0db000300")

	addr := bac_specs.Address{
		Mac: []byte{ip[0], ip[1], ip[2], ip[3], 0xba, 0xc0},
		Net: 0,
		Adr: []byte{},
	}

	for range make([]int, 50) {
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogInput,
		//	Instance: bacnet.ObjectInstance(1),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(175),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(176),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(177),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(178),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(179),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(180),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(181),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(182),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(183),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(184),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(185),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(186),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(187),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(188),
		//})
		//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
		//	Type:     bacnet.AnalogValue,
		//	Instance: bacnet.ObjectInstance(189),
		//})
	}

	wg := sync.WaitGroup{}

	//GetPointList(c, addr, 700900)

	//GetPresentValue(c, addr, 700900, bacnet.ObjectID{
	//	Type:     bacnet.AnalogInput,
	//	Instance: bacnet.ObjectInstance(1),
	//})

	for x := range make([]int, 1) {
		wg.Add(1)

		go func(slp int) {
			time.Sleep(time.Duration(slp*20) * time.Millisecond)

			//c, err := client.NewClient("en0", client.DefaultUDPPort)
			//
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}

			pts := []bac_specs.ObjectID{
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogInput,
					Instance: bac_specs.ObjectInstance(1),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.MultiStateValue,
					Instance: bac_specs.ObjectInstance(76),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(176),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(177),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(178),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(179),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(180),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(181),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(182),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogInput,
					Instance: bac_specs.ObjectInstance(183),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(184),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(185),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(186),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(187),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(188),
				},
				bac_specs.ObjectID{
					Type:     bac_specs.AnalogValue,
					Instance: bac_specs.ObjectInstance(189),
				},
			}

			GetPointDetails(c, addr, 700900, pts)

			c.Close()

			wg.Done()
		}(x)
	}

	wg.Wait()

	// devices, err := c.WhoIs(bacip.WhoIs{
	//    Low:  nil,
	//    High: nil,
	//}, time.Second*10)
	//
	// j, err = json.Marshal(devices)
	//
	// fmt.Println(string(j))

	fmt.Println(time.Now().UnixMilli() - t)

}

func GetPresentValue(c *client.Client, addr bac_specs.Address, instanceID int, objectID bac_specs.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	// Addr := bacnet.AddressFromUDP(net.UDPAddr{
	//    IP:   net.ParseIP(ipAddr),
	//    Port: bacip.DefaultUDPPort,
	//})

	val, err := c.ReadProperty(ctx, makeDevice(addr, instanceID), services.ReadProperty{
		ObjectID: objectID,
		PropertyID: bac_specs.PropertyIdentifier{
			Type: bac_specs.PresentValue,
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)
}

func GetPointList(c *client.Client, addr bac_specs.Address, instanceID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	// Addr := bacnet.AddressFromUDP(net.UDPAddr{
	//    IP:   net.ParseIP(ipAddr),
	//    Port: bacip.DefaultUDPPort,
	//})

	val, err := c.ReadProperty(ctx, makeDevice(addr, instanceID), services.ReadProperty{
		ObjectID: bac_specs.ObjectID{
			Type:     bac_specs.BacnetDevice,
			Instance: bac_specs.ObjectInstance(instanceID),
		},
		PropertyID: bac_specs.PropertyIdentifier{
			Type: bac_specs.ObjectList,
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)

	time.Sleep(5)
}

func GetPointDetails(c *client.Client, addr bac_specs.Address, instanceID int, objectIDs []bac_specs.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	val, err := c.ReadPropertyMultiple(ctx, makeDevice(addr, instanceID), services.ReadPropertyMultiple{
		ObjectIDs: objectIDs,
		PropertyIDs: [][]bac_specs.PropertyIdentifier{{
			bac_specs.PropertyIdentifier{
				Type: bac_specs.ObjectName,
			},
			bac_specs.PropertyIdentifier{
				Type: bac_specs.PresentValue,
			},
			//bacnet.PropertyIdentifier{
			//	Type: bacnet.Description,
			//},
			//bacnet.PropertyIdentifier{
			//	Type: bacnet.Units,
			//},
		}},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)

	time.Sleep(5)
}

func makeDevice(addr bac_specs.Address, instanceID int) bac_specs.Device {
	return bac_specs.Device{
		ID: bac_specs.ObjectID{
			Type:     bac_specs.BacnetDevice,
			Instance: bac_specs.ObjectInstance(instanceID),
		},
		MaxApdu:      0,
		Segmentation: 0,
		Vendor:       0,
		Addr:         addr,
	}
}
