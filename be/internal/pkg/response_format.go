package pkg

import "context"

type ResponseFormat struct {
	Code      int                    `json:"code"`
	ErrorMsg  string                 `json:"error_msg"`
	Result    map[string]interface{} `json:"result"`
	RequestID string                 `json:"request_id"`
}

func NewResponseFormat(ctx context.Context, code int, errorMsg string, result map[string]interface{}) *ResponseFormat {
	return &ResponseFormat{
		Code:      code,
		ErrorMsg:  errorMsg,
		Result:    result,
		RequestID: GetRequestID(ctx),
	}
}
