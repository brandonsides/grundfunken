package models

type InterpreterError struct {
	Message        string
	Underlying     error
	SourceLocation *SourceLocation
}

type SourceLocation struct {
	File         string
	LineNumber   int
	ColumnNumber int
}

func (e *InterpreterError) Error() string {
	return e.Message
}
