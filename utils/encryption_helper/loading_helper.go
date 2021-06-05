package encryption_helper

import (
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/models"
	"github.com/linuzilla/google-license-manager/utils"
)

func ReadEncryptedData(databaseBackend dbBolt.DatabaseBackend) (googleCfg *config.GoogleCfg, credential string, err error) {
	password := ``
	password, err = utils.ReadPassword("Enter Password: ")

	if err != nil {
		return
	}
	err = databaseBackend.ConnectionEstablish(func(connection dbBolt.DatabaseBackendConnection) error {
		masterKey, err := LoadMasterKey(connection, password)

		if err != nil {
			return err
		}

		var encryptedData models.EncryptedData

		err = connection.FindById(models.DataKeyGoogle, &encryptedData)
		if err != nil {
			return err
		}
		decoded, err := Decode(&encryptedData, masterKey)

		if err != nil {
			return err
		}
		googleCfg = decoded.(*config.GoogleCfg)

		err = connection.FindById(models.DataKeyCredential, &encryptedData)
		if err != nil {
			return err
		}
		decoded, err = Decode(&encryptedData, masterKey)
		if err != nil {
			return err
		}
		credential = decoded.(string)
		return nil
	})
	return
}
