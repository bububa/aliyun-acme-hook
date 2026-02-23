package util

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// CreateOpenAPIClient creates an Alibaba Cloud OpenAPI client with proper error handling
func CreateOpenAPIClient(ak, sk, region string) (*openapi.Config, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(ak),
		AccessKeySecret: tea.String(sk),
		RegionId:        tea.String(region),
	}

	// Basic validation of required parameters
	if ak == "" || sk == "" {
		return nil, fmt.Errorf("access key and secret key cannot be empty")
	}

	return config, nil
}
