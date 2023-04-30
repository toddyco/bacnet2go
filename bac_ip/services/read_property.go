package services

import (
	"errors"
	"github.com/toddyco/bacnet2go/bac_specs"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type ReadProperty struct {
	ObjectID   bac_specs.ObjectID
	PropertyID bac_specs.PropertyIdentifier
	Data       interface{} // will contain the response
}

func (rp ReadProperty) MarshalBinary() ([]byte, error) {
	encoder := encoding.NewEncoder()
	encoder.ContextObjectID(0, rp.ObjectID)
	encoder.ContextUnsigned(1, uint32(rp.PropertyID.Type))

	// The array index tag is optional, per BACnet spec
	if rp.PropertyID.ArrayIndex != nil {
		encoder.ContextUnsigned(2, *rp.PropertyID.ArrayIndex)
	}

	return encoder.Bytes(), encoder.Error()
}

func (rp *ReadProperty) UnmarshalBinary(data []byte) error {
	decoder := encoding.NewDecoder(data)
	decoder.ContextObjectID(0, &rp.ObjectID)
	var val uint32
	decoder.ContextValue(1, &val)
	rp.PropertyID.Type = bac_specs.PropertyType(val)
	rp.PropertyID.ArrayIndex = new(uint32)
	decoder.ContextValue(2, rp.PropertyID.ArrayIndex)
	err := decoder.Error()
	var e encoding.ErrorIncorrectTagID

	// The array index tag is optional, per BACnet spec
	if err != nil && errors.As(err, &e) {
		rp.PropertyID.ArrayIndex = nil
		decoder.ResetError()
	}

	decoder.ContextAbstractType(3, &rp.Data)
	return decoder.Error()
}
