package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
)

var providerTokens map[string]string = map[string]string{}

func cloudflareUpdateRecord(ctx *gin.Context, zoneId string, recordId string, value string) error {
	token, err := getProviderToken("cloudflare")
	if err != nil {
		return err
	}

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return err
	}

	_, err = api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneId), cloudflare.UpdateDNSRecordParams{
		ID:      recordId,
		Type:    "A",
		Content: value,
	})
	if err != nil {
		return err
	}

	return nil
}

func getProviderToken(provider string) (string, error) {
	key := fmt.Sprintf("%s_API_TOKEN", strings.ToUpper(provider))

	token, ok := providerTokens[key]
	if ok {
		return token, nil
	}

	token = os.Getenv(key)
	if token == "" {
		return "", fmt.Errorf("Invalid configuration for '%s' provider", provider)
	}

	providerTokens[key] = token

	return token, nil
}
