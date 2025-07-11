package azuredevops

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type AADCred interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type AuthProviderAAD struct {
	cred AADCred
	opts policy.TokenRequestOptions
}

func NewAuthProviderAAD(cred AADCred, opts policy.TokenRequestOptions) AuthProvider {
	return AuthProviderAAD{cred, opts}
}

func (p AuthProviderAAD) GetAuth(ctx context.Context) (string, error) {
	token, err := p.cred.GetToken(ctx, p.opts)
	if err != nil {
		return "", fmt.Errorf("failed to get AAD token: %v", err)
	}
	return "Bearer " + token.Token, nil
}
