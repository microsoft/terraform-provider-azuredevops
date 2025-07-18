package entrauth

import (
	"crypto"
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type CredentialOption interface {
	isCredentialOption()
}

type ClientSecretCredentialOption struct {
	TenantId     string
	ClientId     string
	ClientSecret string

	// Optional
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

func (ClientSecretCredentialOption) isCredentialOption() {}

type ClientCertificateCredentialOption struct {
	TenantId string
	ClientId string
	CertData []*x509.Certificate
	CertKey  crypto.PrivateKey

	// Optional
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
	SendCertificateChain       bool
}

func (ClientCertificateCredentialOption) isCredentialOption() {}

type AssertionPlainCredentialOption struct {
	TenantId  string
	ClientId  string
	Assertion string

	// Optional
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

func (AssertionPlainCredentialOption) isCredentialOption() {}

type AssertionFileCredentialOption struct {
	TenantId      string
	ClientId      string
	AssertionFile string

	// Optional
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

func (AssertionFileCredentialOption) isCredentialOption() {}

type AssertionRequestType string

const (
	AssertionRequestTypeGithub      AssertionRequestType = "Github"
	AssertionRequestTypeAzureDevOps AssertionRequestType = "AzureDevOps"
)

type AssertionRequestCredentialOption struct {
	Type           AssertionRequestType
	PlatformOption AssertionRequestCredentialPlatformOption
}

func (AssertionRequestCredentialOption) isCredentialOption() {}

type AssertionRequestCredentialPlatformOption interface {
	isAssertionRequestCredentialPlatformOption()
}

type AssertionRequestGithubCredentialOption struct {
	TenantId     string
	ClientId     string
	RequestToken string
	RequestUrl   string

	// Optional
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

func (AssertionRequestGithubCredentialOption) isAssertionRequestCredentialPlatformOption() {}

type AssertionRequestAzureDevOpsCredentialOption struct {
	TenantId            string
	ClientId            string
	ServiceConnectionId string
	SystemAccessToken   string

	// Optional
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

func (AssertionRequestAzureDevOpsCredentialOption) isAssertionRequestCredentialPlatformOption() {}

type ManagedIdentityCredentialOption struct {
	// Optional
	azcore.ClientOptions
	ID azidentity.ManagedIDKind
}

func (ManagedIdentityCredentialOption) isCredentialOption() {}

type AzureCLICredentialOption struct {
	// Optional
	TenantId                   string
	SubscriptionId             string
	AdditionallyAllowedTenants []string
}

func (AzureCLICredentialOption) isCredentialOption() {}

type AzureDevCLICredentialOption struct {
	// Optional
	TenantId                   string
	AdditionallyAllowedTenants []string
}

func (AzureDevCLICredentialOption) isCredentialOption() {}
