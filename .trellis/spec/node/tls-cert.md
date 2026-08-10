# TLS certificate issuance (dns / http / self)

## Scenario: SNI / CertDomain change must re-issue

### 1. Scope / Trigger

- Admin changes v2node `tls_settings.server_name` (SNI). Panel clients immediately use the new SNI.
- Node `CertInfo.CertDomain` = `PrimaryServerName()` from that field.
- Symptom if broken: client `x509: certificate is valid for OLD, not NEW` after SNI edit.

### 2. Contracts

| Mode | Behavior |
|------|----------|
| `dns` / `http` | If cert files exist **and** leaf SAN/CN covers `CertDomain` → keep. Else delete + `CreateCert`. |
| `self` | Same domain check; regenerate self-signed when mismatch. |
| `file` / `none` | Unchanged (operator-managed files). |
| Renew task | Domain mismatch → re-issue; else expiry-based renew. |

### 3. Wrong vs Correct

#### Wrong

```go
if file.IsExist(cert.CertFile) && file.IsExist(cert.KeyFile) {
    return nil // ignores CertDomain change
}
```

#### Correct

```go
ok, _ := existingCertMatchesDomain(cert.CertFile, cert.CertDomain)
if ok {
    return nil
}
// remove + CreateCert / generateSelfSslCertificate
```

### 4. Ops note (CDN)

- `host` may be a CDN hostname; `server_name` must be the **certificate** domain.
- Do not set SNI equal to CDN host unless that exact name is on the issued cert.
