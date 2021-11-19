package resources

import (
	"errors"
	"fmt"
)

// Campaign struct
type Campaign struct {
	method   string
	metadata map[string]interface{}
}

func (c *Campaign) objectName() string {
	const obName = "Campaign"
	return obName
}

// NewCampaign writes that new Campaign instance
func NewCampaign(metadata map[string]interface{}) (*Campaign, error) {
	rawMethod, ok := metadata["method"]
	if !ok {
		return nil, errors.New("missing required parameters: method")
	}
	method, ok := rawMethod.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return &Campaign{
		method:   method,
		metadata: metadata,
	}, nil
}

// getMetadata mold campaign get metadata
func (c *Campaign) getMetadata() (map[string]interface{}, error) {
	idIF, ok := c.metadata["id"]
	if !ok {
		return nil, fmt.Errorf("id is missing")
	}
	pathParam, ok := idIF.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return buildMetadata(c.method, c.objectName(), pathParam, nil, "", "campaign_get"), nil
}

// BuildMetadata
func (c *Campaign) BuildMetadata() (map[string]interface{}, error) {
	switch c.method {
	case "get":
		return c.getMetadata()
	}
	return nil, fmt.Errorf("invalid method: %s", c.method)
}

func buildMetadata(method, object, pathParam string, queryParams map[string]string, body string, connectionKey string) map[string]interface{} {
	metadata := map[string]interface{}{
		"method":         method,
		"object":         object,
		"connection_key": connectionKey,
	}
	if len(pathParam) > 0 {
		metadata["path_param"] = pathParam
	}
	if queryParams != nil {
		metadata["query_params"] = queryParams
	}
	if body != "" {
		metadata["body"] = body
	}
	return metadata
}
