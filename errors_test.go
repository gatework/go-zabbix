package zabbix

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTypedErrorsSupportWrappingAndInspection(t *testing.T) {
	t.Parallel()
	checkErrorCategory(t, AuthenticationErrorUnknown, AuthenticationErrorStateChanged, AuthenticationErrorSessionRequired)
	checkErrorCategory(t, MetricErrorUnknown, MetricErrorInvalidQuery, MetricErrorUnavailable, MetricErrorInvalidData)
	checkErrorCategory(t, WebErrorUnknown, WebErrorInvalidRequest, WebErrorAuthenticationRequired,
		WebErrorInteractiveLoginRequired, WebErrorInvalidChartImage)
	if errors.Is(WebErrorInvalidRequest, MetricError(WebErrorInvalidRequest)) {
		t.Fatal("equal messages from different error categories matched")
	}
}

func checkErrorCategory[E interface {
	error
	comparable
}](t *testing.T, values ...E) {
	t.Helper()
	for _, want := range values {
		t.Run(want.Error(), func(t *testing.T) {
			t.Parallel()
			cause := errors.Join(context.Canceled, fmt.Errorf("operation: %w", want))
			err := requestError("test.operation", ErrorKindTransport, cause)
			got, ok := errors.AsType[E](err)
			if !ok || got != want || !errors.Is(err, want) || !errors.Is(err, context.Canceled) {
				t.Fatalf("error category or joined cause lost: %v", err)
			}
			if want.Error() != strings.ToLower(want.Error()) || !strings.HasPrefix(want.Error(), "zabbix: ") {
				t.Fatalf("inconsistent error text: %q", want.Error())
			}
		})
	}
}

func TestValidationErrorsAreTypedAndHideInputValues(t *testing.T) {
	t.Parallel()
	for _, invoke := range []func() error{
		func() error { _, err := New("https://user:password-secret@example.invalid"); return err },
		func() error { _, err := New("https://example.invalid", WithTimeout(-time.Second)); return err },
		func() error {
			client, err := New("https://example.invalid")
			if err != nil {
				return err
			}
			defer client.CloseIdleConnections()
			_, err = client.Host().Create(t.Context(), HostCreateParams{})
			return err
		},
	} {
		err := invoke()
		validation, ok := errors.AsType[ValidationError](err)
		if !ok || strings.Contains(err.Error(), "secret") || strings.Contains(validation.Error(), "secret") {
			t.Fatalf("missing safe validation error: %v", err)
		}
	}
	var request RequestError
	if request.Kind != ErrorKindUnknown || request.Error() != "zabbix: unknown request error" {
		t.Fatalf("uninitialized request kind has misleading text: %v", &request)
	}
}
