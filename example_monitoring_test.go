package zabbix_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gatework/go-zabbix"
)

func ExampleMetricsService_Traffic() {
	client, err := zabbix.New("https://monitor.example/api_jsonrpc.php", zabbix.WithAPIToken("configured-api-token"))
	if err != nil {
		panic(err)
	}
	defer client.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	till := time.Now().UTC().Truncate(time.Second)
	series, err := client.Metrics().Traffic(ctx, zabbix.TrafficQuery{
		HostID: "10084", ItemID: "42001", From: till.Add(-time.Hour).Unix(), Till: till.Unix(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(series.SampleCount, series.Partial)
}

func ExampleWebClient_Chart() {
	// The configured frontend account uses UTC. SSO/MFA users can instead
	// supply a real CookieJar through WebOptions.HTTPClient.
	webClient, err := zabbix.NewWebClient("https://monitor.example/zabbix/", zabbix.WebOptions{Timezone: time.UTC})
	if err != nil {
		panic(err)
	}
	defer webClient.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := webClient.Login(ctx, "configured-user", "configured-password"); err != nil {
		panic(err)
	}
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if err := webClient.Logout(cleanup); err != nil {
			fmt.Println(err)
		}
	}()
	till := time.Now().UTC().Truncate(time.Second)
	data, err := webClient.Chart(ctx, zabbix.ChartRequest{
		ItemIDs: []string{"42001", "42002"}, From: till.Add(-time.Hour), To: till,
	})
	if err != nil {
		panic(err)
	}
	// Pass data to the application's notification attachment adapter.
	fmt.Println("PNG bytes:", len(data))
}
