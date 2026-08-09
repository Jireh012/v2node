# Directory Structure

> Package layout for this Go module.

---

## Layout

```
.
├── main.go                 # entry
├── cmd/                    # CLI / version
├── api/v2board/            # panel HTTP client (ApiPrefix + SM4)
├── conf/                   # config.json load/watch (NodeConfig)
├── common/crypt/           # SM4 derive / envelope / compact e
├── common/{counter,exec,file,format,rate,systime,task}/
├── node/                   # node lifecycle, sync with panel
├── core/                   # xray / inbound wiring
├── limiter/                # speed / device limits
├── script/                 # install.sh, v2node.sh
└── .github/workflows/      # release builds
```

Optional git submodule `sing-box_mod` (see `.gitmodules`) is upstream xray-related; not part of Trellis `packages` unless explicitly enabled.

---

## Conventions

- Panel protocol changes belong in `api/v2board` + `common/crypt` + `conf`; keep install flags in `script/`.
- Prefer small focused packages; do not invent a second SM4 key field.
- Module path `github.com/wyx2685/v2node` stays for Go imports; Git remotes: `origin` = Jireh012, `upstream` = wyx2685.
