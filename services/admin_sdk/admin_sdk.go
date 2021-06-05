package admin_sdk

import (
	"context"
	"fmt"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/constants"
	"github.com/linuzilla/google-license-manager/services/google_credential"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
	"log"
	"sync"
)

type GoogleAdminSdk interface {
	UserInfo(userId string) (*admin.User, error)
	CustomerInfo()
	HasUserInfo() bool
	HasCustomerInfo() bool
}

type googleAdminSdkImpl struct {
	CredentialService google_credential.GoogleCredential `inject:"*"`
	GoogleConfig      *config.GoogleCfg                  `inject:"*"`
	initializeOnce    sync.Once
	service           *admin.Service
}

func New() GoogleAdminSdk {
	return &googleAdminSdkImpl{}
}

func (impl *googleAdminSdkImpl) PostSummerConstruct() {
	impl.CredentialService.HasScope(admin.AdminDirectoryUserReadonlyScope)
}

func (impl *googleAdminSdkImpl) loadService() *admin.Service {
	impl.initializeOnce.Do(func() {
		ctx := context.Background()

		credentialConfig := impl.CredentialService.LoadConfig()
		credentialConfig.Client(ctx)

		service, err := admin.NewService(ctx,
			option.WithHTTPClient(credentialConfig.Client(ctx)),
			option.WithUserAgent(constants.UserAgent),
			option.WithScopes(admin.AdminDirectoryUserReadonlyScope),
			option.WithScopes(admin.AdminDirectoryCustomerReadonlyScope),
		)

		if err != nil {
			log.Fatal(err)
		}
		impl.service = service
	})

	return impl.service
}

func (impl *googleAdminSdkImpl) HasUserInfo() bool {
	return impl.CredentialService.HasScope(admin.AdminDirectoryUserReadonlyScope)
}

func (impl *googleAdminSdkImpl) UserInfo(userId string) (*admin.User, error) {
	userKey := impl.CredentialService.NormalizeUserId(userId)

	return impl.loadService().Users.Get(userKey).Do()
}

func (impl *googleAdminSdkImpl) HasCustomerInfo() bool {
	return impl.CredentialService.HasScope(admin.AdminDirectoryCustomerReadonlyScope)
}

func (impl *googleAdminSdkImpl) CustomerInfo() {
	customer, err := impl.loadService().Customers.Get(impl.GoogleConfig.CustomerId).Do()

	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Domain: %s\n", customer.CustomerDomain)
		fmt.Printf("Language: %s\n", customer.Language)
		fmt.Printf("Phone: %s\n", customer.PhoneNumber)
		fmt.Printf("Alternate Email: %s\n", customer.AlternateEmail)
		fmt.Printf("Contact Name: %s\n", customer.PostalAddress.ContactName)
	}
}
