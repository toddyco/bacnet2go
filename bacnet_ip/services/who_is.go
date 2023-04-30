package services

import (
	"fmt"
	"github.com/toddyco/bacnet2go/internal/encoding"
	"github.com/toddyco/bacnet2go/specs"
)

type WhoIs struct {
	Low, High *uint32 // may be null if we want to check all range
}

func (w WhoIs) MarshalBinary() ([]byte, error) {
	encoder := encoding.NewEncoder()
	if w.Low != nil && w.High != nil {
		if *w.Low > specs.MaxInstance || *w.High > specs.MaxInstance {
			return nil, fmt.Errorf("invalid WhoIs range: [%d, %d]: max value is %d", *w.Low, *w.High, specs.MaxInstance)
		}
		if *w.Low > *w.High {
			return nil, fmt.Errorf("invalid WhoIs range: [%d, %d]: low limit is higher than high limit", *w.Low, *w.High)
		}
		encoder.ContextUnsigned(0, *w.Low)
		encoder.ContextUnsigned(1, *w.High)
	}
	return encoder.Bytes(), encoder.Error()
}

func (w *WhoIs) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		// If data is empty, the whoIs request is a full range
		// check. So keep the low and high pointer nil
		return nil
	}
	w.Low = new(uint32)
	w.High = new(uint32)
	decoder := encoding.NewDecoder(data)
	decoder.ContextValue(byte(0), w.Low)
	decoder.ContextValue(byte(1), w.High)
	return decoder.Error()
}
