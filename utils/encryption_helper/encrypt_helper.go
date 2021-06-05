package encryption_helper

import (
	"encoding/json"
	"errors"
	"fmt"
	dbBolt "github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/models"
	"github.com/linuzilla/google-license-manager/utils"
	"io/ioutil"
)

func Decode(model *models.EncryptedData, password string) (interface{}, error) {
	decrypted, err := utils.Decrypt(password, model.Base64Data)

	if err != nil {
		return nil, err
	}

	if utils.Sha256String(decrypted) != model.Checksum {
		return nil, errors.New("password incorrect")
	}

	switch model.Id {
	case models.DataKeyGoogle:
		var store config.GoogleCfg
		err := json.Unmarshal(decrypted, &store)

		if err != nil {
			return nil, err
		}
		return &store, nil

	case models.DataKeyCredential:
		return string(decrypted), nil

	case models.DataKeyMasterKey:
		return decrypted, nil

	default:
		return nil, errors.New(model.Id + ": unknown Data")
	}
}

func EncodeMasterKey(masterKey []byte, password string) (*models.EncryptedData, error) {
	encrypt, err := utils.Encrypt(password, masterKey)
	if err != nil {
		return nil, err
	}

	return &models.EncryptedData{
		Id:         models.DataKeyMasterKey,
		Checksum:   utils.Sha256String(masterKey),
		Base64Data: encrypt,
	}, nil
}

func EncodeGoogleConfig(cfg *config.GoogleCfg, password string) (*models.EncryptedData, error) {
	jsonBlob, err := json.Marshal(*cfg)
	if err != nil {
		return nil, err
	}

	encrypt, err := utils.Encrypt(password, jsonBlob)
	if err != nil {
		return nil, err
	}

	return &models.EncryptedData{
		Id:         models.DataKeyGoogle,
		Checksum:   utils.Sha256String(jsonBlob),
		Base64Data: encrypt,
	}, nil
}

func EncodeCredential(data []byte, password string) (*models.EncryptedData, error) {
	encrypt, err := utils.Encrypt(password, data)
	if err != nil {
		return nil, err
	}

	return &models.EncryptedData{
		Id:         models.DataKeyCredential,
		Checksum:   utils.Sha256String(data),
		Base64Data: encrypt,
	}, nil
}

func LoadMasterKey(connection dbBolt.DatabaseBackendConnection, password string) (string, error) {
	var encryptedData models.EncryptedData

	err := connection.FindById(models.DataKeyMasterKey, &encryptedData)
	if err != nil {
		return ``, err
	}

	decode, err := Decode(&encryptedData, password)
	if err != nil {
		return ``, err
	}
	return string(decode.([]byte)), nil
}

func SaveMasterKey(connection dbBolt.DatabaseBackendConnection, masterKey string, password string) error {
	encryptedData, err := EncodeMasterKey([]byte(masterKey), password)

	if err != nil {
		return err
	}

	return connection.SaveOrUpdate(encryptedData)
}

func SaveGoogleConfig(connection dbBolt.DatabaseBackendConnection, googleConfig *config.GoogleCfg, masterKey string) error {
	encryptedData, err := EncodeGoogleConfig(googleConfig, masterKey)

	if err != nil {
		return err
	}

	return connection.SaveOrUpdate(encryptedData)
}

func SaveCredentialFormFile(connection dbBolt.DatabaseBackendConnection, fileName string, masterKey string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	} else {
		encodeCredential, err := EncodeCredential(b, masterKey)
		if err != nil {
			return err
		}
		err = connection.SaveOrUpdate(encodeCredential)
		return err
	}
}
