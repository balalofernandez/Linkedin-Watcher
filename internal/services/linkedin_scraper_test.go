package services

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestScrapeLinkedInConnections(t *testing.T) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	fmt.Println("os.Getenv:", os.Getenv("LINKEDIN_USERNAME"))
	fmt.Println("viper:", viper.GetString("LINKEDIN_USERNAME"))
	fmt.Println("All viper settings:")
	for _, k := range viper.AllKeys() {
		fmt.Printf("%s = %s\n", k, viper.GetString(k))
	}

	ctx := context.Background()
	linkedinID := "ACoAABCTJc0BXeaIs6JV7hwE0oSn0eodv8ab_6I"
	// This test assumes cookies or login is handled elsewhere or not required for public profiles
	connections, err := ScrapeLinkedInConnections(ctx, linkedinID)
	if err != nil {
		t.Fatalf("scrape failed: %v", err)
	}
	if len(connections) == 0 {
		t.Errorf("expected at least one connection, got 0")
	}
}
