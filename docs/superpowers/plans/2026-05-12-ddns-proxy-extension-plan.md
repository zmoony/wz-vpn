# DDNS And Proxy Extension Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the current Pi-Gateway skeleton with real backend capabilities for DDNS configuration management and reverse proxy configuration management.

**Architecture:** Add two new backend capability chains: `ddns_service -> ddnsgo manager` and `proxy_service -> nginx manager`. Keep API handlers thin, persist user-facing configuration in SQLite, and let platform adapters own file generation plus system command execution.

**Tech Stack:** Go 1.22, Gin, SQLite, Vue 3, Element Plus, nginx CLI integration, ddns-go managed config/status files

---

## File Map

- Modify: `internal/config/config.go`
- Modify: `internal/store/db.go`
- Modify: `internal/app/app.go`
- Modify: `internal/api/router.go`
- Create: `internal/domain/proxy.go`
- Create: `internal/domain/ddns.go`
- Create: `internal/store/proxy_host_store.go`
- Create: `internal/store/ddns_config_store.go`
- Create: `internal/platform/nginx/manager.go`
- Create: `internal/platform/ddnsgo/manager.go`
- Create: `internal/service/proxy_service.go`
- Create: `internal/service/ddns_service.go`
- Create: `internal/service/proxy_service_test.go`
- Create: `internal/service/ddns_service_test.go`
- Create: `internal/api/proxy.go`
- Create: `internal/api/ddns.go`
- Create: `web/src/api/proxy.ts`
- Create: `web/src/api/ddns.ts`
- Create: `web/src/views/ProxyView.vue`
- Create: `web/src/views/DDNSView.vue`
- Modify: `web/src/router/index.ts`
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

### Task 1: Add failing tests and shared persistence changes

**Files:**
- Create: `internal/service/proxy_service_test.go`
- Create: `internal/service/ddns_service_test.go`
- Modify: `internal/store/db.go`

- [ ] Step 1: Add tests for proxy config rendering/apply orchestration and DDNS config save/sync orchestration.
- [ ] Step 2: Add SQLite tables required by the new services.
- [ ] Step 3: Wire matching domain/store types needed by the tests.

### Task 2: Implement backend services and system adapters

**Files:**
- Create: `internal/domain/proxy.go`
- Create: `internal/domain/ddns.go`
- Create: `internal/store/proxy_host_store.go`
- Create: `internal/store/ddns_config_store.go`
- Create: `internal/platform/nginx/manager.go`
- Create: `internal/platform/ddnsgo/manager.go`
- Create: `internal/service/proxy_service.go`
- Create: `internal/service/ddns_service.go`
- Modify: `internal/config/config.go`
- Modify: `internal/app/app.go`

- [ ] Step 1: Implement domain models and store interfaces.
- [ ] Step 2: Implement nginx adapter with render, write, test, reload, and rollback behavior.
- [ ] Step 3: Implement ddns-go adapter with managed config write, status read, and reload command hooks.
- [ ] Step 4: Implement proxy and DDNS services with validation and orchestration.
- [ ] Step 5: Register the services in application bootstrap.

### Task 3: Expose backend APIs and replace frontend stubs

**Files:**
- Create: `internal/api/proxy.go`
- Create: `internal/api/ddns.go`
- Modify: `internal/api/router.go`
- Create: `web/src/api/proxy.ts`
- Create: `web/src/api/ddns.ts`
- Create: `web/src/views/ProxyView.vue`
- Create: `web/src/views/DDNSView.vue`
- Modify: `web/src/router/index.ts`

- [ ] Step 1: Add REST handlers for proxy CRUD and DDNS read/update/sync.
- [ ] Step 2: Replace the two stub routes with real views.
- [ ] Step 3: Hook the views to the new backend endpoints.

### Task 4: Verify what the current environment can prove

**Files:**
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

- [ ] Step 1: Run frontend build to verify the new pages compile.
- [ ] Step 2: Record why Go compile/test and system-level nginx/ddns-go execution remain unverified here.
- [ ] Step 3: Summarize the exact boundaries of what is real versus environment-dependent.
