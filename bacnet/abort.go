package bacnet

//go:generate stringer -type=AbortReason
type AbortReason uint8

const (
	OtherAbortReason                         AbortReason = 0
	BufferOverflowAbortReason                AbortReason = 1
	InvalidAPDUInThisStateAbortReason        AbortReason = 2
	PreemptedByHigherPriorityTaskAbortReason AbortReason = 3
	SegmentationNotSupportedAbortReason      AbortReason = 4
	SecurityErrorAbortReason                 AbortReason = 5
	InsufficientSecurityAbortReason          AbortReason = 6
	WindowSizeOutOfRangeAbortReason          AbortReason = 7
	ApplicationExceededReplyTimeAbortReason  AbortReason = 8
	OutOfResourcesAbortReason                AbortReason = 9
	TsmTimeoutAbortReason                    AbortReason = 10
	APDUTooLongAbortReason                   AbortReason = 11
)
