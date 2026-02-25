package slb

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	slb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

type LSBListner struct {
	LoadBalancerID    string
	DomainExtensionID string
	ListenerPort      int32
}

func List(ctx context.Context, clt *slb.Client, domain string, regionID string) ([]LSBListner, error) {
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
		RegionId:           tea.String(regionID),
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

func checkSLB(ctx context.Context, clt *slb.Client, instanceID string, domain string) ([]LSBListner, error) {
	attrReq := slb.DescribeLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(instanceID),
	}
	attrResp, err := clt.DescribeLoadBalancerAttribute(&attrReq)
	if err != nil {
		return nil, fmt.Errorf("get SLB instance attribute failed, %w", err)
	}
	if attrResp.Body == nil || attrResp.Body.ListenerPortsAndProtocal == nil {
		return nil, nil
	}
	listners := make([]LSBListner, 0, len(attrResp.Body.ListenerPortsAndProtocal.ListenerPortAndProtocal))
	for _, p := range attrResp.Body.ListenerPortsAndProtocal.ListenerPortAndProtocal {
		if protocal := p.ListenerProtocal; protocal != nil && *protocal == "https" && p.ListenerPort != nil {
			domainReq := slb.DescribeDomainExtensionsRequest{
				LoadBalancerId: attrResp.Body.LoadBalancerId,
				ListenerPort:   p.ListenerPort,
			}
			domainResp, err := clt.DescribeDomainExtensions(&domainReq)
			if err != nil {
				return nil, fmt.Errorf("get SLB instance domain extensions failed, %w", err)
			}
			if domainResp.Body != nil && domainResp.Body.DomainExtensions != nil && len(domainResp.Body.DomainExtensions.DomainExtension) > 0 {
				for _, ext := range domainResp.Body.DomainExtensions.DomainExtension {
					extID := ext.DomainExtensionId
					if extID == nil {
						continue
					}
					if v := ext.Domain; v != nil && strings.HasSuffix(*v, domain) {
						slog.InfoContext(ctx, "found SLB listener", "load_balancer_id", instanceID, "port", *p.ListenerPort, "domain_extension_id", *extID, "domain", *v)
						listners = append(listners, LSBListner{
							LoadBalancerID:    instanceID,
							DomainExtensionID: *extID,
							ListenerPort:      *p.ListenerPort,
						})
						break
					}
				}
			} else {
				listenerDescribeReq := slb.DescribeLoadBalancerHTTPSListenerAttributeRequest{
					LoadBalancerId: attrResp.Body.LoadBalancerId,
					ListenerPort:   p.ListenerPort,
				}
				describeResp, err := clt.DescribeLoadBalancerHTTPSListenerAttribute(&listenerDescribeReq)
				if err != nil {
					return nil, fmt.Errorf("get SLB HTTPS listener attribute failed, %w", err)
				}
				if describeResp.Body != nil && describeResp.Body.Rules != nil {
					for _, rule := range describeResp.Body.Rules.Rule {
						if rule.Domain != nil && strings.HasSuffix(*rule.Domain, domain) {
							slog.InfoContext(ctx, "found SLB listener", "load_balancer_id", instanceID, "port", *p.ListenerPort, "domain", *rule.Domain)
							listners = append(listners, LSBListner{
								LoadBalancerID: instanceID,
								ListenerPort:   *p.ListenerPort,
							})
							break
						}
					}
				}
			}
		} else {
			slog.Info("invalid protocal", "protocal", p)
		}
	}
	return listners, nil
}
