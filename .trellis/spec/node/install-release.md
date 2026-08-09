# Install Script & GitHub Release

> One-click install and binary distribution for this fork.

---

## Scenario: Jireh012 install distribution

### 1. Scope / Trigger

- Trigger: Panel `install_command` and `script/install.sh` download script + Linux zip from GitHub.
- Symptom if broken: scripts mention ApiPrefix but zip comes from upstream → node lacks SM4 wire.
- Gotcha: pushing scripts alone is not enough; install resolves `releases/latest`.

### 2. Signatures

```bash
# Panel / docs
wget -N https://raw.githubusercontent.com/Jireh012/v2node/main/script/install.sh && bash install.sh \
  --api-host '...' --node-id N --api-key '...' --api-prefix '/n/...'

# install.sh
# GET https://api.github.com/repos/Jireh012/v2node/releases/latest → tag_name
# GET https://github.com/Jireh012/v2node/releases/download/${tag}/v2node-linux-${arch}.zip
# manage script: raw .../Jireh012/v2node/main/script/v2node.sh
```

### 3. Contracts

| Piece | Canonical value |
|-------|-----------------|
| Raw scripts branch | `main` |
| Release owner/repo | `Jireh012/v2node` |
| Common assets | `v2node-linux-64.zip`, `v2node-linux-arm64-v8a.zip` (see `.github/build/friendly-filenames.json`) |
| Zip contents | `v2node`, `geoip.dat`, `geosite.dat` (+ README/LICENSE) |
| Config flags from install | `ApiHost`, `NodeID`, `ApiKey`, `ApiPrefix` |
| Forbidden | `wyx2685/v2node` as download source for this panel’s installs |

`--api-host` must be the **panel API origin** (not a Vite frontend port).

CI: `.github/workflows/release.yml` uploads zips on `release: published` (and builds on push to `main`/`master` for Go changes).

### 4. Validation & Error Matrix

| Condition | Result |
|-----------|--------|
| No Release / API rate limit | install exits: 检测版本失败 |
| Wrong asset name | download/unzip fail |
| Binary without ApiPrefix support | cannot talk to obfuscated panel |

### 5. Good / Base / Bad Cases

- **Good**: `v0.5.0` (or newer) with linux-64/arm64 zips; binary contains `ApiPrefix`.
- **Base**: install.sh + v2node.sh URLs all `Jireh012/v2node`.
- **Bad**: fork scripts only; Releases empty or still wyx2685.

### 6. Tests Required

- Grep/CI: scripts must not contain `wyx2685/v2node` download URLs.
- Smoke: `releases/latest` + unzip lists `v2node`, dat files.

### 7. Wrong vs Correct

#### Wrong
```text
raw Jireh012 install.sh → downloads wyx2685 release zip
```

#### Correct
```text
raw Jireh012 install.sh → downloads Jireh012 release zip (SM4 binary)
```

---

## Publish checklist

1. Merge to `main`.
2. `gh release create <tag> --target main` (or tag + publish).
3. Ensure assets uploaded (CI or `gh release upload`).
4. Verify `api.github.com/repos/Jireh012/v2node/releases/latest`.
