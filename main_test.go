package main

import (
	"reflect"
	"testing"
)

func TestMainProcessing(t *testing.T) {
	tests := []struct {
		FromMetadata     map[string]interface{}
		ExpectToMetadata map[string]interface{}
	}{
		// get
		{
			FromMetadata: map[string]interface{}{
				"connection_type": "request",
				"method":          "get",
				"object":          "Campaign",
				"id":              "test_id",
			},
			ExpectToMetadata: map[string]interface{}{
				"method":         "get",
				"object":         "Campaign",
				"connection_key": "campaign_get",
				"path_param":     "test_id",
			},
		},
	}
	for i, tt := range tests {
		got, err := handle(tt.FromMetadata)
		if err != nil {
			t.Errorf("failed to handle: %w", err)
		}
		if !reflect.DeepEqual(got, tt.ExpectToMetadata) {
			t.Errorf("%d# metadata got = %v, want %v", i, got, tt.ExpectToMetadata)
		}
	}
}
