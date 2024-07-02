package api_errors

type InternalServerError struct{}

func (e InternalServerError) Error() string {
	return "Internal server error"
}

type BadRequestError struct {
	Detail string
}

func (e BadRequestError) Error() string {
	message := "Bad Request"
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	return message
}
