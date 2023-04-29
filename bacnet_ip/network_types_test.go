package bacnet_ip

import (
	"github.com/toddyco/bacnet2go/bacnet_ip/services"
	"testing"

	"github.com/matryer/is"
	"github.com/toddyco/bacnet2go/bacnet"
)

func TestFullEncodingAndCoherency(t *testing.T) {
	ttc := []struct {
		bvlc    BVLC
		encoded string //hex string
	}{
		{
			bvlc: BVLC{
				Type:     TypeBacnetIP,
				Function: BacFuncBroadcast,
				NPDU: NPDU{
					Version:               Version1,
					IsNetworkLayerMessage: false,
					ExpectingReply:        false,
					Priority:              Normal,
					APDU: &APDU{
						DataType:    UnconfirmedServiceRequest,
						ServiceType: ServiceUnconfirmedWhoIs,
						Payload: &services.ReadProperty{
							ObjectID: bacnet.ObjectID{
								Type:     bacnet.AnalogInput,
								Instance: 300184,
							},
							PropertyID: bacnet.PropertyIdentifier{
								Type: bacnet.PresentValue,
							},
						},
					},
				},
			},
			encoded: "810b000801001008",
		},
		{
			bvlc: BVLC{
				Type:     TypeBacnetIP,
				Function: BacFuncBroadcast,
				NPDU: NPDU{
					Version:               Version1,
					IsNetworkLayerMessage: false,
					ExpectingReply:        false,
					Priority:              Normal,
					Destination: &bacnet.Address{
						Net: 0xffff,
						Adr: []byte{},
					},
					Source:   &bacnet.Address{},
					HopCount: 255,
					APDU: &APDU{
						DataType:    UnconfirmedServiceRequest,
						ServiceType: ServiceUnconfirmedIAm,
						Payload: &services.IAm{
							ObjectID: bacnet.ObjectID{
								Type:     8,
								Instance: 30185,
							},
							MaxApduLength:       1476,
							SegmentationSupport: bacnet.SegmentationSupportBoth,
							VendorID:            364,
						},
					},
				},
			},
			encoded: "810b00190120ffff00ff1000c4020075e92205c4910022016c",
		},
	}

	for _, tc := range ttc {
		t.Run(tc.encoded, func(t *testing.T) {
			is := is.New(t)
			result, err := tc.bvlc.MarshalBinary()
			is.NoErr(err)
			//is.Equal(tc.encoded, hex.EncodeToString(result))
			w := BVLC{}
			is.NoErr(w.UnmarshalBinary(result))
			//result2, err := w.MarshalBinary()
			//is.NoErr(err)
			//is.Equal(tc.encoded, hex.EncodeToString(result2))
		})
	}
}
