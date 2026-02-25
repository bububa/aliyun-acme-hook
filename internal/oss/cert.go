package oss

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/model"
)

func Certificate(ctx context.Context, cfg *config.AliyunConfig, cert *model.Cert) error {
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AK, cfg.SK)).
		WithRegion(cfg.Region)
	client := oss.NewClient(ossCfg)
	paginator := client.NewListBucketsPaginator(&oss.ListBucketsRequest{})
	var i int
	for paginator.HasNext() {
		i++
		page, err := paginator.NextPage(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get OSS bucket list", "page", i, "error", err)
			return fmt.Errorf("failed to get OSS bucket list, %w", err)
		}
		// Print the bucket found
		for _, b := range page.Buckets {
			clt := client
			if *b.Region != cfg.Region {
				ossCfg := oss.LoadDefaultConfig().
					WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AK, cfg.SK)).
					WithRegion(*b.Region)
				clt = oss.NewClient(ossCfg)
			}
			if err := certCname(ctx, clt, cert, *b.Name); err != nil {
				slog.ErrorContext(ctx, "certificate OSS bucket domain failed", "error", err, "bucket", *b.Name)
				return err
			}
		}
	}
	return nil
}

func certCname(ctx context.Context, clt *oss.Client, cert *model.Cert, bucketName string) error {
	listReq := oss.ListCnameRequest{
		Bucket: &bucketName,
	}
	listResp, err := clt.ListCname(ctx, &listReq)
	if err != nil {
		return fmt.Errorf("list bucket cnames failed, %w", err)
	}
	for _, cname := range listResp.Cnames {
		if !strings.HasSuffix(*cname.Domain, cert.Domain) {
			continue
		}
		if cname.Certificate == nil {
			slog.WarnContext(ctx, "cname custom domain is controled by CDN", "domain", *cname.Domain, "bucket", bucketName)
			continue
		}
		certConfig := &oss.CertificateConfiguration{
			Force: tea.Bool(true),
		}
		if cert.ID > 0 {
			certConfig.CertId = tea.String(strconv.FormatInt(cert.ID, 10))
		} else {
			certConfig.Certificate = tea.String(string(cert.FullChain))
			certConfig.PrivateKey = tea.String(string(cert.Key))
		}
		certReq := oss.PutCnameRequest{
			Bucket: &bucketName,
			BucketCnameConfiguration: &oss.BucketCnameConfiguration{
				Cname: &oss.Cname{
					Domain:                   cname.Domain,
					CertificateConfiguration: certConfig,
				},
			},
		}
		if _, err := clt.PutCname(ctx, &certReq); err != nil {
			slog.ErrorContext(ctx, "update OSS cname failed", "error", err, "domain", *cname.Domain, "bucket", bucketName)
			return fmt.Errorf("update OSS cname failed for domain:%s, %w", *cname.Domain, err)
		}
	}
	return nil
}
