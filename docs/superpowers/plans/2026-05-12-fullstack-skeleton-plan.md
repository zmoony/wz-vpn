# Fullstack Skeleton Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go + Vue single-repo skeleton for Pi-Gateway with working admin auth and WireGuard peer management, while keeping DDNS, cert, proxy, and firewall as expandable module shells.

**Architecture:** The backend is a Gin-based API server with SQLite persistence and a platform adapter layer for system integrations. The frontend is a Vue 3 + Vite SPA served by the backend in production, with authenticated navigation and a real WireGuard peer workflow.

**Tech Stack:** Go 1.22, Gin, SQLite, Vue 3, Vite, Element Plus, Pinia, Vitest

---

## File Map

- Create: `agent_memory/context.md`
- Create: `agent_memory/progress.md`
- Create: `agent_memory/bugs.md`
- Create: `cmd/pi-gateway/main.go`
- Create: `go.mod`
- Create: `internal/app/app.go`
- Create: `internal/config/config.go`
- Create: `internal/api/router.go`
- Create: `internal/api/auth.go`
- Create: `internal/api/sys.go`
- Create: `internal/api/wg.go`
- Create: `internal/api/module_stub.go`
- Create: `internal/api/middleware/auth.go`
- Create: `internal/api/middleware/requestid.go`
- Create: `internal/api/middleware/recovery.go`
- Create: `internal/service/auth_service.go`
- Create: `internal/service/sys_service.go`
- Create: `internal/service/wg_service.go`
- Create: `internal/store/db.go`
- Create: `internal/store/user_store.go`
- Create: `internal/store/wg_peer_store.go`
- Create: `internal/store/setting_store.go`
- Create: `internal/domain/user.go`
- Create: `internal/domain/wg_peer.go`
- Create: `internal/domain/system.go`
- Create: `internal/platform/command_runner.go`
- Create: `internal/platform/systemstats/collector.go`
- Create: `internal/platform/wireguard/manager.go`
- Create: `internal/platform/wireguard/parser.go`
- Create: `internal/platform/wireguard/qr.go`
- Create: `tests` or `internal/..._test.go` files for Go services and handlers
- Create: `web/package.json`
- Create: `web/vite.config.ts`
- Create: `web/tsconfig.json`
- Create: `web/index.html`
- Create: `web/src/main.ts`
- Create: `web/src/App.vue`
- Create: `web/src/router/index.ts`
- Create: `web/src/stores/auth.ts`
- Create: `web/src/api/*.ts`
- Create: `web/src/layouts/AppLayout.vue`
- Create: `web/src/views/*.vue`
- Create: `web/src/components/*.vue`

## Tasks

### Task 1: Establish repository metadata and persistent context

**Files:**
- Create: `agent_memory/context.md`
- Create: `agent_memory/progress.md`
- Create: `agent_memory/bugs.md`

- [ ] Step 1: Record stable project assumptions and current scope.
- [ ] Step 2: Record the active task goal and execution checkpoints.
- [ ] Step 3: Record current blockers, especially lack of local Go runtime.

### Task 2: Define backend skeleton and tests first

**Files:**
- Create: `go.mod`
- Create: `cmd/pi-gateway/main.go`
- Create: `internal/...`
- Test: `internal/service/*_test.go`

- [ ] Step 1: Write failing backend tests for auth password hashing/verification and WireGuard IP allocation.
- [ ] Step 2: Run targeted Go tests when a Go runtime is available and confirm expected failure.
- [ ] Step 3: Implement domain, store interfaces, service logic, and platform adapters with minimal code to satisfy tests.
- [ ] Step 4: Add Gin handlers and router wiring for auth, system stats, WireGuard, and stub module endpoints.
- [ ] Step 5: Re-run the targeted tests and capture pass/fail state.

### Task 3: Build frontend skeleton and page flow

**Files:**
- Create: `web/*`

- [ ] Step 1: Create Vite + Vue app metadata and a failing frontend smoke test if a test runner is configured.
- [ ] Step 2: Implement router, shared API client, auth store, layout, login page, dashboard, WireGuard page, and module stub pages.
- [ ] Step 3: Run frontend build and confirm the SPA compiles.

### Task 4: Document outcomes and verification state

**Files:**
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

- [ ] Step 1: Record what was generated and what remains stubbed.
- [ ] Step 2: Record actual verification commands run and results.
- [ ] Step 3: Record unverified areas caused by missing local runtime or system dependencies.
