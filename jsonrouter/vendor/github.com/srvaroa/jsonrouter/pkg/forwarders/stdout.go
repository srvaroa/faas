package forwarder

import (
	"fmt"
)

type StdoutForwarder struct {
}

// Useful for debugging mostly.
func (f StdoutForwarder) Send(payload *[]byte, target string) error {
	fmt.Printf(string(*payload))
	return nil
}
