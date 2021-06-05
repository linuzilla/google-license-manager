package commands

import (
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/services/admin_sdk"
	"github.com/linuzilla/google-license-manager/services/license_manager"
	"log"
)

type ListCommand struct {
	LicenseManager license_manager.LicenseManager `inject:"*"`
	AdminSdk       admin_sdk.GoogleAdminSdk       `inject:"*"`
	DbBackend      dbBolt.DatabaseBackend         `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*ListCommand)(nil)

func (ListCommand) PostSummerConstruct() {
}

func (ListCommand) ImplementCommandInterface() {
}

func (ListCommand) Command() string {
	return `list`
}

func (cmd *ListCommand) Execute(args ...string) int {
	err := listUsers(cmd.LicenseManager, cmd.AdminSdk, cmd.DbBackend, false)

	if err != nil {
		log.Println(err)
	}

	return 0
}
