package zabbix_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gatework/go-zabbix"
)

func ExampleMonitoringService_Traffic() {
	client, err := zabbix.New("https://monitor.example/api_jsonrpc.php", zabbix.WithAPIToken("configured-api-token"))
	if err != nil {
		panic(err)
	}
	defer client.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	till := time.Now().UTC().Truncate(time.Second)
	series, err := client.Monitoring().Traffic(ctx, zabbix.TrafficQuery{
		HostID: "10084", ItemID: "42001", From: till.Add(-time.Hour).Unix(), Till: till.Unix(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(series.SampleCount, series.Partial)
}

func ExampleFrontend_Chart() {
	// The configured frontend account uses UTC. SSO/MFA users can instead
	// supply a real CookieJar through FrontendOptions.HTTPClient.
	frontend, err := zabbix.NewFrontend("https://monitor.example/zabbix/", zabbix.FrontendOptions{Timezone: time.UTC})
	if err != nil {
		panic(err)
	}
	defer frontend.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := frontend.Login(ctx, "configured-user", "configured-password"); err != nil {
		panic(err)
	}
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if err := frontend.Logout(cleanup); err != nil {
			fmt.Println(err)
		}
	}()
	till := time.Now().UTC().Truncate(time.Second)
	data, err := frontend.Chart(ctx, zabbix.ChartRequest{
		ItemIDs: []string{"42001", "42002"}, From: till.Add(-time.Hour), To: till,
	})
	if err != nil {
		panic(err)
	}
	// Pass data to the application's notification attachment adapter.
	fmt.Println("PNG bytes:", len(data))
}
