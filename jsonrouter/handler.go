package function

import (
	"fmt"
	"os"

	"github.com/srvaroa/jsonrouter/pkg/forwarders"
	"github.com/srvaroa/jsonrouter/pkg/router"
)

func Handle(req []byte) string {

	configData := os.Getenv("config")
	if len(configData) == 0 {
		return "Could not read route configuration (or empty)"
		os.Exit(1)
	}

	configDataBytes := []byte(configData)
	routes, err := router.NewRoutingTable(&configDataBytes)
	if err != nil {
		os.Exit(1)
	}

	json_raw := string(req)
	res, err := routes.FindMatches(&json_raw)
	if err != nil {
		os.Exit(1)
	}

	dump := forwarder.HttpForwarder{}
	for endpoint, _ := range res {
		err = dump.Send(&req, endpoint)
		if err != nil {
			return fmt.Sprintf("Failed to forward event %s", err)
		}
	}

	return "OK"
}
