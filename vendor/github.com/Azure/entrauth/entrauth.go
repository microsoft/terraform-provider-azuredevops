package entrauth

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type NewCredentialOption struct {
	Logger             *log.Logger
	ChainedTokenOption azidentity.ChainedTokenCredentialOptions
}

// NewCredential news a chained token credential. The exact credentials and their orders being chained are determined by the `credOpts`.
func NewCredential(credsOpts []CredentialOption, option *NewCredentialOption) (token *azidentity.ChainedTokenCredential, err error) {
	logger := log.New(io.Discard, "", 0)
	if option.Logger != nil {
		logger = option.Logger
	}

	var creds []azcore.TokenCredential
	for _, option := range credsOpts {
		switch opt := option.(type) {
		case ClientSecretCredentialOption:
			if cred, err := azidentity.NewClientSecretCredential(opt.TenantId, opt.ClientId, opt.ClientSecret,
				&azidentity.ClientSecretCredentialOptions{
					ClientOptions:              opt.ClientOptions,
					AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
					DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
					Cache:                      opt.Cache,
				}); err != nil {
				logger.Printf("failed to new client secret credential: %v", err)
			} else {
				logger.Printf("new client secret credential succeeded")
				creds = append(creds, cred)
			}
		case AssertionPlainCredentialOption:
			if opt.Assertion == "" {
				logger.Printf("failed to new plain assertion credential: empty assertion")
			} else if cred, err := azidentity.NewClientAssertionCredential(opt.TenantId, opt.ClientId,
				func(ctx context.Context) (string, error) { return opt.Assertion, nil },
				&azidentity.ClientAssertionCredentialOptions{
					ClientOptions:              opt.ClientOptions,
					AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
					DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
					Cache:                      opt.Cache,
				}); err != nil {
				logger.Printf("failed to new plain assertion credential: %v", err)
			} else {
				logger.Printf("new plain assertion credential succeeded")
				creds = append(creds, cred)
			}
		case AssertionFileCredentialOption:
			// This might seem to be odd at the first glance. Whilst the WorkloadIdentityCredential effectively
			// (safely) read the client assertion from a file and follow the client assertion flow.
			if cred, err := azidentity.NewWorkloadIdentityCredential(&azidentity.WorkloadIdentityCredentialOptions{
				TenantID:                   opt.TenantId,
				ClientID:                   opt.ClientId,
				TokenFilePath:              opt.AssertionFile,
				ClientOptions:              opt.ClientOptions,
				AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
				DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
				Cache:                      opt.Cache,
			}); err != nil {
				logger.Printf("failed to new file (based) assertion credential: %v", err)
			} else {
				logger.Printf("new file (based) assertion credential succeeded")
				creds = append(creds, cred)
			}
		case ClientCertificateCredentialOption:
			if cred, err := azidentity.NewClientCertificateCredential(opt.TenantId, opt.ClientId, opt.CertData, opt.CertKey,
				&azidentity.ClientCertificateCredentialOptions{
					ClientOptions:              opt.ClientOptions,
					AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
					DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
					Cache:                      opt.Cache,
					SendCertificateChain:       opt.SendCertificateChain,
				}); err != nil {
				logger.Printf("failed to new client certificate credential: %v", err)
			} else {
				logger.Printf("new client certificate credential succeeded")
				creds = append(creds, cred)
			}
		case AssertionRequestCredentialOption:
			var (
				cred azcore.TokenCredential
				err  error
			)
			switch opt.Type {
			case AssertionRequestTypeGithub:
				cred, err = newAssertionRequestGithubCredential(opt)
				if err != nil {
					logger.Printf("failed to new request (based) assertion for Github credential: %v", err)
				} else {
					logger.Printf("new request (based) assertion for Github credential succeeded")
					creds = append(creds, cred)
				}
			case AssertionRequestTypeAzureDevOps:
				cred, err = newAssertionRequestAzureDevOpsCredential(opt)
				if err != nil {
					logger.Printf("failed to new request (based) assertion for AzureDevOps credential: %v", err)
				} else {
					logger.Printf("new request (based) assertion for AzureDevOps credential succeeded")
					creds = append(creds, cred)
				}
			default:
				logger.Printf("unknown request (based) assertion credential type: %v", opt.Type)
			}
		case ManagedIdentityCredentialOption:
			if cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
				ClientOptions: opt.ClientOptions,
				ID:            opt.ID,
			}); err != nil {
				logger.Printf("failed to new managed identity credential: %v", err)
			} else {
				logger.Printf("new managed identity credential succeeded")
				creds = append(creds, cred)
			}
		case AzureCLICredentialOption:
			if cred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{
				AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
				TenantID:                   opt.TenantId,
				Subscription:               opt.SubscriptionId,
			}); err != nil {
				logger.Printf("failed to new Azure CLI credential: %v", err)
			} else {
				logger.Printf("new Azure CLI credential succeeded")
				creds = append(creds, cred)
			}
		case AzureDevCLICredentialOption:
			if cred, err := azidentity.NewAzureDeveloperCLICredential(&azidentity.AzureDeveloperCLICredentialOptions{
				AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
				TenantID:                   opt.TenantId,
			}); err != nil {
				logger.Printf("failed to new Azure Dev CLI credential: %v", err)
			} else {
				logger.Printf("new Azure Dev CLI credential succeeded")
				creds = append(creds, cred)
			}
		default:
			return nil, fmt.Errorf("unexpected entrauth.CredentialOption: %T", opt)
		}
	}

	return azidentity.NewChainedTokenCredential(creds, &option.ChainedTokenOption)
}

func newAssertionRequestAzureDevOpsCredential(opt AssertionRequestCredentialOption) (azcore.TokenCredential, error) {
	if opt.Type != AssertionRequestTypeAzureDevOps {
		return nil, fmt.Errorf("invalid type %s (expect %s)", opt.Type, AssertionRequestTypeAzureDevOps)
	}
	popt, ok := opt.PlatformOption.(AssertionRequestAzureDevOpsCredentialOption)
	if !ok {
		return nil, fmt.Errorf("invalid platform option %T (expect %T)", opt.PlatformOption, AssertionRequestAzureDevOpsCredentialOption{})
	}
	return azidentity.NewAzurePipelinesCredential(popt.TenantId, popt.ClientId, popt.ServiceConnectionId, popt.SystemAccessToken,
		&azidentity.AzurePipelinesCredentialOptions{
			ClientOptions:              popt.ClientOptions,
			AdditionallyAllowedTenants: popt.AdditionallyAllowedTenants,
			DisableInstanceDiscovery:   popt.DisableInstanceDiscovery,
			Cache:                      popt.Cache,
		},
	)
}

func newAssertionRequestGithubCredential(opt AssertionRequestCredentialOption) (azcore.TokenCredential, error) {
	if opt.Type != AssertionRequestTypeGithub {
		return nil, fmt.Errorf("invalid type %s (expect %s)", opt.Type, AssertionRequestTypeGithub)
	}
	popt, ok := opt.PlatformOption.(AssertionRequestGithubCredentialOption)
	if !ok {
		return nil, fmt.Errorf("invalid platform option %T (expect %T)", opt.PlatformOption, AssertionRequestGithubCredentialOption{})
	}
	return NewGithubCredential(popt.TenantId, popt.ClientId, popt.RequestUrl, popt.RequestToken,
		&GithubCredentialOption{
			ClientOptions:              popt.ClientOptions,
			AdditionallyAllowedTenants: popt.AdditionallyAllowedTenants,
			Cache:                      popt.Cache,
		},
	)
}
