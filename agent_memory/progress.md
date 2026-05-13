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

# 下一步

- 在具备 Go 运行时的环境中执行 `go test ./...` 和 `go build ./cmd/pi-gateway`。
- 在真实 Linux + WireGuard 宿主环境中联调 `wg`、`qrencode`、`/etc/wireguard/wg0.conf` 读写与权限。
- 在真实 Linux + nginx / ddns-go 环境中联调 `PI_GATEWAY_NGINX_SITES_DIR`、`PI_GATEWAY_DDNSGO_*` 相关路径和 reload 命令。
- 在真实 Linux + acme.sh / nginx 环境中联调 `PI_GATEWAY_ACME_SH_PATH`、证书安装目录、`nginx -s reload` 命令和阿里云 DNS-01 凭据。
- 在真实 Linux + nftables 环境中联调 `PI_GATEWAY_NFTABLES_RULES_PATH`、待确认状态文件、备份目录、`nft -f` 应用和超时回滚链路。
- 下一轮可继续补 `forward` 链、`conntrack` 查看、以及证书/防火墙的高级配置能力。
