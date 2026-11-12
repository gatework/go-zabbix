package zabbix_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gatework/go-zabbix"
)

func frontendCookie() string {
	// A fixture envelope, not a real signature. The server owns signature checking.
	return url.PathEscape(base64.StdEncoding.EncodeToString([]byte(`{"sessionid":"0123456789abcdef0123456789abcdef","sign":"fixture-signature"}`)))
}

func chartPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	var data bytes.Buffer
	if err := png.Encode(&data, image.NewNRGBA(image.Rect(0, 0, width, height))); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

func frontendClient(t *testing.T, server *httptest.Server, options zabbix.FrontendOptions) *zabbix.Frontend {
	t.Helper()
	hc := server.Client()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	base, err := url.Parse(server.URL + "/zabbix/")
	if err != nil {
		t.Fatal(err)
	}
	jar.SetCookies(base, []*http.Cookie{{Name: "zbx_session", Value: frontendCookie(), Path: "/zabbix/"}})
	hc.Jar = jar
	options.HTTPClient = hc
	client, err := zabbix.NewFrontend(base.String(), options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.CloseIdleConnections)
	return client
}

func chartRequest() zabbix.ChartRequest {
	from := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)
	return zabbix.ChartRequest{ItemIDs: []string{"101", "102"}, From: from, To: from.Add(time.Hour)}
}

func TestFrontend_LoginChartLogout(t *testing.T) {
	t.Parallel()
	picture := chartPNG(t, 940, 280)
	var charts atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/zabbix/index.php":
			if r.Method != http.MethodPost || r.ParseForm() != nil || r.Form.Get("name") != "monitor" || r.Form.Get("password") != "a+b&c" || r.Form.Get("enter") == "" {
				t.Error("invalid login form")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "zbx_session", Value: frontendCookie(), Path: "/zabbix/", Secure: true, HttpOnly: true})
			http.Redirect(w, r, "/zabbix/zabbix.php?action=dashboard.view", http.StatusSeeOther)
		case "/zabbix/chart.php", "/zabbix/chart2.php":
			assertChartRequest(t, r)
			charts.Add(1)
			w.Header().Set("Content-Type", "image/png")
			if _, err := w.Write(picture); err != nil {
				t.Error(err)
			}
		case "/zabbix/api_jsonrpc.php":
			var request struct {
				Method string `json:"method"`
				ID     int64  `json:"id"`
			}
			if json.NewDecoder(r.Body).Decode(&request) != nil || request.Method != "user.logout" || r.Header.Get("Authorization") != "Bearer 0123456789abcdef0123456789abcdef" {
				t.Error("logout did not use the actual frontend session")
			}
			if err := json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": request.ID, "result": true}); err != nil {
				t.Error(err)
			}
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	client, err := zabbix.NewFrontend(server.URL+"/zabbix", zabbix.FrontendOptions{HTTPClient: server.Client(), Timezone: time.FixedZone("UTC+8", 8*3600)})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Login(t.Context(), "monitor", "a+b&c"); err != nil {
		t.Fatal(err)
	}
	request := chartRequest()
	for _, graphID := range []string{"", "200"} {
		request.GraphID = graphID
		if graphID != "" {
			request.ItemIDs = nil
		}
		data, err := client.Chart(t.Context(), request)
		if err != nil || !bytes.Equal(data, picture) {
			t.Fatalf("chart = %d bytes, %v", len(data), err)
		}
	}
	if err := client.Logout(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Chart(t.Context(), request); !errors.Is(err, zabbix.ErrFrontendAuthentication) || charts.Load() != 2 {
		t.Fatalf("logged-out chart = %v; requests=%d", err, charts.Load())
	}
}

func assertChartRequest(t *testing.T, r *http.Request) {
	t.Helper()
	cookie, err := r.Cookie("zbx_session")
	if err != nil || cookie.Value != frontendCookie() || r.Header.Get("Authorization") != "" {
		t.Error("chart must use its cookie without an API bearer token")
	}
	q := r.URL.Query()
	if q.Get("from") != "2026-11-12 08:00:00" || q.Get("to") != "2026-11-12 09:00:00" || q.Get("width") != "900" || q.Get("height") != "240" || q.Get("resolve_macros") != "1" {
		t.Error("chart lost its fixed window, timezone or dimensions")
	}
	if strings.HasSuffix(r.URL.Path, "chart2.php") {
		if q.Get("graphid") != "200" || q.Get("profileIdx") != "web.charts.filter" || len(q["itemids[]"]) != 0 {
			t.Error("wrong configured graph query")
		}
	} else if q.Get("profileIdx") != "web.item.graph.filter" || q.Get("batch") != "1" || strings.Join(q["itemids[]"], ",") != "101,102" {
		t.Error("wrong item graph query")
	}
}

func TestFrontend_LoginRejectsRedirectsAndForms(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, location string
		status         int
		want           error
	}{
		{"login form", "", 200, zabbix.ErrFrontendAuthentication},
		{"external redirect", "https://outside.invalid/steal", 302, zabbix.ErrFrontendAuthentication},
		{"MFA", "/zabbix/index_mfa.php", 302, zabbix.ErrFrontendInteractiveLogin},
		{"wrong directory", "/another/dashboard", 302, zabbix.ErrFrontendAuthentication},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.SetCookie(w, &http.Cookie{Name: "zbx_session", Value: frontendCookie(), Path: "/zabbix/"})
				w.Header().Set("Location", tc.location)
				w.WriteHeader(tc.status)
			}))
			t.Cleanup(server.Close)
			client, err := zabbix.NewFrontend(server.URL+"/zabbix/", zabbix.FrontendOptions{HTTPClient: server.Client()})
			if err != nil {
				t.Fatal(err)
			}
			if err := client.Login(t.Context(), "monitor", "private-secret"); !errors.Is(err, tc.want) || strings.Contains(err.Error(), "private-secret") {
				t.Fatalf("login error = %v", err)
			}
		})
	}
}

func TestFrontend_RejectsNonCharts(t *testing.T) {
	t.Parallel()
	valid := chartPNG(t, 940, 280)
	for _, tc := range []struct {
		name, mediaType string
		data            []byte
		limit           int64
	}{
		{"HTML", "text/html", []byte("<html>login</html>"), 0},
		{"error PNG", "image/png", chartPNG(t, 600, 20), 0},
		{"truncated", "image/png", valid[:len(valid)/2], 0},
		{"appended error", "image/png", append(bytes.Clone(valid), []byte("<html>error</html>")...), 0},
		{"response cap", "image/png", valid, 10},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.mediaType)
				if _, err := w.Write(tc.data); err != nil {
					t.Error(err)
				}
			}))
			t.Cleanup(server.Close)
			client := frontendClient(t, server, zabbix.FrontendOptions{MaxImageBytes: tc.limit})
			data, err := client.Chart(t.Context(), chartRequest())
			if err == nil || data != nil {
				t.Fatalf("invalid image accepted: %d bytes, %v", len(data), err)
			}
		})
	}
}

func TestFrontend_ZeroValueAndCancellation(t *testing.T) {
	t.Parallel()
	var zero zabbix.Frontend
	if _, err := zero.Chart(t.Context(), chartRequest()); !errors.Is(err, zabbix.ErrFrontendInput) {
		t.Fatalf("zero frontend = %v", err)
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Error("cancelled request sent") }))
	t.Cleanup(server.Close)
	client := frontendClient(t, server, zabbix.FrontendOptions{})
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := client.Chart(ctx, chartRequest()); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation = %v", err)
	}
}
