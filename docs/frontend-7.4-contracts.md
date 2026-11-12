# Zabbix 7.4 frontend chart contract

This separate web adapter is based on the official **7.4.16** frontend source.
It is not a JSON-RPC image API. `graph.get` returns graph configuration;
`image.get` returns stored icons/backgrounds, not rendered traffic history.

| Capability | Request | Source |
| --- | --- | --- |
| Login | POST `index.php`: `name`, `password`, `enter`, `form`, `autologin=0` | [index.php](https://github.com/zabbix/zabbix/blob/7.4.16/ui/index.php) |
| Item chart | GET `chart.php`: repeated `itemids[]`, `batch=1`, `profileIdx=web.item.graph.filter` | [chart.php](https://github.com/zabbix/zabbix/blob/7.4.16/ui/chart.php) |
| Configured graph | GET `chart2.php`: `graphid`, `profileIdx=web.charts.filter` | [chart2.php](https://github.com/zabbix/zabbix/blob/7.4.16/ui/chart2.php) |
| Time window | `from`, `to` in `YYYY-MM-DD HH:MM:SS`; explicit `profileIdx` is mandatory | [getTimeSelectorPeriod](https://github.com/zabbix/zabbix/blob/7.4.16/ui/include/func.inc.php) |
| Cookie | Server-issued signed `zbx_session` envelope | [CEncryptedCookieSession](https://github.com/zabbix/zabbix/blob/7.4.16/ui/include/classes/core/CEncryptedCookieSession.php), [CWebUser](https://github.com/zabbix/zabbix/blob/7.4.16/ui/include/classes/user/CWebUser.php) |
| Logout | `user.logout` with the underlying session ID from the actual cookie | [CUser](https://github.com/zabbix/zabbix/blob/7.4.16/ui/include/classes/api/services/CUser.php) |

Without `profileIdx`, 7.4's `getTimeSelectorPeriod` replaces supplied `from/to`
with the default period. `Frontend.Chart` supplies the correct index automatically.
It uses second precision and an explicit time window of at least one minute;
it does not use legacy `stime`/`period` parameters. `FrontendOptions.Timezone`
must match the frontend user's timezone. UTC on both sides avoids ambiguity
during daylight-saving clock changes. Nil options timezone means UTC.

The frontend base URL is the directory, including any reverse-proxy prefix,
for example `https://monitor.example/zabbix/`. It is distinct from the API URL.
Login redirects are inspected but never followed. Cross-origin redirects and
redirects outside this directory are rejected. A login form returned with HTTP
200 is an authentication failure; an MFA/SSO redirect requires an external flow.
Reverse-proxy authentication and custom login pages are outside the form contract.

`NewFrontend` copies cookies from a supplied `http.Client.Jar` for the base URL
into a private jar and borrows the transport. Later changes do not affect the
caller's jar. No cookie or password is exposed through a credential getter.
The server verifies the cookie signature; the client only reads the envelope
to identify the session for logout. A raw API token or `user.login` session ID
cannot be converted into the signed frontend cookie.

Login operations are serialized and session changes invalidate pending attempts.
Failed login preserves the previous session; replacement does not revoke it,
so call `Logout` first when revocation is needed. A superseded successful login
gets a bounded compensating logout. `Logout` only clears the matching local
generation after successful server revocation. `ClearSession` invalidates pending
logins but does not revoke the server session or cancel an already running chart.

The default HTTP timeout is 30 seconds (a positive borrowed timeout is retained).
The PNG cap is 8 MiB by default, configurable up to 64 MiB. Requested plot size
defaults to 900x240, with width 64–4096 and height 64–2048; labels and legends may
increase the actual image. The decoder checks MIME, byte cap, dimensions, a
16-megapixel allocation bound, PNG integrity and absence of trailing output.

Zabbix can return **HTTP 200 with an error PNG**, as shown by `show_messages`
in [func.inc.php](https://github.com/zabbix/zabbix/blob/7.4.16/ui/include/func.inc.php).
Small error images, HTML and malformed PNGs are rejected. These checks cannot
prove that every full-size PNG contains the intended samples; an empty graph is
a valid image. Validate item bindings and data separately when alert semantics
require that proof. Frontend version/customization changes require a real smoke
test before use.

HTTP fixtures cover subdirectory login, escaped passwords, signed-cookie transfer,
both chart endpoints, fixed timezone/window, logout, unsafe redirects, rejected
images and concurrent session cleanup. `make integration` includes a disposable
sample push/read-back and frontend login/chart/logout round trip. Compilation or mocked HTTP tests alone
are not evidence that a deployed frontend accepts the adapter.
