#!/usr/bin/env bash
# Run against a disposable stack; never reuse a development database.
set -euo pipefail

cd "$(dirname "$0")/.."
project="go-zabbix-it-$(date +%s)-$$"
export ZABBIX_API_PORT=0 ZABBIX_SERVER_PORT=0
compose=(docker compose -p "$project")

cleanup() {
  status=$?
  trap - EXIT
  if (( status != 0 )); then
    "${compose[@]}" logs --tail=100 zabbix-server zabbix-web-nginx-pgsql || true
  fi
  "${compose[@]}" down --volumes --remove-orphans || status=1
  exit "$status"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

"${compose[@]}" up -d --pull missing postgresql-server zabbix-server zabbix-web-nginx-pgsql
address=$("${compose[@]}" port zabbix-web-nginx-pgsql 8080)
api="http://${address}/api_jsonrpc.php"
ready=false
for ((attempt=0; attempt<120; attempt++)); do
  response=$(curl --silent --show-error --fail --max-time 3 \
    -H 'Content-Type: application/json-rpc' \
    --data '{"jsonrpc":"2.0","method":"apiinfo.version","params":[],"id":1}' \
    "$api" 2>/dev/null) || response=''
  if [[ "$response" =~ \"result\"[[:space:]]*:[[:space:]]*\"7\.4\. ]]; then
    printf 'Zabbix API: %s\n' "$response"
    ready=true
    break
  fi
  sleep 1
done
if [[ "$ready" != true ]]; then
  printf 'Zabbix 7.4 API did not become ready\n' >&2
  exit 1
fi

for service in postgresql-server zabbix-server zabbix-web-nginx-pgsql; do
  container=$("${compose[@]}" ps -q "$service")
  image=$(docker inspect --format '{{.Image}}' "$container")
  docker image inspect --format '{{json .RepoDigests}} {{.Architecture}}' "$image"
done

GOWORK=off go test -tags=integration -race -count=1 -timeout=10m "$@" ./... \
  -args -zbx.url="$api" -zbx.user=Admin -zbx.pass=zabbix
