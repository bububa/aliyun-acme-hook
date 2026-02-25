package config

type Config struct {
	Accounts []Account `required:"true"`
}

type Account struct {
	// Name account description
	Name string `required:"true"`
	// CAS cas config if do not have will use raw certification/key
	CAS *AliyunConfig
	// CDN cdn config ignore update if do not have
	CDN *AliyunConfig
	// SLB slb config ignore update if do not have
	SLB *AliyunConfig
	// OSS oss config ignore update if do not have
	OSS *AliyunConfig
	// FC serverless config ignore update if do not have
	FC *AliyunConfig
}

type AliyunConfig struct {
	// AK access key
	AK string `required:"true"`
	// SK access secret key
	SK string `required:"true"`
	// STSToken sts token
	STSToken string
	// AccountID aliyun main account id
	AccountID string
	// Region API service region
	Region string
}
