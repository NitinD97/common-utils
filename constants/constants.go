package constants

const (
	RequestId = "requestId"
)

type Enum string

const (
	EnumServiceName Enum = "serviceName"
)

func (s Enum) ToString() string {
	return string(s)
}
