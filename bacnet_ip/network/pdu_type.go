package network

// //go:generate stringer -type=PDUType
type PDUType byte

const (
	ConfirmedServiceRequest   PDUType = 0
	UnconfirmedServiceRequest PDUType = 0x10
	SimpleAck                 PDUType = 0x20
	ComplexAck                PDUType = 0x30
	SegmentAck                PDUType = 0x40
	Error                     PDUType = 0x50
	Reject                    PDUType = 0x60
	Abort                     PDUType = 0x70
)

// IsType is necessary because if the type is abort, then the 8th bit can be
// a 1 or 0, so simply doing an equivalency check for the PDUType with one of
// the above constants is not sufficient.
func (pt PDUType) IsType(t PDUType) bool {
	return pt&t == t
}

func (pt PDUType) SupportsInvokeID() bool {
	return pt == SimpleAck || pt == ComplexAck || pt == Error || pt.IsType(Abort)
}
