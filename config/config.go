package config

type Config struct {
	Accounts []Account `required:"true"`
}

type Account struct {
	Name string `required:"true"`
	CAS  *AliyunConfig
	CDN  *AliyunConfig
	SLB  *AliyunConfig
	OSS  *AliyunConfig
}

type AliyunConfig struct {
	AK       string `required:"true"`
	SK       string `required:"true"`
	STSToken string
	Region   string
}
