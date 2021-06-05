package commands

import (
	"fmt"
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/models"
	"log"
)

type DumpCommand struct {
	DbBackend dbBolt.DatabaseBackend `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DumpCommand)(nil)

func (DumpCommand) PostSummerConstruct() {
}

func (DumpCommand) ImplementCommandInterface() {
}

func (DumpCommand) Command() string {
	return `dump-database`
}

func (cmd *DumpCommand) Execute(args ...string) int {
	err := cmd.DbBackend.ConnectionEstablish(func(connection dbBolt.DatabaseBackendConnection) error {
		var users []models.LicensedUser

		err := connection.FindAll(&users)

		if err == nil {
			fmt.Printf("\nIn database: %d record(s)\n\n", len(users))
			for i, user := range users {
				fmt.Printf("%3d. %s (%s)\n", i+1, user.Id, user.Description)
			}
		}

		return err
	})

	if err != nil {
		log.Println(err)
	}
	return 0
}
