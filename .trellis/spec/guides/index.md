# Thinking Guides

> **Purpose**: Expand thinking before coding. Point to code-specs for executable contracts.

---

## Available Guides

| Guide | Purpose | When to Use |
|-------|---------|-------------|
| [Code Reuse Thinking Guide](./code-reuse-thinking-guide.md) | Patterns and duplication | Repeated logic / helpers |
| [Cross-Layer Thinking Guide](./cross-layer-thinking-guide.md) | Data flow across packages | Panel wire ↔ node ↔ core |

**v2node code-specs**: [node/index.md](../node/index.md)

**Upstream compare**: [upstream.md](../node/upstream.md) → **https://github.com/wyx2685/v2node**

**Deploy/install**: this repo **https://github.com/Jireh012/v2node** → [install-release.md](../node/install-release.md)

Sibling Java panel contracts: `v2board-java-api` `.trellis/spec/backend/server-node.md`.

---

## Quick Reference: Thinking Triggers

### When Aligning / Merging Upstream (“对照更新”)

- [ ] Compare against **only** [wyx2685/v2node](https://github.com/wyx2685/v2node) `main` → [upstream.md](../node/upstream.md)
- [ ] Preserve ApiPrefix + SM4 + Jireh012 install URLs → [panel-wire.md](../node/panel-wire.md) + [install-release.md](../node/install-release.md)
- [ ] Binary change → publish Jireh012 Release assets

### When Changing Panel Client / Crypto

- [ ] Paths stay `{ApiPrefix}/{c|u|p|a|l}`; query only `e`; body `{iv,payload}`
- [ ] Working key = `SHA-256(ApiKey)[:16]` — no separate Sm4Key
- [ ] Keep soft-fail behavior on alivelist HTTP ≥399

### When Changing Install / CI

- [ ] Download URLs must remain `Jireh012/v2node` (never wyx2685 for this panel)
- [ ] Asset names match `friendly-filenames.json` / install.sh `arch` mapping

### When to Think About Cross-Layer Issues

- [ ] Feature touches panel API + node client + config/install
- [ ] Changing identity JSON / action letters / envelope encoding
- [ ] Traffic / alive thresholds from `base_config` enforced on node

→ Read [Cross-Layer Thinking Guide](./cross-layer-thinking-guide.md)

### When to Think About Code Reuse

- [ ] Duplicating SM4 or path helpers outside `common/crypt` / `api/v2board`
- [ ] New config field mirrored in install.sh and conf.NodeConfig

→ Read [Code Reuse Thinking Guide](./code-reuse-thinking-guide.md)
