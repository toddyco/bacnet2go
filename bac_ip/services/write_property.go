package services

import (
	"github.com/toddyco/bacnet2go/bac_specs"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type WriteProperty struct {
	ObjectID      bac_specs.ObjectID
	Property      bac_specs.PropertyIdentifier
	PropertyValue bac_specs.PropertyValue
	Priority      bac_specs.PriorityList
}

func (wp WriteProperty) MarshalBinary() ([]byte, error) {
	encoder := encoding.NewEncoder()
	encoder.ContextObjectID(0, wp.ObjectID)
	encoder.ContextUnsigned(1, uint32(wp.Property.Type))
	if wp.Property.ArrayIndex != nil {
		encoder.ContextUnsigned(2, uint32(*wp.Property.ArrayIndex))
	}
	err := encoder.ContextAbstractType(3, wp.PropertyValue)
	if err != nil {
		return nil, err
	}
	if wp.Priority != 0 {
		encoder.ContextUnsigned(4, uint32(wp.Priority))
	}
	return encoder.Bytes(), encoder.Error()
}

func (wp *WriteProperty) UnmarshalBinary(data []byte) error {
	decoder := encoding.NewDecoder(data)
	return decoder.Error()
}
