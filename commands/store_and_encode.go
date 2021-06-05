package commands

import (
	"fmt"
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/utils"
	"github.com/linuzilla/google-license-manager/utils/encryption_helper"
	"log"
	"strings"
)

type StoreAndEncodeCommand struct {
	DbBackend    dbBolt.DatabaseBackend `inject:"*"`
	GoogleConfig *config.GoogleCfg      `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*StoreAndEncodeCommand)(nil)

func (StoreAndEncodeCommand) PostSummerConstruct() {
}

func (StoreAndEncodeCommand) ImplementCommandInterface() {
}

func (StoreAndEncodeCommand) Command() string {
	return `store-and-encode`
}

func (cmd *StoreAndEncodeCommand) Execute(args ...string) int {
	if len(args) < 1 {
		fmt.Println("usage: store-and-encode password")
	} else {
		password := strings.TrimSpace(strings.Join(args, " "))
		fmt.Printf("Encode with password: [%s]\n", password)

		masterKey := utils.RandomString(64, 72)

		err := cmd.DbBackend.ConnectionEstablish(func(connection dbBolt.DatabaseBackendConnection) error {
			err := encryption_helper.SaveGoogleConfig(connection, cmd.GoogleConfig, masterKey)
			if err != nil {
				return err
			}

			err = encryption_helper.SaveCredentialFormFile(connection, cmd.GoogleConfig.CredentialFile, masterKey)
			if err != nil {
				return err
			}

			return encryption_helper.SaveMasterKey(connection, masterKey, password)
		})
		if err != nil {
			log.Println(err)
		}
	}

	return 0
}
