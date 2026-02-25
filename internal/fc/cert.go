package fc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc "github.com/alibabacloud-go/fc-20230330/v4/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/model"
)

func Certificate(ctx context.Context, cfg *config.AliyunConfig, cert *model.Cert) error {
	apiConfig := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AK),
		AccessKeySecret: tea.String(cfg.SK),
		RegionId:        tea.String(cfg.Region),
		Endpoint:        tea.String(cfg.AccountID + "." + cfg.Region + ".fc.aliyuncs.com"),
	}
	client, err := fc.NewClient(apiConfig)
	if err != nil {
		return fmt.Errorf("failed to create FC client: %w", err)
	}
	var nextToken string
	listReq := &fc.ListCustomDomainsRequest{
		Limit: tea.Int32(100),
	}
	for {
		listReq.SetNextToken(nextToken)
		listResp, err := client.ListCustomDomains(listReq)
		if err != nil {
			return fmt.Errorf("get FC custom domains failed, %w", err)
		}
		if listResp.Body == nil {
			break
		}
		for _, domain := range listResp.Body.CustomDomains {
			if domain.DomainName == nil {
				continue
			}
			domainName := *domain.DomainName
			if !strings.HasSuffix(domainName, cert.Domain) {
				continue
			}
			certReq := fc.UpdateCustomDomainRequest{
				Body: &fc.UpdateCustomDomainInput{
					AuthConfig: domain.AuthConfig,
					CertConfig: &fc.CertConfig{
						CertName:    tea.String(cert.Name),
						Certificate: tea.String(string(cert.FullChain)),
						PrivateKey:  tea.String(string(cert.Key)),
					},
					CorsConfig:  domain.CorsConfig,
					Protocol:    domain.Protocol,
					RouteConfig: domain.RouteConfig,
					TlsConfig:   domain.TlsConfig,
					WafConfig:   domain.WafConfig,
				},
			}
			if _, err := client.UpdateCustomDomain(domain.DomainName, &certReq); err != nil {
				slog.ErrorContext(ctx, "update FC custom domain failed", "domain", domainName, "error", err)
				return fmt.Errorf("update FC custom domain:%s failed, %w", domainName, err)
			} else {
				slog.InfoContext(ctx, "update FC custom domain succeed", "domain", domainName)
			}
		}
		if listResp.Body.NextToken == nil || *listResp.Body.NextToken == "" {
			break
		}
		nextToken = *listResp.Body.NextToken
	}
	return nil
}
