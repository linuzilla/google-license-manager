package license_manager

import (
	"context"
	"fmt"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/constants"
	"github.com/linuzilla/google-license-manager/services/google_credential"
	"google.golang.org/api/licensing/v1"
	"google.golang.org/api/option"
	"log"
	"sync"
)

type LicenseManager interface {
	AddLicenseToUser(userId string)
	RevokeLicenseFromUser(userId string)
	ListLicenseUser() []*licensing.LicenseAssignment
}

type licenseManagerImpl struct {
	CredentialService google_credential.GoogleCredential `inject:"*"`
	GoogleConfig      *config.GoogleCfg                  `inject:"*"`
	service           *licensing.Service
	newServiceOnce    sync.Once
}

func New() LicenseManager {
	return &licenseManagerImpl{}
}

func (impl *licenseManagerImpl) loadService() *licensing.Service {
	impl.newServiceOnce.Do(func() {
		ctx := context.Background()

		credentialConfig := impl.CredentialService.LoadConfig()
		credentialConfig.Client(ctx)

		service, err := licensing.NewService(ctx,
			option.WithHTTPClient(credentialConfig.Client(ctx)),
			option.WithUserAgent(constants.UserAgent),
			option.WithScopes(licensing.AppsLicensingScope))
		if err != nil {
			log.Fatal(err)
		}
		impl.service = service
	})

	return impl.service
}

func (impl *licenseManagerImpl) AddLicenseToUser(userId string) {
	service := impl.loadService()

	if assignment, err := service.LicenseAssignments.Insert(impl.GoogleConfig.ProductId, impl.GoogleConfig.ProductSkuId, &licensing.LicenseAssignmentInsert{
		UserId: impl.CredentialService.NormalizeUserId(userId),
	}).Do(); err != nil {
		log.Print(err)
	} else {
		fmt.Printf("%s added\n", assignment.UserId)
	}
}

func (impl *licenseManagerImpl) RevokeLicenseFromUser(userId string) {
	service := impl.loadService()
	user := impl.CredentialService.NormalizeUserId(userId)

	if assignment, err := service.LicenseAssignments.Delete(impl.GoogleConfig.ProductId, impl.GoogleConfig.ProductSkuId, user).Do(); err != nil {
		log.Print(err)
	} else {
		fmt.Printf("%s revoked (%d)\n", user, assignment.HTTPStatusCode)
	}
}

func (impl *licenseManagerImpl) ListLicenseUser() []*licensing.LicenseAssignment {
	service := impl.loadService()
	var licenses []*licensing.LicenseAssignment
	pageToken := ``

	fields := service.LicenseAssignments.ListForProductAndSku(impl.GoogleConfig.ProductId, impl.GoogleConfig.ProductSkuId, impl.GoogleConfig.CustomerId).
		Fields(`nextPageToken, items(productId, userId, skuId, skuName, productName)`)

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != `` {
		if r, err := fields.PageToken(pageToken).Do(); err != nil {
			log.Print(err)
		} else {
			for _, item := range r.Items {
				licenses = append(licenses, item)
			}

			pageToken = r.NextPageToken
		}
	}

	return licenses
}
