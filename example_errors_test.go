package zabbix_test

import (
	"errors"
	"fmt"

	"github.com/gatework/go-zabbix"
)

func ExampleMetricError() {
	err := fmt.Errorf("read traffic: %w", zabbix.MetricErrorUnavailable)
	if metricErr, ok := errors.AsType[zabbix.MetricError](err); ok {
		fmt.Println(metricErr)
	}
	fmt.Println(errors.Is(err, zabbix.MetricErrorUnavailable))
	// Output:
	// zabbix: metric object unavailable
	// true
}

func ExampleAPIError() {
	err := &zabbix.RequestError{
		Method: "host.get",
		Kind:   zabbix.ErrorKindAPI,
		Err:    &zabbix.APIError{Code: -32602, Message: "upstream detail"},
	}
	if apiErr, ok := errors.AsType[*zabbix.APIError](err); ok {
		fmt.Println(apiErr.Code)
	}
	fmt.Println(err)
	// Output:
	// -32602
	// zabbix: host.get: api error -32602
}
