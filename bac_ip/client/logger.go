package client

type Logger interface {
	Info(...interface{})
	Error(...interface{})
}

type NoOpLogger struct{}

func (NoOpLogger) Info(...interface{})  {}
func (NoOpLogger) Error(...interface{}) {}
