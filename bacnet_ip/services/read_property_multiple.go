package services

import (
	"errors"
	"github.com/toddyco/bacnet2go/bacnet"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type ReadPropertyMultiple struct {
	ObjectIDs   []bacnet.ObjectID
	PropertyIDs [][]bacnet.PropertyIdentifier
	Data        []interface{} // will contain the response
}

func (rpm ReadPropertyMultiple) MarshalBinary() ([]byte, error) {
	rpm.preparePropertyIDsSlice()

	encoder := encoding.NewEncoder()

	for _, objectID := range rpm.ObjectIDs {
		encoder.ContextObjectID(0, objectID)
		encoder.ContextOpening(1)

		for _, propertyID := range rpm.PropertyIDs {
			encoder.ContextUnsigned(0, uint32(propertyID[0].Type))
		}

		encoder.ContextClosing(1)
	}

	return encoder.Bytes(), encoder.Error()
}

func (rpm *ReadPropertyMultiple) UnmarshalBinary(data []byte) error {
	rpm.Data = make([]interface{}, 0, len(rpm.PropertyIDs))

	for objectIdx, objectID := range rpm.ObjectIDs {

		decoder := encoding.NewDecoder(data)
		decoder.ContextObjectID(0, &objectID)
		decoder.ContextOpening(1)

		tag, err := decoder.PeekTag()

		if tag.Closing {
			break
		}

		var val uint32
		decoder.ContextValue(2, &val)

		if decoder.Error() != nil {
			break
		}

		rpm.PropertyIDs[objectIdx] = append(rpm.PropertyIDs[objectIdx], bacnet.PropertyIdentifier{})
		rpm.Data = append(rpm.Data, nil)

		rpm.PropertyIDs[objectIdx][0].Type = bacnet.PropertyType(val)
		rpm.PropertyIDs[objectIdx][0].ArrayIndex = new(uint32)
		decoder.ContextValue(3, rpm.PropertyIDs[objectIdx][0].ArrayIndex)
		err = decoder.Error()
		var e encoding.ErrorIncorrectTagID

		// The array index tag is optional, per BACnet spec
		if err != nil && errors.As(err, &e) {
			rpm.PropertyIDs[objectIdx][0].ArrayIndex = nil
			decoder.ResetError()
		}

		decoder.ContextAbstractType(4, &rpm.Data[objectIdx])

		err = decoder.Error()

		if err != nil {
			return err
		}
	}

	return nil
}

// preparePropertyIDsSlice ensures that rpm.PropertyIDs has one slice of PropertyIdentifier
// structs for every rpm.ObjectIDs
func (rpm *ReadPropertyMultiple) preparePropertyIDsSlice() {
	if len(rpm.ObjectIDs) > len(rpm.PropertyIDs) {
		for i := range rpm.ObjectIDs {
			if i == 0 {
				continue
			}

			rpm.PropertyIDs[i] = rpm.PropertyIDs[0]
		}
	}
}
