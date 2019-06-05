package forwarder

type Forwarder interface {
	Send(payload *[]byte, target string) error
}
