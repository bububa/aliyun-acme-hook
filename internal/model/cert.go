package model

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Cert struct {
	ID        int64
	Name      string
	Domain    string
	Key       []byte
	FullChain []byte
}

func (c *Cert) LoadFromEnv() error {
	// 1. 获取 acme.sh 提供的环境变量
	// acme.sh 会在执行 hook 时自动 export 这些路径
	// certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("CERT_KEY_PATH")
	fullChainPath := os.Getenv("FULLCHAIN_PATH")
	c.Domain = os.Getenv("CERT_DOMAIN") // acme.sh 传入的主域名

	if fullChainPath == "" || keyPath == "" {
		return errors.New("missing certificate path environment variables")
	}

	// 2. Validate domain name to prevent injection or malicious input
	if err := c.validateDomain(); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	// 3. 读取证书内容
	certContent, err := os.ReadFile(fullChainPath)
	if err != nil {
		return err
	}
	c.FullChain = certContent
	keyContent, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}
	c.Key = keyContent
	return nil
}

// validateDomain validates the domain name to prevent injection attacks
func (c *Cert) validateDomain() error {
	if c.Domain == "" {
		return nil // Allow empty domain, may be optional in some cases
	}

	// Basic domain validation using regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]{1,61}[a-zA-Z0-9](\.[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])*\.?$`)
	if !domainRegex.MatchString(c.Domain) {
		return fmt.Errorf("invalid domain format: %s", c.Domain)
	}

	// Additional validation: ensure domain is not too long
	if len(c.Domain) > 253 {
		return fmt.Errorf("domain name too long: %s", c.Domain)
	}

	// Additional validation: ensure each label is not too long
	labels := strings.Split(strings.TrimRight(c.Domain, "."), ".")
	for _, label := range labels {
		if len(label) > 63 {
			return fmt.Errorf("domain label too long: %s", label)
		}
	}

	return nil
}
