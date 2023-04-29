package bacnet

// ErrorClass is the type of error send by an object after a request
type ErrorClass uint16

// ErrorCode identifies precisely the error
type ErrorCode uint32

//go:generate stringer -type=ErrorClass
const (
	DeviceError        ErrorClass = 0x00
	ObjectError        ErrorClass = 0x01
	PropertyError      ErrorClass = 0x02
	ResourcesError     ErrorClass = 0x03
	SecurityError      ErrorClass = 0x04
	ServicesError      ErrorClass = 0x05
	VTError            ErrorClass = 0x06
	CommunicationError ErrorClass = 0x07
)

//go:generate stringer -type=ErrorCode
const (
	Other                              ErrorCode = 0x00
	DeviceBusy                         ErrorCode = 0x03
	ConfigurationInProgress            ErrorCode = 0x02
	OperationalProblem                 ErrorCode = 0x19
	DynamicCreationNotSupported        ErrorCode = 0x04
	NoObjectsOfSpecifiedType           ErrorCode = 0x11
	ObjectDeletionNotPermitted         ErrorCode = 0x17
	ObjectIdentifierAlreadyExists      ErrorCode = 0x18
	ReadAccessDenied                   ErrorCode = 0x1B
	UnknownObject                      ErrorCode = 0x1F
	UnsupportedObjectType              ErrorCode = 0x24
	CharacterSetNotSupported           ErrorCode = 0x29
	DatatypeNotSupported               ErrorCode = 0x2F
	InconsistentSelectionCriterion     ErrorCode = 0x08
	InvalidArrayIndex                  ErrorCode = 0x2A
	InvalidDataType                    ErrorCode = 0x09
	NotCovProperty                     ErrorCode = 0x2C
	OptionalFunctionalityNotSupported  ErrorCode = 0x2D
	PropertyIsNotAnArray               ErrorCode = 0x32
	UnknownProperty                    ErrorCode = 0x20
	ValueOutOfRange                    ErrorCode = 0x25
	WriteAccessDenied                  ErrorCode = 0x28
	NoSpaceForObject                   ErrorCode = 0x12
	NoSpaceToAddListElement            ErrorCode = 0x13
	NoSpaceToWriteProperty             ErrorCode = 0x14
	AuthenticationFailed               ErrorCode = 0x01
	IncompatibleSecurityLevels         ErrorCode = 0x06
	InvalidOperatorName                ErrorCode = 0x0C
	KeyGenerationError                 ErrorCode = 0x0F
	PasswordFailure                    ErrorCode = 0x1A
	SecurityNotSupported               ErrorCode = 0x1C
	Timeout                            ErrorCode = 0x1E
	CovSubscriptionFailed              ErrorCode = 0x2B
	DuplicateName                      ErrorCode = 0x30
	DuplicateObjectID                  ErrorCode = 0x31
	FileAccessDenied                   ErrorCode = 0x05
	InconsistentParameters             ErrorCode = 0x07
	InvalidConfigurationData           ErrorCode = 0x2E
	InvalidFileAccessMethod            ErrorCode = 0x0A
	InvalidFileStartPosition           ErrorCode = 0x0B
	InvalidParameterDataType           ErrorCode = 0x0D
	InvalidTimeStamp                   ErrorCode = 0x0E
	MissingRequiredParameter           ErrorCode = 0x10
	PropertyIsNotAList                 ErrorCode = 0x16
	ServiceRequestDenied               ErrorCode = 0x1D
	UnknownVtClass                     ErrorCode = 0x22
	UnknownVtSession                   ErrorCode = 0x23
	NoVtSessionsAvailable              ErrorCode = 0x15
	VtSessionAlreadyClosed             ErrorCode = 0x26
	VtSessionTerminationFailure        ErrorCode = 0x27
	AbortBufferOverflow                ErrorCode = 0x33
	AbortInvalidApduInThisState        ErrorCode = 0x34
	AbortPreemptedByHigherPriorityTask ErrorCode = 0x35
	AbortSegmentationNotSupported      ErrorCode = 0x36
	AbortProprietary                   ErrorCode = 0x37
	AbortOther                         ErrorCode = 0x38
	InvalidTag                         ErrorCode = 0x39
	NetworkDown                        ErrorCode = 0x3A
	RejectBufferOverflow               ErrorCode = 0x3B
	RejectInconsistentParameters       ErrorCode = 0x3C
	RejectInvalidParameterDataType     ErrorCode = 0x3D
	RejectInvalidTag                   ErrorCode = 0x3E
	RejectMissingRequiredParameter     ErrorCode = 0x3F
	RejectParameterOutOfRange          ErrorCode = 0x40
	RejectTooManyArguments             ErrorCode = 0x41
	RejectUndefinedEnumeration         ErrorCode = 0x42
	RejectUnrecognizedService          ErrorCode = 0x43
	RejectProprietary                  ErrorCode = 0x44
	RejectOther                        ErrorCode = 0x45
	UnknownDevice                      ErrorCode = 0x46
	UnknownRoute                       ErrorCode = 0x47
	ValueNotInitialized                ErrorCode = 0x48
	InvalidEventState                  ErrorCode = 0x49
	NoAlarmConfigured                  ErrorCode = 0x4A
	LogBufferFull                      ErrorCode = 0x4B
	LoggedValuePurged                  ErrorCode = 0x4C
	NoPropertySpecified                ErrorCode = 0x4D
	NotConfiguredForTriggeredLogging   ErrorCode = 0x4E
	UnknownSubscription                ErrorCode = 0x4F
	ParameterOutOfRange                ErrorCode = 0x50
	ListElementNotFound                ErrorCode = 0x51
	Busy                               ErrorCode = 0x52
	CommunicationDisabled              ErrorCode = 0x53
	Success                            ErrorCode = 0x54
	AccessDenied                       ErrorCode = 0x55
	BadDestinationAddress              ErrorCode = 0x56
	BadDestinationDeviceID             ErrorCode = 0x57
	BadSignature                       ErrorCode = 0x58
	BadSourceAddress                   ErrorCode = 0x59
	BadTimestamp                       ErrorCode = 0x5A
	CannotUseKey                       ErrorCode = 0x5B
	CannotVerifyMessageID              ErrorCode = 0x5C
	CorrectKeyRevision                 ErrorCode = 0x5D
	DestinationDeviceIDRequired        ErrorCode = 0x5E
	DuplicateMessage                   ErrorCode = 0x5F
	EncryptionNotConfigured            ErrorCode = 0x60
	EncryptionRequired                 ErrorCode = 0x61
	IncorrectKey                       ErrorCode = 0x62
	InvalidKeyData                     ErrorCode = 0x63
	KeyUpdateInProgress                ErrorCode = 0x64
	MalformedMessage                   ErrorCode = 0x65
	NotKeyServer                       ErrorCode = 0x66
	SecurityNotConfigured              ErrorCode = 0x67
	SourceSecurityRequired             ErrorCode = 0x68
	TooManyKeys                        ErrorCode = 0x69
	UnknownAuthenticationType          ErrorCode = 0x6A
	UnknownKey                         ErrorCode = 0x6B
	UnknownKeyRevision                 ErrorCode = 0x6C
	UnknownSourceMessage               ErrorCode = 0x6D
	NotRouterToDnet                    ErrorCode = 0x6E
	RouterBusy                         ErrorCode = 0x6F
	UnknownNetworkMessage              ErrorCode = 0x70
	MessageTooLong                     ErrorCode = 0x71
	SecurityErrorCode                  ErrorCode = 0x72
	AddressingError                    ErrorCode = 0x73
	WriteBdtFailed                     ErrorCode = 0x74
	ReadBdtFailed                      ErrorCode = 0x75
	RegisterForeignDeviceFailed        ErrorCode = 0x76
	ReadFdtFailed                      ErrorCode = 0x77
	DeleteFdtEntryFailed               ErrorCode = 0x78
	DistributeBroadcastFailed          ErrorCode = 0x79
	UnknownFileSize                    ErrorCode = 0x7A
	AbortApduTooLong                   ErrorCode = 0x7B
	AbortApplicationExceededReplyTime  ErrorCode = 0x7C
	AbortOutOfResources                ErrorCode = 0x7D
	AbortTSMTimeout                    ErrorCode = 0x7E
	AbortWindowSizeOutOfRange          ErrorCode = 0x7F
	FileFull                           ErrorCode = 0x80
	InconsistentConfiguration          ErrorCode = 0x81
	InconsistentObjectType             ErrorCode = 0x82
	InternalError                      ErrorCode = 0x83
	NotConfigured                      ErrorCode = 0x84
	OutOfMemory                        ErrorCode = 0x85
	ValueTooLong                       ErrorCode = 0x86
	AbortInsufficientSecurity          ErrorCode = 0x87
	AbortSecurityError                 ErrorCode = 0x88
)
