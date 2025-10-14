package azuredevops

import (
	"context"
	"encoding/base64"
)

type AuthProviderPAT struct {
	pat string
}

func NewAuthProviderPAT(pat string) AuthProvider {
	return AuthProviderPAT{pat}
}

func (p AuthProviderPAT) GetAuth(_ context.Context) (string, error) {
	auth := "_:" + p.pat
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)), nil
}
