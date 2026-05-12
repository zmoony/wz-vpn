# 已知问题 / 风险

- 本地终端缺少可直接调用的 `go` 运行时，后端编译和 Go 测试可能无法在当前环境完成。
- 当前机器大概率不存在真实的 WireGuard 宿主环境，因此 WireGuard 系统级行为只能通过适配层设计和非破坏性命令约定来保证结构正确，无法完成真实联调。
- 前端已完成构建验证，但产物体积较大，Vite 给出 chunk size warning；后续可以考虑路由懒加载和手动分包。
- 当前 WireGuard 客户端配置依赖环境变量 `PI_GATEWAY_WG_PUBLIC_KEY` 提供服务端公钥；若未配置，生成的客户端配置只能作为占位示例，不能直接导入使用。
- 当前 DDNS 适配层写入的是 Pi-Gateway 管理的受管配置文件，并通过 `PI_GATEWAY_DDNSGO_RELOAD_COMMAND` 触发外部 reload；是否与现网 ddns-go 原生配置完全一致，仍需在目标环境核对。
- 反向代理适配层依赖本机存在可执行的 `nginx` 命令，以及目标目录权限正确；本轮未在真实 nginx 环境中验证 `-t` / `reload`。
- DDNS 编辑时前端不会回显历史 `AccessKey Secret`，更新现有配置时需要重新输入 Secret。
