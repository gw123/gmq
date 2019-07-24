package webEvent

type WebError struct {
	Inner error
	*Response
}

func NewWebError(response *Response, err error) *WebError {
	this := &WebError{
		Response: response,
		Inner:    err,
	}
	return this
}

func (this *WebError) Error() string {
	return this.Inner.Error()
}

func (this *WebError) GetResponse() *Response {
	return this.Response
}
