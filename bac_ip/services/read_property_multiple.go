package services

import (
	"errors"
	"github.com/toddyco/bacnet2go/bac_specs"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type ReadPropertyMultiple struct {
	ObjectIDs   []bac_specs.ObjectID
	PropertyIDs [][]bac_specs.PropertyIdentifier
	Data        [][]interface{} // will contain the response
}

type ObjectAndProperties struct {
	ObjectID    bac_specs.ObjectID
	PropertyIDs []bac_specs.PropertyIdentifier
	Data        []interface{}
}

func (rpm ReadPropertyMultiple) MarshalBinary() ([]byte, error) {
	//rpm.preparePropertyIDsSlice()

	encoder := encoding.NewEncoder()

	for _, objectID := range rpm.ObjectIDs {
		encoder.ContextObjectID(0, objectID)
		encoder.ContextOpening(1)

		for _, propertyID := range rpm.PropertyIDs[0] {
			encoder.ContextUnsigned(0, uint32(propertyID.Type))
		}

		encoder.ContextClosing(1)
	}

	return encoder.Bytes(), encoder.Error()
}

func (rpm *ReadPropertyMultiple) UnmarshalBinary(data []byte) error {
	rpm.ObjectIDs = []bac_specs.ObjectID{}
	rpm.PropertyIDs = [][]bac_specs.PropertyIdentifier{}
	rpm.Data = [][]interface{}{}

	var propIdx int

	decoder := encoding.NewDecoder(data)

	for { // Loop over objects
		objAndProps := ObjectAndProperties{
			ObjectID:    bac_specs.ObjectID{},
			PropertyIDs: []bac_specs.PropertyIdentifier{},
			Data:        []interface{}{},
		}

		decoder.ContextObjectID(0, &objAndProps.ObjectID)

		if decoder.Error() != nil {
			decoder.ResetError()
			break
		}

		decoder.ContextOpening(1)

		propIdx = 0

		for { // Loop over properties
			var val uint32
			decoder.ContextValue(2, &val)

			if decoder.Error() != nil {
				decoder.ResetError()
				break
			}

			objAndProps.PropertyIDs = append(objAndProps.PropertyIDs, bac_specs.PropertyIdentifier{})
			objAndProps.Data = append(objAndProps.Data, nil)

			objAndProps.PropertyIDs[propIdx].Type = bac_specs.PropertyType(val)
			objAndProps.PropertyIDs[propIdx].ArrayIndex = new(uint32)
			decoder.ContextValue(3, objAndProps.PropertyIDs[propIdx].ArrayIndex)
			err := decoder.Error()
			var e encoding.ErrorIncorrectTagID

			// The array index tag is optional, per BACnet spec
			if err != nil && errors.As(err, &e) {
				objAndProps.PropertyIDs[propIdx].ArrayIndex = nil
				decoder.ResetError()
			}

			decoder.ContextAbstractType(4, &objAndProps.Data[propIdx])

			propIdx++
		}

		decoder.ContextClosing(1)

		if err := decoder.Error(); err != nil {
			return err
		}

		rpm.ObjectIDs = append(rpm.ObjectIDs, objAndProps.ObjectID)
		rpm.PropertyIDs = append(rpm.PropertyIDs, objAndProps.PropertyIDs)
		rpm.Data = append(rpm.Data, objAndProps.Data)
	}

	return nil
}

//// preparePropertyIDsSlice ensures that rpm.PropertyIDs has one slice of PropertyIdentifier
//// structs for every rpm.ObjectIDs
//func (rpm *ReadPropertyMultiple) preparePropertyIDsSlice() {
//	lenDiff := len(rpm.ObjectIDs) - len(rpm.PropertyIDs)
//
//	for x := lenDiff; x > 0; x-- {
//		rpm.PropertyIDs = append(rpm.PropertyIDs, rpm.PropertyIDs[0])
//	}
//}
