package services

import (
	"fmt"

	"github.com/toddyco/bacnet2go/bacnet"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type APDUError struct {
	Class bacnet.ErrorClass
	Code  bacnet.ErrorCode
}

func (e APDUError) Error() string {
	return fmt.Sprintf("apdu error class %v code %v", e.Class, e.Code)
}
func (e APDUError) MarshalBinary() ([]byte, error) {
	panic("not implemented")
}

func (e *APDUError) UnmarshalBinary(data []byte) error {
	decoder := encoding.NewDecoder(data)
	decoder.AppData(&e.Class, nil)
	decoder.AppData(&e.Code, nil)
	return decoder.Error()
}

// Todo http://kargs.net/BACnet/BACnet_Essential_Objects_Services.pdf -> Time synchro, Reinitialize device, DeviceCommunicationControl
