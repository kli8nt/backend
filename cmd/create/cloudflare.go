package create

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

// create a cname record

func CreateRecord(client *cloudflare.API, zoneID string, name string, target string, proxied bool) error {
	ctx := context.Background()
	record := cloudflare.CreateDNSRecordParams{
		Content: target,
		Name:    name,
		Type:    "A",
		Comment: "Kli8nt",
		Proxied: &proxied,
		TTL:     3600,
	}

	r, err := client.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), record)
	if err != nil {
		return err
	}
	fmt.Println(r)

	return nil
}

// TODO: delete a canme record
