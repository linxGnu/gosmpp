package controller

import "github.com/microcosm-cc/bluemonday"

// ResponseError ...
type ResponseError struct {
	Code    int
	Message string
}

// Response ...
type Response struct {
	Error ResponseError
	Data  interface{}
}

// SetData ...
func (c *Response) SetData(_dat interface{}) {
	c.Data = _dat
}

// SetCodeMessage ...
func (c *Response) SetCodeMessage(code int, message string) {
	c.Error.Code = code
	c.Error.Message = string(message)
}

var santizer = bluemonday.UGCPolicy()

// Sanitize do antixss to input
func Sanitize(s string) string {
	return santizer.Sanitize(s)
}

// SanitizeBytes do antixss with bytes
func SanitizeBytes(data []byte) []byte {
	return santizer.SanitizeBytes(data)
}
