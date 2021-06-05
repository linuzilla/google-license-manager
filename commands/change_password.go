package commands

import (
	"errors"
	"fmt"
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/utils"
	"github.com/linuzilla/google-license-manager/utils/encryption_helper"
	"log"
)

type ChangePasswordCommand struct {
	DbBackend    dbBolt.DatabaseBackend `inject:"*"`
	GoogleConfig *config.GoogleCfg      `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*ChangePasswordCommand)(nil)

func (ChangePasswordCommand) PostSummerConstruct() {
}

func (ChangePasswordCommand) ImplementCommandInterface() {
}

func (ChangePasswordCommand) Command() string {
	return `change-password`
}

func (cmd *ChangePasswordCommand) Execute(args ...string) int {
	password, err := utils.ReadPassword("Enter Password: ")
	if err != nil {
		log.Println(err)
	}

	err = cmd.DbBackend.ConnectionEstablish(func(connection dbBolt.DatabaseBackendConnection) error {
		masterKey, err := encryption_helper.LoadMasterKey(connection, password)
		if err != nil {
			return err
		}

		for i := 0; i < 3; i++ {
			newPassword, err := utils.ReadPassword("Enter new password: ")
			if err != nil {
				return err
			}
			retypePassword, err := utils.ReadPassword("Retype new password: ")
			if err != nil {
				return err
			}
			if newPassword == retypePassword {
				return encryption_helper.SaveMasterKey(connection, masterKey, newPassword)
			}
			fmt.Println("Sorry, passwords do not match")
			fmt.Println()
		}
		return errors.New("password unchanged")
	})
	if err != nil {
		log.Println(err)
	}
	return 0
}
