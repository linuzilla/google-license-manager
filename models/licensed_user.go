package models

import dbBolt "github.com/linuzilla/go-boltdb"

type LicensedUser struct {
	Id          string
	Description string
}

var _ dbBolt.BoltModel = (*LicensedUser)(nil)
var licensedUserBucket = []byte("users")

func (model *LicensedUser) PrimaryKey() string {
	return model.Id
}

func (model *LicensedUser) Bucket() []byte {
	return licensedUserBucket
}
