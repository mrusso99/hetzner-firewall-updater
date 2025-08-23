package firewall

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func GetFirewall(ctx context.Context, client *hcloud.Client, name string) (*hcloud.Firewall, error) {
	fw, _, err := client.Firewall.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if fw == nil {
		return nil, fmt.Errorf("firewall %s not found", name)
	}
	return fw, nil
}

func UpdateFirewallRules(ctx context.Context, client *hcloud.Client, firewall *hcloud.Firewall) {

	rulesChanged := false
	currentIp, err := getIp(ctx)

	for i := range firewall.Rules {
		for j := range firewall.Rules[i].SourceIPs {
			if !firewall.Rules[i].SourceIPs[j].IP.Equal(net.ParseIP("0.0.0.0")) && !firewall.Rules[i].SourceIPs[j].IP.Equal(currentIp) {
				firewall.Rules[i].SourceIPs[j] = net.IPNet{IP: currentIp, Mask: net.CIDRMask(32, 32)}
				rulesChanged = true
			}
		}
	}

	if !rulesChanged {
		log.Println("Rules not changed, exiting from function")
		return
	}

	updatedRules := hcloud.FirewallSetRulesOpts{
		Rules: firewall.Rules,
	}

	act, resp, err := client.Firewall.SetRules(ctx, firewall, updatedRules)
	if err != nil {
		log.Fatal("Error:", err)
	} else {
		log.Printf("Action: %+v\n", act)
		log.Printf("HTTP Response Status: %s\n", resp.Status)
	}
}
