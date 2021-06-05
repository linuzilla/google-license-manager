package commands

import (
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/services/admin_sdk"
	"github.com/linuzilla/google-license-manager/services/license_manager"
	"log"
)

type SyncCommand struct {
	LicenseManager license_manager.LicenseManager `inject:"*"`
	AdminSdk       admin_sdk.GoogleAdminSdk       `inject:"*"`
	DbBackend      dbBolt.DatabaseBackend         `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*SyncCommand)(nil)

func (SyncCommand) PostSummerConstruct() {
}

func (SyncCommand) ImplementCommandInterface() {
}

func (SyncCommand) Command() string {
	return `sync`
}

func (cmd *SyncCommand) Execute(args ...string) int {
	err := listUsers(cmd.LicenseManager, cmd.AdminSdk, cmd.DbBackend, true)

	if err != nil {
		log.Println(err)
	}

	return 0
}
