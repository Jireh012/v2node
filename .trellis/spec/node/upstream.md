# Upstream Reference (wyx2685/v2node)

> Canonical upstream for compare / merge / “对照更新”.

---

## Convention: Upstream repository

**What**: Treat **[wyx2685/v2node](https://github.com/wyx2685/v2node)** (`main`) as the only default upstream when comparing, cherry-picking, or merging into this fork.

**Why**: This tree is [Jireh012/v2node](https://github.com/Jireh012/v2node). Without a fixed upstream,对照更新 drifts across random forks.

**Do**:
- `git remote add upstream https://github.com/wyx2685/v2node.git` (once); `git fetch upstream` then diff `upstream/main`.
- Focus paths: `api/v2board/`, `node/`, `conf/`, `script/`, `core/`, `common/`, `.github/workflows/`.
- After merge: re-check [panel-wire.md](./panel-wire.md) divergences; if binaries change, publish a **Jireh012** Release ([install-release.md](./install-release.md)).

**Don't**:
- Use another fork as the compare baseline unless Trellis is updated.
- Point `script/install.sh` download URLs at wyx2685 Releases (install stays Jireh012).
- Drop ApiPrefix/SM4 when resolving merge conflicts.

---

## Roles (do not mix)

| Role | Repo | Use for |
|------|------|---------|
| Upstream (compare) | [wyx2685/v2node](https://github.com/wyx2685/v2node) | Diff, merge,对照更新 |
| This fork (runtime) | [Jireh012/v2node](https://github.com/Jireh012/v2node) | Panel install_command, production nodes |

---

## Workflow checklist

1. Fetch/compare **wyx2685/v2node** `main`.
2. List intentional divergences to keep (see [panel-wire.md](./panel-wire.md) + install URLs → Jireh012).
3. Merge/cherry-pick into `main`; fix conflicts without dropping SM4/prefix.
4. `go test` on touched packages.
5. Publish Release if install zip contents change.
6. Document new intentional divergences in this Trellis layer.

---

## Design Decision: Single canonical upstream

**Decision**: Compare/merge upstream = `https://github.com/wyx2685/v2node` @ `main`. Deploy/install = this repo. Change only via explicit Trellis update.
