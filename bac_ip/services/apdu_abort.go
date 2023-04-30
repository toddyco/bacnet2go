package services

import (
	"fmt"
	"github.com/toddyco/bacnet2go/bac_specs"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type APDUAbort struct {
	Reason bac_specs.AbortReason
}

func (a APDUAbort) Error() string {
	return fmt.Sprintf("APDU abort reason %v", a.Reason)
}

func (a APDUAbort) MarshalBinary() ([]byte, error) {
	panic("not implemented")
}

func (a *APDUAbort) UnmarshalBinary(data []byte) error {
	decoder := encoding.NewDecoder(data)
	decoder.AppData(&a.Reason, nil)
	return decoder.Error()
}
