package customerror

import "fmt"

type CustomError struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("error code %d: %s", ce.ErrorCode, ce.ErrorMessage)
}
