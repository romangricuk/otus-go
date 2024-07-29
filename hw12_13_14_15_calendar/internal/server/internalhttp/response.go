package internalhttp

type Response struct {
	Data      interface{} `json:"data,omitempty"`
	Errors    []string    `json:"errors,omitempty"`
	Status    int         `json:"status"`
	RequestID string      `json:"requestId"`
}

func NewResponse(data interface{}, errors []string, status int) Response {
	return Response{
		Data:   data,
		Errors: errors,
		Status: status,
	}
}
