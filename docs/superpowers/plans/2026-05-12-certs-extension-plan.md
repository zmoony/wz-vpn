# Certificate Management Extension Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend Pi-Gateway with real certificate management for a single root domain wildcard certificate using acme.sh, certificate status parsing, and nginx reload integration.

**Architecture:** Add a `cert_service -> acme manager` backend chain that persists one-or-more certificate definitions in SQLite, runs acme.sh issue/install/renew commands, and reads installed certificate files for status reporting. Reuse the current panel structure so frontend certificate management fits beside DDNS and reverse proxy pages.

**Tech Stack:** Go 1.22, Gin, SQLite, Vue 3, Element Plus, acme.sh CLI, x509 certificate parsing

---

## File Map

- Modify: `internal/config/config.go`
- Modify: `internal/store/db.go`
- Modify: `internal/app/app.go`
- Modify: `internal/api/router.go`
- Create: `internal/domain/cert.go`
- Create: `internal/store/cert_config_store.go`
- Create: `internal/platform/acme/manager.go`
- Create: `internal/service/cert_service.go`
- Create: `internal/service/cert_service_test.go`
- Create: `internal/api/cert.go`
- Create: `web/src/api/certs.ts`
- Create: `web/src/views/CertsView.vue`
- Modify: `web/src/router/index.ts`
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

### Task 1: Add failing tests and persistence schema

**Files:**
- Create: `internal/service/cert_service_test.go`
- Modify: `internal/store/db.go`

- [ ] Step 1: Add service tests covering wildcard config creation, issue orchestration, renew orchestration, and status hydration.
- [ ] Step 2: Add the certificate configuration table to SQLite migrations.

### Task 2: Implement backend certificate services and acme adapter

**Files:**
- Create: `internal/domain/cert.go`
- Create: `internal/store/cert_config_store.go`
- Create: `internal/platform/acme/manager.go`
- Create: `internal/service/cert_service.go`
- Modify: `internal/config/config.go`
- Modify: `internal/app/app.go`

- [ ] Step 1: Implement domain/store structures for certificate config and status.
- [ ] Step 2: Implement acme.sh issue/install/renew command adapter plus installed certificate parsing.
- [ ] Step 3: Implement service validation and orchestration.
- [ ] Step 4: Register the new service in app bootstrap.

### Task 3: Expose API and replace frontend certificate stub

**Files:**
- Create: `internal/api/cert.go`
- Modify: `internal/api/router.go`
- Create: `web/src/api/certs.ts`
- Create: `web/src/views/CertsView.vue`
- Modify: `web/src/router/index.ts`

- [ ] Step 1: Add list/create/update/renew handlers.
- [ ] Step 2: Replace the cert stub route with a real page.
- [ ] Step 3: Hook the page to issue and renew actions.

### Task 4: Verify what the current environment can prove

**Files:**
- Modify: `agent_memory/progress.md`
- Modify: `agent_memory/bugs.md`

- [ ] Step 1: Run frontend build and confirm the certificate page compiles.
- [ ] Step 2: Record that Go compile/test and live acme.sh execution remain unverified here.
