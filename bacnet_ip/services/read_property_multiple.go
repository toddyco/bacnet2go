package services

import (
	"github.com/toddyco/bacnet2go/bacnet"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type ReadPropertyMultiple struct {
	ObjectID    bacnet.ObjectID
	PropertyIDs []bacnet.PropertyIdentifier
	Data        interface{} // will contain the response
}

func (rpm *ReadPropertyMultiple) MarshalBinary() ([]byte, error) {
	encoder := encoding.NewEncoder()
	encoder.ContextObjectID(0, rpm.ObjectID)
	encoder.ContextOpening(1)

	for _, propertyID := range rpm.PropertyIDs {
		encoder.ContextUnsigned(0, uint32(propertyID.Type))
	}

	encoder.ContextClosing(1)

	return encoder.Bytes(), encoder.Error()
}

func (rpm *ReadPropertyMultiple) UnmarshalBinary(data []byte) error {
	return nil
}
