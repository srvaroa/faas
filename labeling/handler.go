package function

import (
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v26/github"
	"github.com/openfaas-incubator/go-function-sdk"
)

func Handle(req handler.Request) (handler.Response, error) {
	log.Println(req.Header) // Check function logs for the request headers

	var err error
	eventName := req.Header["Http_X_Github_Event"][0]

	// Wrap the shared transport for use with the integration ID (33575)
	// authenticating with installation ID 1199612, which is my test
	// repo.  I assume this should be set based on the repo we want to
	// act on (or the owner)
	// TODO: get this right
	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		1199612, 33576,
		"private_key_file.pem")

	if err != nil {
		log.Printf("Failed to load key: %s", err)
		return failure("Oops, something broke", err)
	}

	gh := github.NewClient(&http.Client{Transport: itr})
	l := NewLabeller(gh)

	if err = l.HandleEvent(eventName, &req.Body); err != nil {
		log.Printf("Failed to process event: %s", err)
		return failure("Oops, something broke", err)
	}
	return success("OK")
}

func failure(msg string, err error) (handler.Response, error) {
	return response(msg, err)
}

func success(msg string) (handler.Response, error) {
	return response(msg, nil)
}

func response(msg string, err error) (handler.Response, error) {
	return handler.Response{
		Body: []byte(msg),
		Header: map[string][]string{
			"X-Served-By": []string{"srvaroa.o6s.io/labeling"},
		},
	}, err
}
