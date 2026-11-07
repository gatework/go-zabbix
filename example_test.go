package zabbix_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gatework/go-zabbix"
)

func ExampleClient() {
	client, err := zabbix.New("https://zabbix.example.com/api_jsonrpc.php",
		zabbix.WithAPIToken("configured-api-token"))
	if err != nil {
		panic(err)
	}
	defer client.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	hosts, err := client.Host().Get(ctx, zabbix.HostGetParams{
		Output: []string{"hostid", "name"}, Limit: 10, SelectInterfaces: "count",
	})
	if err != nil {
		panic(err)
	}
	for _, host := range hosts {
		if host.Interfaces.Count != nil {
			fmt.Println(host.HostID, host.Name, *host.Interfaces.Count)
		}
	}
}

func ExampleClient_Login() {
	client, err := zabbix.New("https://zabbix.example.com/api_jsonrpc.php")
	if err != nil {
		panic(err)
	}
	defer client.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := client.Login(ctx, "configured-user", "configured-password"); err != nil {
		panic(err)
	}
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if err := client.Logout(cleanup); err != nil {
			fmt.Println(err)
		}
	}()
}

func ExampleHostService_Update() {
	client, err := zabbix.New("https://zabbix.example.com/api_jsonrpc.php",
		zabbix.WithAPIToken("configured-api-token"))
	if err != nil {
		panic(err)
	}
	defer client.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	description := ""
	_, err = client.Host().Update(ctx, zabbix.HostUpdateParams{
		HostID: "10084", Description: &description, Tags: []zabbix.TagInput{},
	})
	var apiErr *zabbix.APIError
	if errors.As(err, &apiErr) {
		fmt.Println("Zabbix rejected the update:", apiErr.Code)
	} else if err != nil {
		fmt.Println(err)
	}
}
