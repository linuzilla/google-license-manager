package commands

import (
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/services/admin_sdk"
	"github.com/linuzilla/google-license-manager/services/google_credential"
)

type CustomerCommand struct {
	AdminSdk admin_sdk.GoogleAdminSdk `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*CustomerCommand)(nil)
var _ ConditionalRegister = (*CustomerCommand)(nil)

func (CustomerCommand) PostSummerConstruct() {
}

func (CustomerCommand) ImplementCommandInterface() {
}

func (CustomerCommand) CanRegister(googleCredential google_credential.GoogleCredential) bool {
	return googleCredential.HasCustomerInfo()
}

func (CustomerCommand) Command() string {
	return `customer`
}

func (cmd *CustomerCommand) Execute(args ...string) int {
	cmd.AdminSdk.CustomerInfo()
	return 0
}
