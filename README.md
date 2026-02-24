# Aliyun acme.sh deploy hook

[![Go Reference](https://pkg.go.dev/badge/github.com/bububa/aliyun-acme-hook.svg)](https://pkg.go.dev/github.com/bububa/aliyun-acme-hook)
[![Go](https://github.com/bububa/aliyun-acme-hook/actions/workflows/go.yml/badge.svg)](https://github.com/bububa/aliyun-acme-hook/actions/workflows/go.yml)
[![goreleaser](https://github.com/bububa/aliyun-acme-hook/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/bububa/aliyun-acme-hook/actions/workflows/goreleaser.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/bububa/aliyun-acme-hook.svg)](https://github.com/bububa/aliyun-acme-hook)
[![GoReportCard](https://goreportcard.com/badge/github.com/bububa/aliyun-acme-hook)](https://goreportcard.com/report/github.com/bububa/aliyun-acme-hook)
[![GitHub license](https://img.shields.io/github/license/bububa/aliyun-acme-hook.svg)](https://github.com/bububa/aliyun-acme-hook/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/bububa/aliyun-acme-hook.svg)](https://GitHub.com/bububa/aliyun-acme-hook/releases/)

## issue domain with acme.sh

```bash
acme.sh --issue -d example.com -d *.example.com --dns dns_ali --keylength 2048
```

## configuration

Create `/etc/aliyun-acme-hook.toml` with your Alibaba Cloud credentials:

```toml
[[Accounts]]
Name="my-account-name"  # Replace with your account identifier

[Accounts.CAS]
AK="YOUR_ACCESS_KEY_HERE"      # Replace with your actual access key
SK="YOUR_SECRET_KEY_HERE"      # Replace with your actual secret key
Region="cn-zhangjiakou"        # Specify your region

[Accounts.CDN]
AK="YOUR_ACCESS_KEY_HERE"      # Replace with your actual access key
SK="YOUR_SECRET_KEY_HERE"      # Replace with your actual secret key
Region="cn-zhangjiakou"        # Specify your region

[Accounts.SLB]
AK="YOUR_ACCESS_KEY_HERE"      # Replace with your actual access key
SK="YOUR_SECRET_KEY_HERE"      # Replace with your actual secret key
Region="cn-zhangjiakou"        # Specify your region

[Accounts.OSS]
AK="YOUR_ACCESS_KEY_HERE"      # Replace with your actual access key
SK="YOUR_SECRET_KEY_HERE"      # Replace with your actual secret key
Region="cn-zhangjiakou"        # Specify your region
```

⚠️ **Security Warning**: Never commit real credentials to version control. Store this file securely with appropriate permissions (e.g., `chmod 600 /etc/aliyun-acme-hook.toml`).

## deploy script

Create `~/.acme.sh/deploy/aliyun_acme_hook.sh`:

```bash
#!/bin/bash

# acme.sh 会自动调用以 _deploy 结尾的函数
aliyun_acme_hook_deploy() {
  # 1. 强制定义路径
  # $domain 是 acme.sh 提供的当前域名 (example.com)
  # $CERT_HOME 是 acme.sh 的安装根目录 (通常是 /root/.acme.sh)

  REAL_FULLCHAIN="$CERT_HOME/$domain/fullchain.cer"
  REAL_KEY="$CERT_HOME/$domain/$domain.key"

  # 2. 打印调试信息，确保我们在日志里能看到路径
  _info "Resolved FULLCHAIN_PATH: $REAL_FULLCHAIN"
  _info "Resolved CERT_KEY_PATH: $REAL_KEY"

  # 3. 验证文件是否真的存在
  if [ ! -f "$REAL_FULLCHAIN" ] || [ ! -f "$REAL_KEY" ]; then
    _err "Critical Error: Certificate files not found in RSA directory!"
    return 1
  fi

  # 4. 显式导出变量，让子进程 Go 能够读取
  export CERT_KEY_PATH="$REAL_KEY"
  export FULLCHAIN_PATH="$REAL_FULLCHAIN"
  export CERT_DOMAIN="$domain" # acme.sh 内部的主域名变量是 $domain
  _info "Starting upload to Alibaba Cloud services (CAS, CDN, SLB)..."

  /usr/local/bin/aliyun-acme-hook -c /etc/aliyun-acme-hook.toml certificate

  if [ $? -eq 0 ]; then
    _info "Aliyun Certificate Deployment Success."
    return 0
  else
    _err "Aliyun Certificate Deployment Failed."
    return 1
  fi
}
```

Make sure the script is executable:

```bash
chmod +x ~/.acme.sh/deploy/aliyun_acme_hook.sh
```

## deploy command

```bash
acme.sh --deploy -d example.com --deploy-hook aliyun_acme_hook
```

## command-line usage

You can also run the tool directly to update certificates:

```bash
aliyun-acme-hook -c /etc/aliyun-acme-hook.toml certificate
```

This will:
1. Load certificate information from environment variables (set by acme.sh)
2. Upload the certificate to Alibaba Cloud CAS (Certificate Authority Service)
3. Deploy the certificate to CDN domains if CDN configuration is present
4. Deploy the certificate to SLB (Server Load Balancer) if SLB configuration is present

## services supported

This hook supports deploying certificates to:
- **CAS** (Certificate Authority Service): Primary certificate storage
- **CDN**: Content Delivery Network SSL certificates  
- **SLB**: Server Load Balancer SSL certificates
- **OSS**: Object Storage Service SSL certificates for custom domains

The service will automatically determine which services to deploy to based on your configuration file.
```
