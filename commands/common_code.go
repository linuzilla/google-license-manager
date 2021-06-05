package commands

import (
	"fmt"
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/google-license-manager/models"
	"github.com/linuzilla/google-license-manager/services/admin_sdk"
	"github.com/linuzilla/google-license-manager/services/license_manager"
	"log"
	"sort"
	"strconv"
)

func listUsers(licenseManager license_manager.LicenseManager,
	adminSdk admin_sdk.GoogleAdminSdk,
	dbBackend dbBolt.DatabaseBackend,
	syncToDb bool) error {
	return dbBackend.ConnectionEstablish(func(connection dbBolt.DatabaseBackendConnection) error {
		var users []models.LicensedUser

		err := connection.FindAll(&users)

		if err == nil {
			databaseMap := make(map[string]*models.LicensedUser)

			for i, user := range users {
				databaseMap[user.Id] = &users[i]
			}
			licenseUsers := licenseManager.ListLicenseUser()

			if len(licenseUsers) == 0 {
				fmt.Println("No product/User found.")
			} else {
				fmt.Printf("Product: %s, [ %s ]\n\n", licenseUsers[0].SkuId, licenseUsers[0].ProductName)

				sort.Slice(licenseUsers, func(i, j int) bool {
					return licenseUsers[i].UserId < licenseUsers[j].UserId
				})
				maxLen := 0

				for _, item := range licenseUsers {
					if len(item.UserId) > maxLen {
						maxLen = len(item.UserId)
					}
				}
				format := "%3d. %" + strconv.Itoa(maxLen) + "s (%s) %s\n"

				for i, item := range licenseUsers {
					entry, found := databaseMap[item.UserId]
					tag := ` `
					description := ``

					if found {
						description = entry.Description
						delete(databaseMap, item.UserId)
					} else {
						tag = `+`

						if syncToDb {
							newEntry := models.LicensedUser{
								Id: item.UserId,
							}
							if adminSdk.HasUserInfo() {
								if info, err := adminSdk.UserInfo(item.UserId); err != nil {
									fmt.Println(err)
								} else {
									newEntry.Description = info.Name.FullName
									description = info.Name.FullName
								}
							}
							if err := connection.Persist(&newEntry); err != nil {
								log.Println(err)
							}
						}
					}

					fmt.Printf(format, i+1, item.UserId, description, tag)
				}
			}

			if len(databaseMap) > 0 {
				i := 0
				fmt.Println()
				for id, data := range databaseMap {
					i = i + 1
					fmt.Printf("%3d. %s (%s) - \n", i, id, data.Description)
					if syncToDb {
						if err := connection.Delete(&models.LicensedUser{
							Id: id,
						}); err != nil {
							log.Println(err)
						}
					}
				}
			}
			fmt.Println()
			fmt.Printf("License use: %d\n", len(licenseUsers))
		}

		return err
	})
}
