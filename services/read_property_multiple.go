package services

import (
	"errors"
	"github.com/toddyco/bacnet2go/internal/encoding"
	"github.com/toddyco/bacnet2go/specs"
)

type ReadPropertyMultiple struct {
	ObjectIDs   []specs.ObjectID
	PropertyIDs [][]specs.PropertyIdentifier
	Data        [][]interface{} // will contain the response
	Errors      []error
}

type readPropertyMultipleItem struct {
	ObjectID    specs.ObjectID
	PropertyIDs []specs.PropertyIdentifier
	Data        []interface{}
	Errors      []error
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
	rpm.ObjectIDs = []specs.ObjectID{}
	rpm.PropertyIDs = [][]specs.PropertyIdentifier{}
	rpm.Data = [][]interface{}{}

	var propIdx int
	var err, propErr error
	var e encoding.ErrorIncorrectTagID

	decoder := encoding.NewDecoder(data)

	for { // Loop over objects
		objAndProps := readPropertyMultipleItem{
			ObjectID:    specs.ObjectID{},
			PropertyIDs: []specs.PropertyIdentifier{},
			Data:        []interface{}{},
		}

		decoder.ContextObjectID(0, &objAndProps.ObjectID)
		err = decoder.Error()

		if err != nil {
			if errors.As(err, &e) {
				decoder.ResetError()
				continue
			} else {
				break
			}
		}

		decoder.ContextOpening(1)
		err = decoder.Error()

		if err != nil {
			decoder.ResetError()
			continue
		}

		// There can be 1...n properties requested
		propIdx = 0
		propErr = nil

		for { // Loop over each of the 1...n properties
			var propertyType uint32
			decoder.ContextValue(2, &propertyType)

			if decoder.Error() != nil {
				decoder.ResetError()
				break
			}

			objAndProps.PropertyIDs = append(objAndProps.PropertyIDs, specs.PropertyIdentifier{})
			objAndProps.Data = append(objAndProps.Data, nil)
			objAndProps.PropertyIDs[propIdx].Type = specs.PropertyType(propertyType)
			objAndProps.PropertyIDs[propIdx].ArrayIndex = new(uint32)

			// Tag 3 is an array index tag; it's optional
			decoder.ContextValue(3, objAndProps.PropertyIDs[propIdx].ArrayIndex)
			err = decoder.Error()

			if err != nil && errors.As(err, &e) {
				objAndProps.PropertyIDs[propIdx].ArrayIndex = nil
				decoder.ResetError()
			}

			// Tag 4 is an expected return value. Tag 4 will not appear if tag 5 appears.
			decoder.ContextAbstractType(4, &objAndProps.Data[propIdx])
			err = decoder.Error()

			// Tag 5 is an error. Tag 5 will not appear if tag 4 appears.
			if err != nil && errors.As(err, &e) {
				if e.Got == 5 {
					decoder.ResetError()
					decoder.ContextAbstractType(5, &objAndProps.Data[propIdx])
					propErr = err
				} else {
					return err
				}
			} else {

			}

			propIdx++
		}

		decoder.ContextClosing(1)
		err = decoder.Error()

		if err != nil && errors.As(err, &e) {
			decoder.ResetError()
			continue
		}

		rpm.ObjectIDs = append(rpm.ObjectIDs, objAndProps.ObjectID)
		rpm.PropertyIDs = append(rpm.PropertyIDs, objAndProps.PropertyIDs)
		rpm.Data = append(rpm.Data, objAndProps.Data)
		rpm.Errors = append(rpm.Errors, propErr)
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
