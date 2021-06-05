package models

import (
	dbBolt "github.com/linuzilla/go-boltdb"
)

const (
	DataKeyMasterKey  = `master-key`
	DataKeyGoogle     = `google`
	DataKeyCredential = `credential`
)

type EncryptedData struct {
	Id         string
	Checksum   string
	Base64Data string
}

var _ dbBolt.BoltModel = (*EncryptedData)(nil)
var encryptedData = []byte("data")

func (model *EncryptedData) PrimaryKey() string {
	return model.Id
}

func (model *EncryptedData) Bucket() []byte {
	return encryptedData
}

//func (model *EncryptedData) Decode(password string) (interface{}, error) {
//	return utils.DecodeInDatabaseData(model, password)
//}
