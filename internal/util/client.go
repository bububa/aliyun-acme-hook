package util

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
)

// CreateOpenAPIClient creates an Alibaba Cloud OpenAPI client with proper error handling
func CreateOpenAPIClient(ak, sk, region string) (*openapi.Config, error) {
	// Basic validation of required parameters
	if ak == "" || sk == "" {
		return nil, fmt.Errorf("access key and secret key cannot be empty")
	}
	config := new(openapi.Config)
	config.SetAccessKeyId(ak).SetAccessKeySecret(sk).SetRegionId(region)

	return config, nil
}
