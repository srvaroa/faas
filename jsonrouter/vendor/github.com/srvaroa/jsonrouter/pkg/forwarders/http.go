package forwarder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpForwarder struct {
}

func (f HttpForwarder) Send(payload *[]byte, url string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(*payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Target %s Returned %s", url, resp.Status)
	}

	return nil
}
