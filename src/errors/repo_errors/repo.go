package repo_errors

import "fmt"

type OperationError struct{}

func (e OperationError) Error() string {
	return "Error while performing operation"
}

type ObjectNotFoundError struct {
	ParamName string
	Value     interface{}
}

func (e ObjectNotFoundError) Error() string {
	return fmt.Sprintf("Object not found. %s: %v", e.ParamName, e.Value)
}

type ObjectAlreadyExistsError struct{}

func (e ObjectAlreadyExistsError) Error() string {
	return "Object already exists"
}
