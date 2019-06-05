package router

import (
	"encoding/json"

	"github.com/oliveagle/jsonpath"
)

// TODO: reuse lookup patterns
type Route struct {
	Path   string `json:"path"`
	Target string `json:"target"`
}

type RoutingTable struct {
	Routes []Route `json:"routes"`
}

func NewRoutingTable(data *[]byte) (*RoutingTable, error) {
	var routes RoutingTable
	err := json.Unmarshal(*data, &routes)
	return &routes, err
}

func (r *RoutingTable) FindMatches(data *string) (map[string]bool, error) {
	targets := map[string]bool{}
	var json_data interface{}
	err := json.Unmarshal([]byte(*data), &json_data)
	if err != nil {
		return targets, err
	}

	for _, route := range r.Routes {
		_, err := jsonpath.JsonPathLookup(json_data, route.Path)
		if err != nil {
			continue
		}
		targets[route.Target] = true
	}
	return targets, nil
}
