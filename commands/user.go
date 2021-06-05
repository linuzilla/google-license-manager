package commands

import (
	"fmt"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/google-license-manager/services/admin_sdk"
	"github.com/linuzilla/google-license-manager/services/google_credential"
)

type UserCommand struct {
	AdminSdk admin_sdk.GoogleAdminSdk `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*UserCommand)(nil)
var _ ConditionalRegister = (*UserCommand)(nil)

func (UserCommand) PostSummerConstruct() {
}

func (UserCommand) ImplementCommandInterface() {
}

func (cmd *UserCommand) CanRegister(googleCredential google_credential.GoogleCredential) bool {
	return googleCredential.HasUserInfo()
}

func (UserCommand) Command() string {
	return `user`
}

func (cmd *UserCommand) Execute(args ...string) int {
	if len(args) == 0 {
		fmt.Println("usage: user user ...")
	} else {
		for _, user := range args {
			if info, err := cmd.AdminSdk.UserInfo(user); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Account          : %s\n", info.PrimaryEmail)
				fmt.Printf("Full Name        : %s\n", info.Name.FullName)
				fmt.Printf("Recover Email    : %s\n", info.RecoveryEmail)
				fmt.Printf("Recover Phone    : %s\n", info.RecoveryPhone)
				fmt.Printf("Suspended        : %v\n", info.Suspended)
				fmt.Printf("Suspension Reason: %s\n", info.SuspensionReason)
				fmt.Printf("Is Admin:        : %v\n", info.IsAdmin)
				fmt.Printf("Creation Time    : %s\n", info.CreationTime)
				fmt.Printf("Last Login Time  : %s\n", info.LastLoginTime)
			}
		}
	}
	return 0
}
