package k8s

import (
	"encoding/json"
	"github.com/golang/glog"
	"strings"
)

type filterObject struct {
	filterType string
	path       []string
	value      interface{}
}

type FilterCollection struct {
	filters map[string][]filterObject
}

func NewFilterCollection() *FilterCollection {
	return &FilterCollection{filters: make(map[string][]filterObject)}
}

func (collection FilterCollection) addFilter(kind string, filter filterObject) {
	filters := collection.filters[kind]
	if filters == nil {
		filters = make([]filterObject, 1)
		filters[0] = filter
	} else {
		filters = append(filters, filter)
	}

	collection.filters[kind] = filters
}

func (collection FilterCollection) AddRemoveFilter(kind string, path ...string) {
	filter := filterObject{filterType: "remove", path: path}

	collection.addFilter(kind, filter)
}

func (collection FilterCollection) AddReplaceFilter(value interface{}, kind string, path ...string) {
	filter := filterObject{filterType: "replace", path: path, value: value}

	collection.addFilter(kind, filter)
}

func (collection FilterCollection) AddEmptyFilter(kind string, path ...string) {
	filter := filterObject{filterType: "replace", path: path, value: ""}

	collection.addFilter(kind, filter)
}

func (collection FilterCollection) ApplyToJSON(body []byte) ([]byte, bool) {
	// Get the Kind object, fails gracefully
	kind, err := KindFromJSON(body)
	if err != nil {
		glog.Warningf("%v", err)
		return body, false
	}

	// Apply filters if necessary
	if !collection.filter(kind.data) {
		// No filters applied, don't marshal the JSON and returns the unchanged body
		return body, false
	}

	// Return the filtered JSON, fails gracefully
	filteredBody, err := json.Marshal(kind.data)
	if err != nil {
		glog.Warningf("%v", err)
		return body, false
	}

	return filteredBody, true
}

func (collection FilterCollection) filter(data map[string]interface{}) bool {
	// Check if we have a valid Kind object
	kindType, ok := data["kind"].(string)
	if !ok {
		return false
	}

	// If it is a list filter the items.
	// A List has the suffix List (e.g. List, PodList, ServiceList, NodeList, etc.).
	// At the moment we can't apply filters to lists, only yo list items
	if strings.HasSuffix(kindType, "List") {
		var isGenericList = kindType == "List"

		var itemsFilters = []filterObject{}
		if !isGenericList {
			// Check if we have filters for the expected list item type
			itemsType := strings.TrimSuffix(kindType, "List")
			itemsFilters = collection.filters[itemsType]
			if itemsFilters == nil || len(itemsFilters) == 0 {
				return false
			}
		}

		// Check if we have items and items is an array
		items, ok := data["items"].([]interface{})
		if !ok {
			return false
		}

		// Iterate over all the items
		var filtered = false
		for _, item := range items {
			// Check if an item is a JSON object
			data, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			// Generic lists can have any kind of objects
			if isGenericList {
				// Filter the object recursively
				if collection.filter(data) {
					filtered = true
				}
			} else {
				// Filter the object
				if filterData(data, itemsFilters) {
					filtered = true
				}
			}
		}

		return filtered
	} else {
		// If it isn't a list apply the filters (if any)
		filters := collection.filters[kindType]
		if filters != nil && len(filters) > 0 {
			return filterData(data, filters)
		} else {
			// We have no filters, just return
			return false
		}
	}
}

func filterData(data map[string]interface{}, filters []filterObject) bool {
	var filtered = false
	for _, filter := range filters {
		switch filter.filterType {

		case "remove":
			if remove(data, filter.path) {
				filtered = true
			}

		case "replace":
			if replace(data, filter.value, filter.path) {
				filtered = true
			}

		}
	}

	return filtered
}

func remove(data map[string]interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	} else if data[path[0]] == nil {
		return false
	} else if len(path) == 1 {
		delete(data, path[0])

		return true
	}

	first, ok := data[path[0]].(map[string]interface{})
	if !ok {
		return false
	}
	return remove(first, path[1:])
}

func replace(data map[string]interface{}, value interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	}

	if len(path) == 1 {
		if data[path[0]] != nil {
			data[path[0]] = value
			return true
		}

		return false
	} else if data[path[0]] == nil {
		return false
	}

	first, ok := data[path[0]].(map[string]interface{})
	if !ok {
		return false
	}
	return replace(first, value, path[1:])
}
