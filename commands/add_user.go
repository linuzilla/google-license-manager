package commands

import (
	"fmt"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/services/license_manager"
)

type AddUserCommand struct {
	LicenseManager license_manager.LicenseManager `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*AddUserCommand)(nil)

func (AddUserCommand) PostSummerConstruct() {
}

func (AddUserCommand) ImplementCommandInterface() {
}

func (AddUserCommand) Command() string {
	return `add-user`
}

func (cmd *AddUserCommand) Execute(args ...string) int {
	if len(args) == 0 {
		fmt.Println("usage: add-user user ...")
	} else {
		for _, user := range args {
			cmd.LicenseManager.AddLicenseToUser(user)
		}
	}
	return 0
}
