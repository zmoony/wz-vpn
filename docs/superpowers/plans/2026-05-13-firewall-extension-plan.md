# Firewall Input Chain Extension Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add input-chain-only firewall management with nftables preview, real apply, and 30-second automatic rollback.

**Architecture:** Add a `firewall_service -> nftables manager` chain that persists input rules in SQLite, renders a complete nftables ruleset with built-in safety rules, applies it via `nft -f`, and records a pending confirmation state on disk. The frontend replaces the current stub with rule CRUD, preview, apply, and confirm flows.

**Tech Stack:** Go 1.22, Gin, SQLite, Vue 3, Element Plus, nftables CLI

---

## File Map

- Modify: `internal/config/config.go`
- Modify: `internal/store/db.go`
- Modify: `internal/app/app.go`
- Modify: `internal/api/router.go`
- Create: `internal/domain/firewall.go`
- Create: `internal/store/firewall_rule_store.go`
- Create: `internal/platform/nftables/manager.go`
- Create: `internal/service/firewall_service.go`
- Create: `internal/service/firewall_service_test.go`
- Create: `internal/api/firewall.go`
- Create: `web/src/api/firewall.ts`
- Create: `web/src/views/FirewallView.vue`
- Modify: `web/src/router/index.ts`
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

### Task 1: Add tests and persistence schema

**Files:**
- Create: `internal/service/firewall_service_test.go`
- Modify: `internal/store/db.go`

- [ ] Step 1: Add tests for rule preview ordering, apply flow, confirm flow, and rollback orchestration.
- [ ] Step 2: Add SQLite table for firewall rules.

### Task 2: Implement backend firewall service and nftables manager

**Files:**
- Create: `internal/domain/firewall.go`
- Create: `internal/store/firewall_rule_store.go`
- Create: `internal/platform/nftables/manager.go`
- Create: `internal/service/firewall_service.go`
- Modify: `internal/config/config.go`
- Modify: `internal/app/app.go`

- [ ] Step 1: Implement domain/store types for input rules and pending apply state.
- [ ] Step 2: Implement nftables manager for render/apply/confirm/rollback.
- [ ] Step 3: Implement firewall service validation, sort ordering, pending state protection, and rollback timer resumption.
- [ ] Step 4: Register the service in app bootstrap.

### Task 3: Expose API and replace frontend firewall stub

**Files:**
- Create: `internal/api/firewall.go`
- Modify: `internal/api/router.go`
- Create: `web/src/api/firewall.ts`
- Create: `web/src/views/FirewallView.vue`
- Modify: `web/src/router/index.ts`

- [ ] Step 1: Add handlers for list/create/update/delete/preview/apply/confirm/pending.
- [ ] Step 2: Replace the stub route with a real firewall page.
- [ ] Step 3: Hook the page to preview, apply, and confirm flows.

### Task 4: Verify what the current environment can prove

**Files:**
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

- [ ] Step 1: Run frontend build and confirm the firewall page compiles.
- [ ] Step 2: Record that live nft execution remains unverified in this environment.
