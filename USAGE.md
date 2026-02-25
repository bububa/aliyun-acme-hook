# Usage Guide for Aliyun ACME Hook

## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Configuration](#configuration)
- [Integration with acme.sh](#integration-with-acmesh)
- [Direct Usage](#direct-usage)
- [Supported Services](#supported-services)
- [Troubleshooting](#troubleshooting)

## Overview

Aliyun ACME Hook is a Go-based application that serves as an ACME (Automatic Certificate Management Environment) hook for Alibaba Cloud services, enabling automated SSL/TLS certificate management for domains using Alibaba Cloud's CDN, SLB (Server Load Balancer), CAS (Certificate Authority Service), OSS (Object Storage Service), and FC (Function Compute).

## Installation

### Prerequisites
- Go 1.25.0+ (as specified in go.mod)
- Access to Alibaba Cloud account with appropriate permissions
- acme.sh (for automatic certificate renewal)

### Building from Source
```bash
# Clone the repository
git clone https://github.com/bububa/aliyun-acme-hook.git
cd aliyun-acme-hook

# Build the application
make app
# or directly with Go
go build -o dist/aliyun-acme-hook -ldflags="-s -w -extldflags \"-static\"" ./cmd/app

# Install the binary
sudo make all
```

### Direct Binary Installation
Download the latest release from the GitHub releases page and install to `/usr/local/bin/`.

## Configuration

Create `/etc/aliyun-acme-hook.toml` with your Alibaba Cloud credentials. You can configure multiple services per account:

```toml
[[Accounts]]
Name="production"  # Replace with your account identifier

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

[Accounts.FC]
AK="YOUR_ACCESS_KEY_HERE"      # Replace with your actual access key
SK="YOUR_SECRET_KEY_HERE"      # Replace with your actual secret key
Region="cn-hangzhou"           # Specify your region (FC may require different region)
```

### Configuration Options

- `Name`: A descriptive name for the account configuration
- `AK`: Alibaba Cloud Access Key ID
- `SK`: Alibaba Cloud Access Key Secret
- `Region`: Alibaba Cloud region for the service
- `STSToken`: Optional STS Token for temporary credentials
- `AccountID`: Optional Alibaba Cloud main account ID

### Security Notes
⚠️ **Security Warning**: Never commit real credentials to version control. Store this file securely with appropriate permissions (e.g., `chmod 600 /etc/aliyun-acme-hook.toml`).

## Integration with acme.sh

### Issue Domain Certificates

```bash
acme.sh --issue -d example.com -d *.example.com --dns dns_ali --keylength 2048
```

### Deploy Script

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
  if [ ! -f "$REAL_FULLCHAIN" ] || [ ! -f "$REAL_KEY" ] ; then
    _err "Critical Error: Certificate files not found in RSA directory!"
    return 1
  fi

  # 4. 显式导出变量，让子进程 Go 能够读取
  export CERT_KEY_PATH="$REAL_KEY"
  export FULLCHAIN_PATH="$REAL_FULLCHAIN"
  export CERT_DOMAIN="$domain" # acme.sh 内部的主域名变量是 $domain
  _info "Starting upload to Alibaba Cloud services (CAS, CDN, SLB, OSS, FC)..."

  /usr/local/bin/aliyun-acme-hook -c /etc/aliyun-acme-hook.toml certificate

  if [ $? -eq 0 ] ; then
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

### Deploy Command

```bash
acme.sh --deploy -d example.com --deploy-hook aliyun_acme_hook
```

## Direct Usage

You can also run the tool directly to update certificates:

```bash
aliyun-acme-hook -c /etc/aliyun-acme-hook.toml certificate
```

### Command-line Options

The application accepts the following command-line options:

```bash
NAME:
   aliyun-acme-hook - aliyun acme deploy hook

USAGE:
   aliyun-acme-hook [global options] command [command options]

VERSION:
   v1.2.0

COMMANDS:
   help, h  Shows a list of commands or help for one command
   Aliyun:
     certificate  certificate update in aliyun

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --help, -h              show help
   --version, -v           print the version
```

This command will:
1. Load certificate information from environment variables (set by acme.sh)
2. Upload the certificate to Alibaba Cloud CAS (Certificate Authority Service)
3. Deploy the certificate to CDN domains if CDN configuration is present
4. Deploy the certificate to SLB (Server Load Balancer) if SLB configuration is present
5. Deploy the certificate to OSS (Object Storage Service) if OSS configuration is present
6. Deploy the certificate to FC (Function Compute) if FC configuration is present

## Supported Services

This hook supports deploying certificates to:

### CAS (Certificate Authority Service)
Primary certificate storage in Alibaba Cloud. Certificates are uploaded here first before being distributed to other services.

### CDN (Content Delivery Network)
SSL certificates for CDN domains are updated with the new certificate.

### SLB (Server Load Balancer)
SSL certificates for Server Load Balancer listeners are updated.

### OSS (Object Storage Service)
SSL certificates for custom domains on Object Storage Service are updated.

### FC (Function Compute)
SSL certificates for custom domains on Function Compute are updated.

The service will automatically determine which services to deploy to based on your configuration file. Only services with configuration will be processed.

## Environment Variables

The application expects the following environment variables (typically set by acme.sh):
- `CERT_KEY_PATH`: Path to the private key file
- `FULLCHAIN_PATH`: Path to the full certificate chain file
- `CERT_DOMAIN`: The domain name for the certificate

## Troubleshooting

### Common Issues

1. **Permission denied errors**
   - Ensure your configuration file has proper permissions: `chmod 600 /etc/aliyun-acme-hook.toml`
   - Verify that the aliyun-acme-hook binary is executable: `chmod +x /usr/local/bin/aliyun-acme-hook`

2. **Invalid credentials errors**
   - Double-check your Access Key (AK) and Secret Key (SK) are correct
   - Verify the region is properly set for each service
   - Ensure your Alibaba Cloud account has necessary permissions for each service

3. **Certificate not found errors**
   - Verify that the environment variables `CERT_KEY_PATH`, `FULLCHAIN_PATH`, and `CERT_DOMAIN` are properly set
   - Check that acme.sh generates certificate files in the expected location

4. **Service-specific deployment failures**
   - Some services may require additional permissions beyond basic access keys
   - Verify that the domain is properly registered with the respective service (CDN, SLB, etc.)

### Debugging Tips

- Enable detailed logging by setting the `SLOG_LEVEL` environment variable to `debug`
- Check the Alibaba Cloud console to confirm successful certificate uploads
- Verify service-specific configurations (domain binding, listeners, etc.) are properly set up in Alibaba Cloud console