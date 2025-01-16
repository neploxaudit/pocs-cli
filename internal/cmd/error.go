package cmd

type UsageError struct {
	Command string
	Message string
}

func (e *UsageError) Error() string {
	return e.Message
}

type ExecError struct {
	Command string
	Message string
}

func (e *ExecError) Error() string {
	return e.Message
}
