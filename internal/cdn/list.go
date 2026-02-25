package cdn

import (
	"context"
	"fmt"

	cdn "github.com/alibabacloud-go/cdn-20180510/v9/client"
	"github.com/alibabacloud-go/tea/dara"
)

func GetDomains(ctx context.Context, clt *cdn.Client, domain string) ([]string, error) {
	var (
		pageNumber int32 = 1
		pageSize   int32 = 500
	)

	// Validate pagination parameters to prevent excessive API calls
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 500 // Default to a reasonable value
	}
	req := new(cdn.DescribeUserDomainsRequest)
	req.SetDomainName(domain).SetDomainSearchType("suf_match").SetDomainStatus("online").SetPageSize(pageSize)
	ret := make([]string, 0, pageSize)
	for {
		req.SetPageNumber(pageNumber)
		resp, err := clt.DescribeUserDomainsWithContext(ctx, req, new(dara.RuntimeOptions))
		if err != nil {
			return nil, fmt.Errorf("get user domains failed, %w", err)
		}
		if resp.Body == nil {
			break
		}
		if resp.Body.Domains != nil {
			for _, page := range resp.Body.Domains.PageData {
				if v := page.DomainName; v != nil {
					ret = append(ret, *v)
				}
			}
		}
		var total int64
		if v := resp.Body.TotalCount; v != nil {
			total = *v
		}
		if total <= int64(pageSize*pageNumber) {
			break
		}
		pageNumber += 1
	}
	return ret, nil
}
