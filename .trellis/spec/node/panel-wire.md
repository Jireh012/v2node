# Panel Wire (ApiPrefix + SM4)

> HTTP contract between this node and the Java panel obfuscated node API.

---

## Scenario: Obfuscated paths + SM4 client

### 1. Scope / Trigger

- Trigger: Node calls `{ApiHost}{ApiPrefix}/{c|u|p|a|l}` with encrypted identity query `e` and SM4 body envelopes.
- Symptom if broken: falls back to classic UniProxy paths; plaintext `token`/`node_id`; plaintext JSON bodies.

### 2. Signatures

```go
// conf.NodeConfig
ApiHost, NodeID, ApiKey, ApiPrefix, Timeout, RetryCount

// panel.New(c *conf.NodeConfig) (*Client, error)  // requires ApiPrefix + ApiKey
// common/crypt.DeriveNodeWorkingKey(apiKey string) []byte  // SHA-256(UTF-8)[:16]
```

| Action | Method | Client entry |
|--------|--------|--------------|
| `c` | GET | `GetNodeInfo` |
| `u` | GET | `GetUserList` |
| `p` | POST | traffic push |
| `a` | POST | alive report |
| `l` | GET | `GetUserAlive` |

### 3. Contracts

- **ApiPrefix**: required; normalized leading `/`, no trailing `/`. Paths = `ApiPrefix + "/" + action`.
- **Working key**: `SHA-256(UTF-8(ApiKey))[0:16]` — same as panel `server_token` / Java `NodeSm4Codec`. **Not** panel public `SM4_KEY`.
- **Query**: only `e` = `base64url(iv).base64url(ciphertext)` of `{"k":"<ApiKey>","i":<NodeID>,"t":"vn"}`.
- **Body**: JSON `{iv,payload}` (standard base64) for POST bodies and success responses. No msgpack on these paths.
- **Config JSON** (`/etc/v2node/config.json` Nodes[]): `ApiHost`, `NodeID`, `ApiKey`, `ApiPrefix` (mapstructure tags as in `conf.NodeConfig`).

### 4. Validation & Error Matrix

| Condition | Result |
|-----------|--------|
| Empty `ApiPrefix` | `panel.New` error `"ApiPrefix is required"` |
| Empty `ApiKey` | `panel.New` error `"ApiKey is required"` |
| Response missing `iv`/`payload` | decrypt error; caller fails or soft-fails (alivelist) |
| HTTP ≥399 on alivelist | soft-fail empty map (do not hard-crash node) |

### 5. Good / Base / Bad Cases

- **Good**: `ApiPrefix=/n/xxxxxxxxxxxx`, actions `/n/xxxxxxxxxxxx/u`, `e` present, body envelope round-trips.
- **Base**: install writes all four fields; `t` always `"vn"`.
- **Bad**: classic `/api/v1/server/UniProxy/user?token=...`; missing ApiPrefix; separate Sm4Key config field.

### 6. Tests Required

- `api/v2board/path_test.go` — normalize prefix; `New` requires ApiPrefix; `actionPath`.
- `common/crypt` — key derive, compact `e`, envelope encrypt/decrypt (keep parity with Java vectors when present).

### 7. Wrong vs Correct

#### Wrong
```json
{ "ApiHost": "https://panel.example", "NodeID": 1, "ApiKey": "..." }
```
(no ApiPrefix; classic paths)

#### Correct
```json
{
  "ApiHost": "https://api.example.com",
  "NodeID": 1,
  "ApiKey": "<server_token>",
  "ApiPrefix": "/n/xxxxxxxxxxxx"
}
```

---

## Design Decision: No separate Sm4Key

**Context**: Avoid a second secret in install flags / config.

**Decision**: Derive SM4 working key from `ApiKey` only. Rotating panel 通讯密钥 requires updating node `ApiKey`.
