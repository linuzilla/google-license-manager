package commands

import (
	"fmt"
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/models"
	"github.com/linuzilla/google-license-manager/services/google_credential"
	"log"
	"strings"
)

type DescribeCommand struct {
	DbBackend         dbBolt.DatabaseBackend             `inject:"*"`
	CredentialService google_credential.GoogleCredential `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DescribeCommand)(nil)

func (DescribeCommand) PostSummerConstruct() {
}

func (DescribeCommand) ImplementCommandInterface() {
}

func (DescribeCommand) Command() string {
	return `describe`
}

func (cmd *DescribeCommand) Execute(args ...string) int {
	if len(args) < 2 {
		fmt.Println("usage: describe user description ...")
	} else {
		userId := cmd.CredentialService.NormalizeUserId(args[0])

		err := cmd.DbBackend.ConnectionEstablish(func(connection dbBolt.DatabaseBackendConnection) error {
			var entry models.LicensedUser

			if err := connection.FindById(userId, &entry); err != nil {
				return err
			} else {
				entry.Description = strings.Join(args[1:], " ")
				return connection.SaveOrUpdate(&entry)
			}
		})

		if err != nil {
			log.Println(err)
		}
	}
	return 0
}
