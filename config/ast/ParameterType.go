package ast

import "fmt"

// ParameterType represents the type of a parameter
type ParameterType int

const (
	// StringParameter represents a parameter that is an arbitrary string
	StringParameter ParameterType = iota
	// ConnectionParameter represents a parameter that consists of an IP or hostname and port number
	ConnectionParameter ParameterType = iota
	// BooleanParameter represents a parameter that is either true or false
	BooleanParameter ParameterType = iota
)

// StringParameterData is the data contained in a StringParameter
type StringParameterData struct {
	Value string
}

func (p StringParameterData) String() string {
	return p.Value
}

// ConnectionParameterData is the data contained in a ConnectionParameter
type ConnectionParameterData struct {
	Host string
	Port int
}

func (p ConnectionParameterData) String() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

// BooleanParameterData is the data contained in a BooleanParameter
type BooleanParameterData struct {
	Value bool
}

func (b BooleanParameterData) String() string {
	return fmt.Sprint(b.Value)
}
