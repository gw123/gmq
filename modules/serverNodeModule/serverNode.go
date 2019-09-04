package serverNodeModule

type ServerNode struct {
	Url         string
	Host        string
	NodeName    string
	httpRequest SignHttpRequest
}

func (this *ServerNode) Login() error {

	return nil
}
