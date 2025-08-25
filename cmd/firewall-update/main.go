package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/mrusso99/hetzner-firewall-update/internal/firewall"
	"github.com/prometheus-community/pro-bing"
	"github.com/spf13/viper"
)

func main() {
	log.Println("Running Hetzner Firewall Updater...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	viper.AutomaticEnv()

	// Optional: load a config file
	viper.SetConfigName(".env") // file name without extension
	viper.SetConfigType("env")
	viper.AddConfigPath(".") // current directory

	_ = viper.ReadInConfig() // ignore error if file doesn't exist

	// Get Hetzner token
	token := viper.GetString("HETZNER_TOKEN")
	if token == "" {
		log.Fatal("HETZNER_TOKEN is required")
	}

	host := viper.GetString("HOST")
	if host == "" {
		log.Fatal("HOST is required")
	}

	if len(os.Args) < 2 {
		log.Fatal("You must specify at least one firewall name as argument")
	}
	firewallNames := os.Args[1:]

	interval := 30 * time.Second

	for {
		pinger, err := probing.NewPinger(host)
		if err != nil {
			log.Println("Error creating pinger:", err)
			time.Sleep(interval)
			continue
		}

		pinger.Count = 3
		pinger.Timeout = 5 * time.Second
		err = pinger.Run()
		if err != nil {
			log.Println("Ping failed:", err)
		} else {
			stats := pinger.Statistics()
			if stats.PacketsRecv == 0 {
				// No responses within timeout → treat as failure
				log.Printf("Ping timeout: no reply from %s within %v", host, pinger.Timeout)
				err = fmt.Errorf("ping timeout")
			}
		}
		if err != nil {
			log.Println("Ping failed:", err)
			client := hcloud.NewClient(hcloud.WithToken(token))

			for _, fwName := range firewallNames {
				log.Printf("Processing firewall: %s", fwName)

				fw, err := firewall.GetFirewall(ctx, client, fwName)
				if err != nil {
					log.Printf("Error fetching firewall %s: %s", fwName, err)
					continue
				}

				log.Printf("Firewall ID: %d, Name: %s\n", fw.ID, fw.Name)

				firewall.UpdateFirewallRules(ctx, client, fw)

				log.Printf("Firewall %s updated successfully", fwName)
			}
		} else {
			stats := pinger.Statistics()
			log.Printf("Ping to %s: %v", host, stats.AvgRtt)
			// trigger firewall update if needed
		}

		time.Sleep(interval)
	}

}
