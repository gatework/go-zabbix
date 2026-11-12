package zabbix_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gatework/go-zabbix"
)

func TestFrontend_InvalidRequestsDoNotReachServer(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	t.Cleanup(server.Close)
	client := frontendClient(t, server, zabbix.FrontendOptions{})
	for _, invalid := range []func(*zabbix.ChartRequest){
		func(r *zabbix.ChartRequest) { r.ItemIDs = nil },
		func(r *zabbix.ChartRequest) { r.GraphID = "200" },
		func(r *zabbix.ChartRequest) { r.ItemIDs = []string{"../index.php"} },
		func(r *zabbix.ChartRequest) { r.Width = 1 },
		func(r *zabbix.ChartRequest) { r.To = r.From.Add(59 * time.Second) },
		func(r *zabbix.ChartRequest) { r.To = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC) },
	} {
		request := chartRequest()
		invalid(&request)
		if _, err := client.Chart(t.Context(), request); !errors.Is(err, zabbix.ErrFrontendInput) {
			t.Fatalf("invalid request = %v", err)
		}
	}
	if calls.Load() != 0 {
		t.Fatal("invalid request performed network I/O")
	}
}

func TestFrontend_ChartCancellation(t *testing.T) {
	t.Parallel()
	started, release := make(chan struct{}), make(chan struct{})
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client := frontendClient(t, server, zabbix.FrontendOptions{})
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := client.Chart(ctx, chartRequest()); done <- err }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("chart did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("chart did not cancel")
	}
}
