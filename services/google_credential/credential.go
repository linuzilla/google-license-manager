package google_credential

import (
	"fmt"
	"github.com/linuzilla/google-license-manager/config"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/licensing/v1"
	"io/ioutil"
	"log"
	"strings"
)

type GoogleCredential interface {
	LoadConfig() *jwt.Config
	NormalizeUserId(userId string) string
	HasScope(scope string) bool
	HasUserInfo() bool
	HasCustomerInfo() bool
}

type googleCredentialImpl struct {
	googleConfig *config.GoogleCfg
	cfg          *jwt.Config
	scopesMap    map[string]bool
}

var singleton = googleCredentialImpl{
	scopesMap: make(map[string]bool),
}

func (impl *googleCredentialImpl) initScopes() {
	for _, scope := range impl.googleConfig.Scopes {
		impl.scopesMap[scope] = true
	}
	impl.scopesMap[licensing.AppsLicensingScope] = true
}
func (impl *googleCredentialImpl) postConstruct(cfg *jwt.Config) GoogleCredential {
	impl.cfg = cfg
	impl.cfg.Subject = impl.googleConfig.Admin

	return impl
}

func New(googleConfig *config.GoogleCfg) GoogleCredential {
	singleton.googleConfig = googleConfig
	singleton.initScopes()

	fmt.Printf("Loading credential: \"%s\"\n", googleConfig.CredentialFile)
	cfg, err := singleton.loadConfig(googleConfig.CredentialFile)

	if err != nil {
		log.Fatal(err)
	}

	return singleton.postConstruct(cfg)
}

func FromData(googleConfig *config.GoogleCfg, credential string) GoogleCredential {
	singleton.googleConfig = googleConfig
	singleton.initScopes()

	cfg, err := singleton.jwtConfig([]byte(credential))

	if err != nil {
		log.Fatal(err)
	}

	return singleton.postConstruct(cfg)
}

func (impl *googleCredentialImpl) HasScope(scope string) bool {
	_, found := impl.scopesMap[scope]
	return found
}

func (impl *googleCredentialImpl) HasUserInfo() bool {
	return impl.HasScope(admin.AdminDirectoryUserReadonlyScope)
}

func (impl *googleCredentialImpl) HasCustomerInfo() bool {
	return impl.HasScope(admin.AdminDirectoryCustomerReadonlyScope)
}

func (impl *googleCredentialImpl) jwtConfig(b []byte) (*jwt.Config, error) {
	var scopes []string

	for k := range impl.scopesMap {
		scopes = append(scopes, k)
	}

	return google.JWTConfigFromJSON(b, scopes...)
}

func (impl *googleCredentialImpl) loadConfig(credentialFile string) (*jwt.Config, error) {
	if b, err := ioutil.ReadFile(credentialFile); err != nil {
		log.Fatalf("Unable to read client secret file: %v\n", err)
		return nil, err
	} else {
		return impl.jwtConfig(b)
	}
}

func (impl *googleCredentialImpl) LoadConfig() *jwt.Config {
	return impl.cfg
}

func (impl *googleCredentialImpl) NormalizeUserId(userId string) string {
	if !strings.HasSuffix(userId, impl.googleConfig.Domain) {
		return userId + `@` + impl.googleConfig.Domain
	} else {
		return userId
	}
}
