package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/toddyco/bacnet2go/bacip"
    "github.com/toddyco/bacnet2go/bacnet"
    "time"
)

func main() {
    client, err := bacip.NewClient("en7", bacip.DefaultUDPPort)

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

    //ctx, cancel := context.WithTimeout(context.Background(), 800*time.Second)
    //defer cancel()
    //
    //data, err := client.ReadProperty()

    //getPropertyValue(client, "10.1.1.64", 700900)
    getObjects(client, "10.1.1.64", 700900)

    //devices, err := client.WhoIs(bacip.WhoIs{
    //    Low:  nil,
    //    High: nil,
    //}, time.Second*10)
    //
    //j, err = json.Marshal(devices)
    //
    //fmt.Println(string(j))

}

func getPropertyValue(c *bacip.Client, ipAddr string, instanceID int) {
    ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
    defer cancel()

    //Addr := bacnet.AddressFromUDP(net.UDPAddr{
    //    IP:   net.ParseIP(ipAddr),
    //    Port: bacip.DefaultUDPPort,
    //})

    Addr := bacnet.Address{
        Mac: []byte{10, 1, 1, 64, 186, 192},
        Net: 0,
        Adr: nil,
    }

    val, err := c.ReadProperty(ctx, bacnet.Device{
        ID: bacnet.ObjectID{
            Type:     bacnet.BacnetDevice,
            Instance: bacnet.ObjectInstance(instanceID),
        },
        MaxApdu:      0,
        Segmentation: 0,
        Vendor:       0,
        Addr:         Addr,
    }, bacip.ReadProperty{
        ObjectID: bacnet.ObjectID{
            Type:     bacnet.AnalogInput,
            Instance: bacnet.ObjectInstance(1),
        },
        Property: bacnet.PropertyIdentifier{
            Type: bacnet.PresentValue,
        },
        Data: nil,
    })

    fmt.Printf("%v \n", val)
    fmt.Printf("%v \n", err)
}

func getObjects(c *bacip.Client, ipAddr string, instanceID int) {
    ctx, cancel := context.WithTimeout(context.Background(), 800000*time.Second)
    defer cancel()

    //Addr := bacnet.AddressFromUDP(net.UDPAddr{
    //    IP:   net.ParseIP(ipAddr),
    //    Port: bacip.DefaultUDPPort,
    //})

    Addr := bacnet.Address{
        Mac: []byte{10, 1, 1, 64, 186, 192},
        Net: 0,
        Adr: nil,
    }

    val, err := c.ReadProperty(ctx, bacnet.Device{
        ID: bacnet.ObjectID{
            Type:     bacnet.BacnetDevice,
            Instance: bacnet.ObjectInstance(instanceID),
        },
        MaxApdu:      0,
        Segmentation: 0,
        Vendor:       0,
        Addr:         Addr,
    }, bacip.ReadProperty{
        ObjectID: bacnet.ObjectID{
            Type:     bacnet.BacnetDevice,
            Instance: bacnet.ObjectInstance(instanceID),
        },
        Property: bacnet.PropertyIdentifier{
            Type: bacnet.ObjectList,
        },
        Data: nil,
    })

    fmt.Printf("%v \n", val)
    fmt.Printf("%v \n", err)

    time.Sleep(5)
}
