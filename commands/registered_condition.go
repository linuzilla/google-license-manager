package commands

import "github.com/linuzilla/google-license-manager/services/google_credential"

type ConditionalRegister interface {
	CanRegister(googleCredential google_credential.GoogleCredential) bool
}
