package model

import "fmt"

// ValidationError 自定义错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Field '%s' is invalid: %s", e.Field, e.Message)
}
