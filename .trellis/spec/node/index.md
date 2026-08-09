# v2node Development Guidelines

> Coding contracts for this Go node daemon (panel client + xray-based core).

---

## Overview

This repo is **[Jireh012/v2node](https://github.com/Jireh012/v2node)** — a fork of [wyx2685/v2node](https://github.com/wyx2685/v2node) with panel wire changes (ApiPrefix + SM4) for the Java panel in sibling `v2board-java-api`.

| Role | URL | Branch |
|------|-----|--------|
| **This deploy fork** | https://github.com/Jireh012/v2node | `main` + Releases |
| **Upstream (compare/merge)** | https://github.com/wyx2685/v2node | `main` |

Go module path remains `github.com/wyx2685/v2node` (imports / ldflags). That does **not** change which GitHub URL is upstream vs install source.

Panel-side contracts: sibling repo `.trellis/spec/backend/server-node.md`.

---

## Guidelines Index

| Guide | Description | Status |
|-------|-------------|--------|
| [Upstream](./upstream.md) | wyx2685 compare/merge vs Jireh012 deploy | Active |
| [Panel Wire](./panel-wire.md) | ApiPrefix, actions `c\|u\|p\|a\|l`, SM4 query/body | Active |
| [Install & Release](./install-release.md) | install.sh URLs, asset names, publish flow | Active |
| [Directory Structure](./directory-structure.md) | Package layout | Active |

---

## Pre-Development Checklist

- [ ] “对照更新” / merge upstream → [upstream.md](./upstream.md) ([wyx2685/v2node](https://github.com/wyx2685/v2node) `main`)
- [ ] Changing panel HTTP client / crypto → [panel-wire.md](./panel-wire.md); keep parity with Java `NodeSm4Codec`
- [ ] Changing install scripts or CI release → [install-release.md](./install-release.md); download host must stay **Jireh012**
- [ ] Keep intentional divergences: `ApiPrefix` required, SM4 from `ApiKey`, short actions, no classic UniProxy paths
- [ ] After binary-affecting changes: publish a GitHub Release with `v2node-linux-64.zip` / `v2node-linux-arm64-v8a.zip`
- [ ] Tests: `go test ./api/v2board/... ./common/crypt/... ./conf/...`

---

## Language

Spec documents in this directory: **English**. Task PRDs may use Chinese.
