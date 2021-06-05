package commands

import (
	"fmt"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/services/license_manager"
)

type RevokeUserCommand struct {
	LicenseManager license_manager.LicenseManager `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*RevokeUserCommand)(nil)

func (RevokeUserCommand) PostSummerConstruct() {
}

func (RevokeUserCommand) ImplementCommandInterface() {
}

func (RevokeUserCommand) Command() string {
	return `revoke-user`
}

func (cmd *RevokeUserCommand) Execute(args ...string) int {
	if len(args) == 0 {
		fmt.Println("usage: revoke-user user ...")
	} else {
		for _, user := range args {
			cmd.LicenseManager.RevokeLicenseFromUser(user)
		}
	}
	return 0
}
