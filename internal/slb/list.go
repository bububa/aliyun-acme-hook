package slb

import (
	"context"
	"fmt"
	"strings"

	slb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

type LSBListner struct {
	LoadBalancerID    string
	DomainExtensionID string
}

func List(ctx context.Context, clt *slb.Client, domain string) ([]LSBListner, error) {
	var (
		pageNumber int32 = 1
		pageSize   int32 = 100
	)

	// Validate pagination parameters to prevent excessive API calls
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 100 // Default to a reasonable value
	}
	listReq := slb.DescribeLoadBalancersRequest{
		LoadBalancerStatus: tea.String("active"),
		PageSize:           tea.Int32(pageSize),
	}
	ret := make([]LSBListner, 0, pageSize)
	for {
		listReq.PageNumber = tea.Int32(pageNumber)
		resp, err := clt.DescribeLoadBalancers(&listReq)
		if err != nil {
			return nil, fmt.Errorf("get SLB instances list failed, %w", err)
		}
		if resp.Body == nil {
			break
		}
		if resp.Body.LoadBalancers != nil {
			for _, instance := range resp.Body.LoadBalancers.LoadBalancer {
				if v := instance.LoadBalancerId; v != nil {
					listners, err := checkSLB(ctx, clt, *v, domain)
					if err != nil {
						return nil, err
					}
					ret = append(ret, listners...)
				}
			}
		}
		var total int32
		if v := resp.Body.TotalCount; v != nil {
			total = *v
		}
		if total <= pageSize*pageNumber {
			break
		}
		pageNumber += 1
	}
	return ret, nil
}

func checkSLB(_ context.Context, clt *slb.Client, instanceID string, domain string) ([]LSBListner, error) {
	attrReq := slb.DescribeLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(instanceID),
	}
	attrResp, err := clt.DescribeLoadBalancerAttribute(&attrReq)
	if err != nil {
		return nil, fmt.Errorf("get SLB instance attribute failed, %w", err)
	}
	if attrResp.Body == nil || attrResp.Body.ListenerPortsAndProtocal != nil {
		return nil, nil
	}
	listners := make([]LSBListner, 0, len(attrResp.Body.ListenerPortsAndProtocal.ListenerPortAndProtocal))
	for _, p := range attrResp.Body.ListenerPortsAndProtocal.ListenerPortAndProtocal {
		if protocal := p.ListenerProtocal; protocal != nil && *protocal == "https" && p.ListenerPort != nil {
			domainReq := slb.DescribeDomainExtensionsRequest{
				LoadBalancerId: tea.String(instanceID),
				ListenerPort:   p.ListenerPort,
			}
			domainResp, err := clt.DescribeDomainExtensions(&domainReq)
			if err != nil {
				return nil, fmt.Errorf("get SLB instance domain extensions failed, %w", err)
			}
			if domainResp.Body != nil && domainResp.Body.DomainExtensions != nil {
				for _, ext := range domainResp.Body.DomainExtensions.DomainExtension {
					extID := ext.DomainExtensionId
					if extID == nil {
						continue
					}
					if v := ext.Domain; v != nil && strings.HasSuffix(*v, domain) {
						listners = append(listners, LSBListner{
							LoadBalancerID:    instanceID,
							DomainExtensionID: *extID,
						})
						break
					}
				}
			}
		}
	}
	return listners, nil
}
