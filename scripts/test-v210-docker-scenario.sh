#!/usr/bin/env bash
set -euo pipefail

mode="${1:-check}"
release="v2.1.0"
project="${ONEPROXY_V210_PROJECT:-oneproxy-v210-local}"
panel_repo="${ONEPROXY_PANEL_IMAGE_REPO:-ghcr.io/stanleysun233/oneproxy-panel}"
node_repo="${ONEPROXY_NODE_IMAGE_REPO:-ghcr.io/stanleysun233/oneproxy-node}"
tag="${ONEPROXY_IMAGE_TAG:-v2.1.0-rc.local}"
mysql_image="${ONEPROXY_MYSQL_IMAGE:-mysql:8.4}"
redis_image="${ONEPROXY_REDIS_IMAGE:-redis:7}"
guacd_image="${ONEPROXY_GUACD_IMAGE:-guacamole/guacd:1.6.0}"
mysql_db="${ONEPROXY_V210_DB:-oneproxy_v210}"
mysql_root_password="${ONEPROXY_V210_MYSQL_ROOT_PASSWORD:-oneproxy-v210-root}"
jwt_key="${ONEPROXY_V210_JWT_SIGNING_KEY:-oneproxy-v210-jwt}"
admin_password="${ONEPROXY_V210_ADMIN_PASSWORD:-oneproxy-v210-admin}"
tenant_name="${ONEPROXY_V210_TENANT_NAME:-OneProxy v2.1.0 Local}"
node_name="${ONEPROXY_V210_NODE_NAME:-oneproxy-v210-local-node}"
panel_host_port="${ONEPROXY_V210_PANEL_PORT:-12886}"
node_http_host_port="${ONEPROXY_V210_NODE_HTTP_PORT:-12988}"
node_https_host_port="${ONEPROXY_V210_NODE_HTTPS_PORT:-12989}"
node_tcp_host_port="${ONEPROXY_V210_NODE_TCP_PORT:-12990}"
node_udp_host_port="${ONEPROXY_V210_NODE_UDP_PORT:-12991}"
node_direct_host_port="${ONEPROXY_V210_NODE_DIRECT_PORT:-12992}"
mysql_host_port="${ONEPROXY_V210_MYSQL_PORT:-13306}"
redis_host_port="${ONEPROXY_V210_REDIS_PORT:-16379}"
compose_file=""

case "$mode" in
  check|build|run|clean)
    ;;
  *)
    echo "usage: $0 [check|build|run|clean]" >&2
    exit 2
    ;;
esac

cleanup() {
  if [ -n "$compose_file" ]; then
    rm -f "$compose_file"
  fi
}
trap cleanup EXIT

ensure_docker_compose() {
  if docker compose version >/dev/null 2>&1; then
    return 0
  fi
  if command -v docker-compose >/dev/null 2>&1; then
    return 0
  fi
  echo "docker compose is required" >&2
  exit 2
}

compose() {
  if docker compose version >/dev/null 2>&1; then
    docker compose -p "$project" -f "$compose_file" "$@"
  else
    docker-compose -p "$project" -f "$compose_file" "$@"
  fi
}

write_compose() {
  node_bootstrap_token="${1:-}"
  if [ -n "$compose_file" ]; then
    rm -f "$compose_file"
  fi
  compose_file="$(mktemp)"
  cat > "$compose_file" <<COMPOSE
services:
  mysql:
    image: ${mysql_image}
    environment:
      MYSQL_ROOT_PASSWORD: "${mysql_root_password}"
      MYSQL_DATABASE: "${mysql_db}"
    ports:
      - "127.0.0.1:${mysql_host_port}:3306"
    volumes:
      - mysql-data:/var/lib/mysql
  redis:
    image: ${redis_image}
    ports:
      - "127.0.0.1:${redis_host_port}:6379"
    volumes:
      - redis-data:/data
  guacd:
    image: ${guacd_image}
  panel:
    image: ${panel_repo}:${tag}
    depends_on:
      - mysql
      - redis
      - guacd
    environment:
      MYSQL_DSN: "root:${mysql_root_password}@tcp(mysql:3306)/${mysql_db}?charset=utf8mb4&parseTime=true&loc=UTC"
      ADMIN_PASSWORD: "${admin_password}"
      JWT_SIGNING_KEY: "${jwt_key}"
      REDIS_URL: "redis://redis:6379/0"
      GUACD_ADDR: "guacd:4822"
      PORT: "2886"
      HTTP_ADDR: "0.0.0.0:2887"
      CONTROL_PLANE_URL: "http://panel:2887"
      PUBLIC_CERT_PROVIDER: "manual"
      TZ: "UTC"
    ports:
      - "127.0.0.1:${panel_host_port}:2886"
    volumes:
      - panel-data:/app/data
  node:
    image: ${node_repo}:${tag}
    depends_on:
      - panel
    environment:
      NODE_LISTEN_ADDR: ":2988"
      NODE_HTTPS_LISTEN_ADDR: ":2989"
      NODE_TCP_ACCESS_LISTEN_ADDR: ":2990"
      NODE_UDP_ACCESS_LISTEN_ADDR: ":2991"
      NODE_DIRECT_LISTEN_ADDR: ":2992"
      NODE_NAME: "${node_name}"
      NODE_MODE: "edge"
      NODE_BOOTSTRAP_TOKEN: "${node_bootstrap_token}"
      NODE_PARENT_URL: "http://panel:2887"
      NODE_PUBLIC_HOST: "127.0.0.1"
      NODE_HEARTBEAT_INTERVAL: "2s"
      NODE_DIRECT_STUN_SERVERS: ""
      NODE_DIRECT_REFRESH_INTERVAL: "5s"
      NODE_POLICY_STATE_PATH: "runtime/node-policy-state.json"
      NODE_RUNTIME_CONFIG_PATH: "runtime/node-runtime.json"
      PUBLIC_CERT_PROVIDER: "manual"
      TZ: "UTC"
    ports:
      - "127.0.0.1:${node_http_host_port}:2988"
      - "127.0.0.1:${node_https_host_port}:2989"
      - "127.0.0.1:${node_tcp_host_port}:2990"
      - "127.0.0.1:${node_udp_host_port}:2991/udp"
      - "127.0.0.1:${node_direct_host_port}:2992/udp"
    volumes:
      - node-runtime:/app/runtime
volumes:
  mysql-data:
  redis-data:
  panel-data:
  node-runtime:
COMPOSE
}

print_plan() {
  echo "release=${release}"
  echo "mode=${mode}"
  echo "project=${project}"
  echo "services=mysql,redis,guacd,panel,node"
  echo "panel_image=${panel_repo}:${tag}"
  echo "node_image=${node_repo}:${tag}"
  echo "mysql_image=${mysql_image}"
  echo "redis_image=${redis_image}"
  echo "guacd_image=${guacd_image}"
  echo "tenant_name=${tenant_name}"
  echo "node_name=${node_name}"
  echo "host_ports=panel:${panel_host_port},node_http:${node_http_host_port},node_https:${node_https_host_port},node_tcp:${node_tcp_host_port},node_udp:${node_udp_host_port},node_direct:${node_direct_host_port},mysql:${mysql_host_port},redis:${redis_host_port}"
  echo "destructive_modes=build,run,clean"
}

wait_for_http() {
  url="$1"
  label="$2"
  tries=0
  while :; do
    if curl -fsS "$url" >/dev/null 2>&1; then
      echo "${label}=ok"
      return 0
    fi
    tries=$((tries + 1))
    if [ "$tries" -ge 60 ]; then
      echo "${label}=failed" >&2
      return 1
    fi
    sleep 2
  done
}

wait_for_node_bound() {
  tries=0
  while :; do
    body="$(curl -fsS "http://127.0.0.1:${node_http_host_port}/healthz" 2>/dev/null || true)"
    if printf '%s' "$body" | grep -q '"controlPlaneBound":true'; then
      echo "node_bound=ok"
      return 0
    fi
    tries=$((tries + 1))
    if [ "$tries" -ge 60 ]; then
      echo "node_bound=failed body=${body}" >&2
      return 1
    fi
    sleep 2
  done
}

json_get() {
  path="$1"
  node -e 'let raw="";process.stdin.on("data",chunk=>raw+=chunk).on("end",()=>{let value=JSON.parse(raw);for(const part of process.argv[1].split(".").filter(Boolean)){value=Array.isArray(value)?value[Number(part)]:value?.[part];if(value==null)break}if(value!=null)process.stdout.write(typeof value==="object"?JSON.stringify(value):String(value))})' "$path"
}

sha256_hex() {
  node -e 'process.stdout.write(require("node:crypto").createHash("sha256").update(process.argv[1]).digest("hex"))' "$1"
}

api_request() {
  method="$1"
  path="$2"
  access_token="${3:-}"
  tenant_id="${4:-}"
  body="${5:-}"
  headers=(-H "Content-Type: application/json")
  if [ -n "$access_token" ]; then
    headers+=(-H "X-One-Proxy-Access-Token: ${access_token}")
  fi
  if [ -n "$tenant_id" ]; then
    headers+=(-H "X-One-Proxy-Tenant-ID: ${tenant_id}")
  fi
  if [ -n "$body" ]; then
    curl -fsS -X "$method" "http://127.0.0.1:${panel_host_port}/api${path}" "${headers[@]}" -d "$body"
  else
    curl -fsS -X "$method" "http://127.0.0.1:${panel_host_port}/api${path}" "${headers[@]}"
  fi
}

node_api_request() {
  method="$1"
  path="$2"
  node_token="$3"
  body="${4:-}"
  if [ -n "$body" ]; then
    curl -fsS -X "$method" "http://127.0.0.1:${panel_host_port}/api${path}" -H "Content-Type: application/json" -H "X-One-Proxy-Node-Token: ${node_token}" -d "$body"
  else
    curl -fsS -X "$method" "http://127.0.0.1:${panel_host_port}/api${path}" -H "X-One-Proxy-Node-Token: ${node_token}"
  fi
}

create_control_plane_state() {
  login_response="$(api_request POST /auth/login "" "" "{\"account\":\"admin\",\"password\":\"${admin_password}\"}")"
  access_token="$(printf '%s' "$login_response" | json_get data.accessToken)"
  account_id="$(printf '%s' "$login_response" | json_get data.account.id)"

  tenant_response="$(api_request POST /tenants "$access_token" "" "{\"name\":\"${tenant_name}\",\"initialAdminAccountId\":\"${account_id}\"}")"
  tenant_id="$(printf '%s' "$tenant_response" | json_get data.tenant.id)"

  scope_response="$(api_request POST /proxy/scopes "$access_token" "$tenant_id" "{\"name\":\"v2.1.0 local scope\",\"description\":\"isolated release scenario\"}")"
  scope_id="$(printf '%s' "$scope_response" | json_get data.id)"

  bootstrap_response="$(api_request POST /nodes/bootstrap/token "$access_token" "$tenant_id" "{\"targetType\":\"node\",\"nodeName\":\"${node_name}\",\"nodeMode\":\"edge\",\"scopeKey\":\"${scope_id}\",\"publicHost\":\"127.0.0.1\",\"publicPort\":${node_http_host_port}}")"
  node_bootstrap_token="$(printf '%s' "$bootstrap_response" | json_get data.token)"
  echo "tenant_created=${tenant_id}"
  echo "scope_created=${scope_id}"
  echo "bootstrap_token_issued=yes"
}

wait_for_pending_node() {
  tries=0
  while :; do
    pending_response="$(api_request GET /nodes/pending "$access_token" "$tenant_id")"
    node_id="$(printf '%s' "$pending_response" | json_get data.0.id)"
    if [ -n "$node_id" ]; then
      echo "node_pending=${node_id}"
      return 0
    fi
    tries=$((tries + 1))
    if [ "$tries" -ge 60 ]; then
      echo "node_pending=failed" >&2
      return 1
    fi
    sleep 2
  done
}

create_route_state() {
  chain_response="$(api_request POST /proxy "$access_token" "$tenant_id" "{\"name\":\"v2.1.0 local chain\",\"destinationScope\":\"${scope_id}\",\"hopGroups\":[{\"candidates\":[\"${node_id}\"]}]}")"
  chain_id="$(printf '%s' "$chain_response" | json_get data.id)"

  http_path_response="$(api_request POST /proxy/paths "$access_token" "$tenant_id" "{\"chainId\":\"${chain_id}\",\"name\":\"HTTP access\",\"mode\":\"forward\",\"protocol\":\"http\",\"serviceType\":\"http_forward_proxy\",\"remoteProtocol\":\"\",\"listenHost\":\"127.0.0.1\",\"listenPort\":${node_http_host_port},\"targetProtocol\":\"http\",\"targetHost\":\"example.test\",\"targetPort\":80,\"targetSni\":\"\",\"tlsMode\":\"off\",\"authMode\":\"proxy_token\",\"options\":{}}")"
  http_path_id="$(printf '%s' "$http_path_response" | json_get data.id)"
  ssh_path_response="$(api_request POST /proxy/paths "$access_token" "$tenant_id" "{\"chainId\":\"${chain_id}\",\"name\":\"SSH access\",\"mode\":\"tcp\",\"protocol\":\"tcp\",\"serviceType\":\"tcp_access\",\"remoteProtocol\":\"ssh\",\"listenHost\":\"127.0.0.1\",\"listenPort\":${node_tcp_host_port},\"targetProtocol\":\"ssh\",\"targetHost\":\"ssh.example.test\",\"targetPort\":22,\"targetSni\":\"\",\"tlsMode\":\"off\",\"authMode\":\"proxy_token\",\"options\":{}}")"
  ssh_path_id="$(printf '%s' "$ssh_path_response" | json_get data.id)"
  rdp_path_response="$(api_request POST /proxy/paths "$access_token" "$tenant_id" "{\"chainId\":\"${chain_id}\",\"name\":\"RDP access\",\"mode\":\"tcp\",\"protocol\":\"tcp\",\"serviceType\":\"tcp_access\",\"remoteProtocol\":\"rdp\",\"listenHost\":\"127.0.0.1\",\"listenPort\":${node_tcp_host_port},\"targetProtocol\":\"rdp\",\"targetHost\":\"rdp.example.test\",\"targetPort\":3389,\"targetSni\":\"\",\"tlsMode\":\"off\",\"authMode\":\"proxy_token\",\"options\":{}}")"
  rdp_path_id="$(printf '%s' "$rdp_path_response" | json_get data.id)"

  api_request PUT /remote/defaults/ssh "$access_token" "$tenant_id" "{\"accessPathId\":\"${ssh_path_id}\"}" >/dev/null
  api_request PUT /remote/defaults/rdp "$access_token" "$tenant_id" "{\"accessPathId\":\"${rdp_path_id}\"}" >/dev/null

  group_response="$(api_request POST /proxy/route-groups "$access_token" "$tenant_id" "{\"name\":\"v2.1.0 local routes\",\"description\":\"isolated release scenario\"}")"
  route_group_id="$(printf '%s' "$group_response" | json_get data.id)"

  route_body="{\"groupId\":\"${route_group_id}\",\"priority\":10,\"matchType\":\"domain_suffix\",\"matchValue\":\".example.test\",\"actionType\":\"chain\",\"chainId\":\"${chain_id}\",\"destinationScope\":\"\"}"
  route_response="$(api_request POST /proxy/routes "$access_token" "$tenant_id" "$route_body")"
  route_id="$(printf '%s' "$route_response" | json_get data.id)"
  echo "chain_created=${chain_id}"
  echo "access_paths_created=http:${http_path_id},ssh:${ssh_path_id},rdp:${rdp_path_id}"
  echo "route_created=${route_id}"
}

publish_policy() {
  publish_response="$(api_request POST /policies/publish "$access_token" "$tenant_id" "{}")"
  policy_revision="$(printf '%s' "$publish_response" | json_get data.version)"
  if [ -z "$policy_revision" ]; then
    echo "policy_publish=failed response=${publish_response}" >&2
    return 1
  fi
  echo "policy_published=${policy_revision}"
}

validate_latest_bootstrap() {
  bootstrap_response="$(api_request GET /proxy/extension/bootstrap "$access_token" "$tenant_id")"
  printf '%s' "$bootstrap_response" | node -e 'let raw="";process.stdin.on("data",chunk=>raw+=chunk).on("end",()=>{const payload=JSON.parse(raw).data;if(payload.schemaVersion!=="v2.1.0"||"groups" in payload||!payload.nodes?.length||!payload.accessPaths?.length||!payload.routes?.length||!payload.routes[0].accessPathId)process.exit(1);if(!payload.accessPaths.every(path=>Array.isArray(path.entrypoints)&&Array.isArray(path.topologyGroups)))process.exit(1);console.log(`bootstrap_schema=v2.1.0 access_paths=${payload.accessPaths.length} routes=${payload.routes.length}`)})'
  proxy_token="$(printf '%s' "$bootstrap_response" | json_get data.proxyToken)"
  bootstrap_access_path_id="$(printf '%s' "$bootstrap_response" | json_get data.accessPaths.0.id)"
  bootstrap_route_id="$(printf '%s' "$bootstrap_response" | json_get data.routes.0.id)"
  bootstrap_route_access_path_id="$(printf '%s' "$bootstrap_response" | json_get data.routes.0.accessPathId)"
  if [ "$bootstrap_access_path_id" != "$bootstrap_route_access_path_id" ]; then
    echo "bootstrap_access_path_mismatch path=${bootstrap_access_path_id} route=${bootstrap_route_access_path_id}" >&2
    return 1
  fi
  echo "bootstrap_route_access_path=ok"

  defaults_response="$(api_request GET /remote/defaults "$access_token" "$tenant_id")"
  printf '%s' "$defaults_response" | node -e 'let raw="";process.stdin.on("data",chunk=>raw+=chunk).on("end",()=>{const items=JSON.parse(raw).data;const ssh=items.find(item=>item.protocol==="ssh");const rdp=items.find(item=>item.protocol==="rdp");if(!ssh||!rdp||ssh.accessPathId===rdp.accessPathId)process.exit(1);console.log(`remote_defaults=ssh:${ssh.accessPathId},rdp:${rdp.accessPathId}`)})'
}

validate_node_policy() {
  policy_response="$(node_api_request GET /node/agent/policy "$node_access_token")"
  printf '%s' "$policy_response" | node -e 'let raw="";process.stdin.on("data",chunk=>raw+=chunk).on("end",()=>{const data=JSON.parse(raw).data;const payload=JSON.parse(data.payloadJson);const snapshots=payload.snapshots||[];const routeRules=snapshots[0]?.payload?.routeRules||[];const chains=snapshots[0]?.payload?.chains||[];if(!snapshots.length||!routeRules.length||!routeRules[0].accessPathId||!chains.length||!Array.isArray(chains[0].hopGroups)||"hops" in chains[0])process.exit(1);console.log(`node_policy_revision=${data.policyRevisionId} access_path_id=${routeRules[0].accessPathId} routes=${routeRules.length}`)})'
}

validate_proxy_token() {
  token_hash="$(sha256_hex "$proxy_token")"
  validate_body="{\"tokenHash\":\"${token_hash}\",\"accessPathId\":\"${bootstrap_access_path_id}\",\"targetHost\":\"example.test\",\"targetPort\":80,\"protocol\":\"http\",\"routeId\":\"${bootstrap_route_id}\"}"
  node_api_request POST /node/agent/proxy/token/validate "$node_access_token" "$validate_body" >/dev/null
  echo "proxy_token_validation=ok"
}

run_db_evidence() {
  mysql_query() {
    compose exec -T mysql env MYSQL_PWD="${mysql_root_password}" mysql -uroot "${mysql_db}" -N -B -e "$1" || true
  }
  mysql_query "SELECT id, status, enabled, public_host, public_port FROM nodes ORDER BY id;"
  mysql_query "SELECT node_id, transport_type, direction, address, status FROM node_transports ORDER BY node_id, transport_type;"
  mysql_query "SELECT id, chain_id, remote_protocol, listen_port, target_host, target_port, enabled FROM node_access_paths ORDER BY id;"
  mysql_query "SELECT chain_id, hop_index, candidate_index, node_id FROM chain_hops ORDER BY chain_id, hop_index, candidate_index;"
  mysql_query "SELECT tenant_id, remote_protocol, access_path_id FROM tenant_remote_access_defaults ORDER BY tenant_id, remote_protocol;"
  mysql_query "SELECT id, priority, match_type, match_value, action_type, chain_id, destination_scope, enabled FROM route_rules ORDER BY priority, id;"
  mysql_query "SELECT id, access_token_hash REGEXP '^[0-9a-f]{64}$' AS access_hash_shape, refresh_token_hash REGEXP '^[0-9a-f]{64}$' AS refresh_hash_shape FROM sessions ORDER BY id;"
  mysql_query "SELECT id, token_hash REGEXP '^[0-9a-f]{64}$' AS token_hash_shape FROM node_api_tokens ORDER BY node_id, id;"
  mysql_query "SELECT id, token_hash REGEXP '^[0-9a-f]{64}$' AS token_hash_shape, consumed_at IS NOT NULL AS consumed FROM bootstrap_tokens ORDER BY id;"
}

case "$mode" in
  check)
    print_plan
    if command -v docker >/dev/null 2>&1; then
      echo "docker=found"
    else
      echo "docker=missing"
    fi
    if docker compose version >/dev/null 2>&1 || command -v docker-compose >/dev/null 2>&1; then
      echo "docker_compose=found"
    else
      echo "docker_compose=missing"
    fi
    command -v curl >/dev/null 2>&1 && echo "curl=found" || echo "curl=missing"
    command -v node >/dev/null 2>&1 && echo "node=found" || echo "node=missing"
    ;;
  build)
    ensure_docker_compose
    print_plan
    docker build -f docker/one-proxy-panel-base.Dockerfile -t oneproxy-panel-base:"$tag" .
    docker build -f docker/one-proxy-panel.Dockerfile --build-arg PANEL_BASE_IMAGE=oneproxy-panel-base:"$tag" -t "${panel_repo}:${tag}" .
    docker build -f docker/one-proxy-node-base.Dockerfile -t oneproxy-node-base:"$tag" .
    docker build -f docker/one-proxy-node.Dockerfile --build-arg NODE_BASE_IMAGE=oneproxy-node-base:"$tag" -t "${node_repo}:${tag}" .
    ;;
  run)
    ensure_docker_compose
    command -v curl >/dev/null 2>&1 || {
      echo "curl is required for run mode" >&2
      exit 2
    }
    command -v node >/dev/null 2>&1 || {
      echo "node is required for run mode" >&2
      exit 2
    }
    write_compose ""
    print_plan
    compose down --volumes --remove-orphans
    compose up -d mysql redis panel
    wait_for_http "http://127.0.0.1:${panel_host_port}/healthz" "panel_health"
    create_control_plane_state
    write_compose "$node_bootstrap_token"
    compose up -d node
    wait_for_pending_node
    approve_response="$(api_request POST "/nodes/${node_id}/approve" "$access_token" "$tenant_id" "{}")"
    node_access_token="$(printf '%s' "$approve_response" | json_get data.accessToken)"
    echo "node_approved=${node_id}"
    wait_for_node_bound
    create_route_state
    publish_policy
    validate_latest_bootstrap
    validate_node_policy
    validate_proxy_token
    run_db_evidence
    ;;
  clean)
    ensure_docker_compose
    write_compose ""
    print_plan
    compose down --volumes --remove-orphans
    ;;
esac
