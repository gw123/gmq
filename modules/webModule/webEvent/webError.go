package webEvent

type WebError struct {
	Inner error
	Msg   string
}

func (this *WebError) Error() string {
	return this.Inner.Error()
}

func (this *WebError) GetMsg() string {
	return this.Msg
}
