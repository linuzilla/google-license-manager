package config

type GoogleCfg struct {
	Admin          string   `yaml:"admin"`
	Domain         string   `yaml:"domain"`
	CustomerId     string   `yaml:"customer-id"`
	ProductId      string   `yaml:"product-id"`
	ProductSkuId   string   `yaml:"product-sku-id"`
	CredentialFile string   `yaml:"credential-file"`
	Scopes         []string `yaml:"scopes"`
}
