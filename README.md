# v2node

A v2board backend based on modified xray-core.  
基于修改版 xray 内核的节点服务端（对接 Java 面板混淆节点 API）。

## 面板对接（硬切换）

节点只走面板配置的中性前缀，不再使用 `/api/v2/server/**` 或 `/api/v1/server/UniProxy/**`：

```text
{ApiPrefix}/{c|u|p|a|l}?e=<SM4 compact>
```

| 配置项 | 说明 |
|--------|------|
| `ApiHost` | 面板根 URL |
| `ApiKey` | 通讯密钥（≥16）；同时用于 SM4 工作密钥派生 `SHA-256(ApiKey)[:16]` |
| `ApiPrefix` | 与面板 `server.server_api_prefix` 一致（如 `/n/xxxxxxxxxxxx`） |
| `NodeID` | 节点 ID |

成功响应与 POST 体均为 SM4 信封 `{iv,payload}`；身份仅通过密文 query `e` 传递。

## 一键安装

```bash
wget -N https://raw.githubusercontent.com/Jireh012/v2node/main/script/install.sh && bash install.sh \
  --api-host 'https://panel.example.com' \
  --node-id 1 \
  --api-key '<通讯密钥>' \
  --api-prefix '/n/xxxxxxxxxxxx'
```

无 `--sm4-key`：SM4 由 `--api-key` 派生。

## 构建

```bash
GOEXPERIMENT=jsonv2 go build -v -o build_assets/v2node -trimpath -ldflags "-X 'github.com/wyx2685/v2node/cmd.version=$version' -s -w -buildid="
```

## 测试

```bash
go test ./common/crypt/... ./api/v2board/... ./conf/...
```
