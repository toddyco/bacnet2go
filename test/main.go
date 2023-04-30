package main

import (
	"context"
	"fmt"
	"github.com/toddyco/bacnet2go/client"
	services2 "github.com/toddyco/bacnet2go/services"
	"github.com/toddyco/bacnet2go/specs"
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

	addr := specs.Address{
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

			pts := []specs.ObjectID{
				specs.ObjectID{
					Type:     specs.AnalogInput,
					Instance: specs.ObjectInstance(1),
				},
				specs.ObjectID{
					Type:     specs.MultiStateValue,
					Instance: specs.ObjectInstance(76),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(176),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(177),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(178),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(179),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(180),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(181),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(182),
				},
				specs.ObjectID{
					Type:     specs.AnalogInput,
					Instance: specs.ObjectInstance(183),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(184),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(185),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(186),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(187),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(188),
				},
				specs.ObjectID{
					Type:     specs.AnalogValue,
					Instance: specs.ObjectInstance(189),
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

func GetPresentValue(c *client.Client, addr specs.Address, instanceID int, objectID specs.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	// Addr := bacnet.AddressFromUDP(net.UDPAddr{
	//    IP:   net.ParseIP(ipAddr),
	//    Port: bacip.DefaultUDPPort,
	//})

	val, err := c.ReadProperty(ctx, makeDevice(addr, instanceID), services2.ReadProperty{
		ObjectID: objectID,
		PropertyID: specs.PropertyIdentifier{
			Type: specs.PresentValue,
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)
}

func GetPointList(c *client.Client, addr specs.Address, instanceID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	// Addr := bacnet.AddressFromUDP(net.UDPAddr{
	//    IP:   net.ParseIP(ipAddr),
	//    Port: bacip.DefaultUDPPort,
	//})

	val, err := c.ReadProperty(ctx, makeDevice(addr, instanceID), services2.ReadProperty{
		ObjectID: specs.ObjectID{
			Type:     specs.BacnetDevice,
			Instance: specs.ObjectInstance(instanceID),
		},
		PropertyID: specs.PropertyIdentifier{
			Type: specs.ObjectList,
		},
		Data: nil,
	})

	fmt.Printf("%v \n", val)
	fmt.Printf("%v \n", err)

	time.Sleep(5)
}

func GetPointDetails(c *client.Client, addr specs.Address, instanceID int, objectIDs []specs.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
	defer cancel()

	val, err := c.ReadPropertyMultiple(ctx, makeDevice(addr, instanceID), services2.ReadPropertyMultiple{
		ObjectIDs: objectIDs,
		PropertyIDs: [][]specs.PropertyIdentifier{{
			specs.PropertyIdentifier{
				Type: specs.ObjectName,
			},
			specs.PropertyIdentifier{
				Type: specs.PresentValue,
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

func makeDevice(addr specs.Address, instanceID int) specs.Device {
	return specs.Device{
		ID: specs.ObjectID{
			Type:     specs.BacnetDevice,
			Instance: specs.ObjectInstance(instanceID),
		},
		MaxApdu:      0,
		Segmentation: 0,
		Vendor:       0,
		Addr:         addr,
	}
}
