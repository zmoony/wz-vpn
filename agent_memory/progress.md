# 当前任务

- 任务：基于现有项目文档生成首轮可运行代码骨架。
- 当前阶段：规划已确认，开始生成实现文件。

# 阶段结论

- 已确认采用“单体后端 + 内嵌前端”方案。
- 已确认本轮真实能力优先模块为登录鉴权、WireGuard 管理。
- 已读取需求、架构、功能清单、技术选型等核心文档。
- 已生成 Go 后端项目骨架：配置、应用装配、路由、中间件、SQLite store、登录鉴权服务、系统状态服务、WireGuard 服务与平台适配层。
- 已生成 Vue 前端项目骨架：登录页、应用布局、系统总览页、WireGuard 管理页，以及其余模块占位页。
- 已完成前端依赖安装与生产构建验证，`web/dist` 已产出。
- 已新增反向代理真实后端能力：`proxy_hosts` 落库、nginx vhost 渲染、`nginx -t` 校验、`nginx -s reload`、失败回滚。
- 已新增 DDNS 真实后端能力：`ddns_configs` 落库、受管 ddns-go 配置文件写入、状态文件读取、reload/sync 命令入口。
- 已把前端 `Proxy`、`DDNS` 页面替换为真实表单与列表页，并接入对应后端 API。
- 已新增证书管理真实后端能力：`cert_configs` 落库、`acme.sh` issue/install/renew 编排、已安装证书到期时间与剩余天数解析。
- 已把前端 `Certs` 页面替换为真实页面，并接入申请、更新、强制续期接口。
- 已新增防火墙真实后端能力：`firewall_rules` 落库、`input` 链规则渲染、规则预览、真实 apply、30 秒待确认状态、确认/自动回滚链路。
- 已把前端 `Firewall` 页面替换为真实页面，并接入规则 CRUD、预览、应用、确认和 pending 状态查询接口。
- 已为防火墙补充最小 `forward` 链能力：支持配置单个内网 CIDR，并与当前 `input + forward` 规则一起预览、应用、确认和回滚。
- 已新增 `conntrack` 只读查看能力：后端通过系统命令读取连接明细，前端在防火墙页面提供当前连接表格和按来源 IP 过滤入口。
- 已把防火墙 `forward` 链升级为独立规则列表：新增 `firewall_forward_rules` 表、后端独立 CRUD/API、nftables 渲染支持显式 forward 规则，并保留旧的单网段配置作为兼容回退。
- 已新增最小安装链路文件：仓库根目录 `.env.example`、`deploy/pi-gateway.service`、`scripts/install.sh`，用于 Linux 宿主机下的目录初始化、文件安装与 systemd 启用。
- 已把“系统设置”页从占位页补成最小可用版：后端新增 `/api/settings` 真实接口，前端新增 `SettingsView`，可保存 WireGuard 关键参数并展示当前运行环境摘要。
- 已修正防火墙页在 `conntrack` 缺失时的降级行为：当前连接区块会显示不可用提示，不再把整个页面体验拖坏。
- 已把 WireGuard 服务接入 settings 表读取：接口名、子网、服务端地址、客户端 DNS、公网 Endpoint、服务端公钥可通过设置页覆盖，并影响后续 Peer 操作。

# 下一步

- 在具备 Go 运行时的环境中执行 `go test ./...` 和 `go build ./cmd/pi-gateway`。
- 在真实 Linux + WireGuard 宿主环境中联调 `wg`、`qrencode`、`/etc/wireguard/wg0.conf` 读写与权限。
- 在真实 Linux + nginx / ddns-go 环境中联调 `PI_GATEWAY_NGINX_SITES_DIR`、`PI_GATEWAY_DDNSGO_*` 相关路径和 reload 命令。
- 在真实 Linux + acme.sh / nginx 环境中联调 `PI_GATEWAY_ACME_SH_PATH`、证书安装目录、`nginx -s reload` 命令和阿里云 DNS-01 凭据。
- 在真实 Linux + nftables 环境中联调 `PI_GATEWAY_NFTABLES_RULES_PATH`、待确认状态文件、备份目录、`nft -f` 应用和超时回滚链路。
- 在具备 `conntrack-tools` 的 Linux 宿主环境中联调 `conntrack -L` 输出解析、来源 IP 过滤效果和大连接量下的页面展示表现。
- 在具备 Go 运行时的环境中执行新增的防火墙服务测试，重点验证独立 `forward` 规则与旧兼容配置的优先级行为。
- 在真实 Linux 宿主环境中验证 `scripts/install.sh`、`deploy/pi-gateway.service`、环境变量模板路径、以及 `systemctl enable/restart` 的实际行为。
- 在浏览器里完成一次人工联调：确认“系统设置”页保存后再次进入页面能看到持久化值，且防火墙页在无 `conntrack` 命令的宿主机上只影响当前连接区块。
- 下一轮可继续补 `forward` 规则的接口字段、NAT、连接聚合视图，以及证书/防火墙的高级配置能力。
- 2026-05-13：反向代理证书选择链路已切到“按根域名引用证书配置”。
  - 后端 `proxy_hosts` 新增 `certificate_root_domain` 字段，并兼容根据旧证书路径反推根域名。
  - `ProxyService` 现在会校验证书根域名是否存在，并从证书配置解析 `fullchain.pem` / `privkey.pem` 路径后再生成 nginx vhost。
  - 前端反向代理页改为下拉选择证书根域名，只读展示解析出的证书和私钥路径；没有可用证书时禁止提交并提示先去证书管理申请。
  - 验证已通过：`go test ./...`、`go build -o pi-gateway.exe ./cmd/pi-gateway`、`npm run build`。
