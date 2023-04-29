package services

import (
	"github.com/toddyco/bacnet2go/bacnet"
	"github.com/toddyco/bacnet2go/internal/encoding"
)

type IAm struct {
	ObjectID            bacnet.ObjectID
	MaxApduLength       uint32
	SegmentationSupport bacnet.SegmentationSupport
	VendorID            uint32
}

func (iam IAm) MarshalBinary() ([]byte, error) {
	encoder := encoding.NewEncoder()
	encoder.AppData(iam.ObjectID)
	encoder.AppData(iam.MaxApduLength)
	encoder.AppData(iam.SegmentationSupport)
	encoder.AppData(iam.VendorID)
	return encoder.Bytes(), encoder.Error()
}

func (iam *IAm) UnmarshalBinary(data []byte) error {
	decoder := encoding.NewDecoder(data)
	decoder.AppData(&iam.ObjectID, nil)
	decoder.AppData(&iam.MaxApduLength, nil)
	decoder.AppData(&iam.SegmentationSupport, nil)
	decoder.AppData(&iam.VendorID, nil)
	return decoder.Error()
}
