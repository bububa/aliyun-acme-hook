package slb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	slb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/model"
	"github.com/bububa/aliyun-acme-hook/internal/util"
)

func Certificate(ctx context.Context, cfg *config.AliyunConfig, cert *model.Cert) error {
	if cert.Domain == "" {
		return errors.New("missing cert domain environment variables")
	}

	// 3. 初始化阿里云客户端
	// 建议通过环境变量获取 Ali_Key 和 Ali_Secret
	apiConfig, err := util.CreateOpenAPIClient(cfg.AK, cfg.SK, cfg.Region)
	if err != nil {
		return fmt.Errorf("failed to create SLB API config: %w", err)
	}
	client, err := slb.NewClient(apiConfig)
	if err != nil {
		return fmt.Errorf("failed to create SLB client: %w", err)
	}
	req := slb.UploadServerCertificateRequest{
		AliCloudCertificateName: tea.String(cert.Name),
		AliCloudCertificateId:   tea.String(strconv.FormatInt(cert.ID, 10)),
	}
	uploadResp, err := client.UploadServerCertificate(&req)
	if err != nil {
		slog.ErrorContext(ctx, "upload LSB cert failed", "domain", cert.Domain, "error", err)
		return fmt.Errorf("set SLB domain cert failed, %w", err)
	}
	list, err := List(ctx, client, cert.Domain)
	if err != nil {
		slog.ErrorContext(ctx, "list LSB listeners failed", "domain", cert.Domain, "error", err)
		return fmt.Errorf("upload SLB domain cert failed, %w", err)
	}
	for _, v := range list {
		certReq := slb.SetDomainExtensionAttributeRequest{
			DomainExtensionId:   tea.String(v.LoadBalancerID),
			ServerCertificateId: uploadResp.Body.ServerCertificateId,
		}
		if _, err := client.SetDomainExtensionAttribute(&certReq); err != nil {
			slog.ErrorContext(ctx, "set LSB domain extension failed", "domain", cert.Domain, "error", err)
			return fmt.Errorf("set SLB domain extension cert failed, %w", err)
		}
	}
	slog.InfoContext(ctx, "certicate SLB domain done!", "domain", cert.Domain)
	return nil
}
