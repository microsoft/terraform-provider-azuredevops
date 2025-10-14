package aztfauth

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/entrauth"
)

type Option struct {
	Logger *log.Logger

	TenantId     string
	TenantIdFile string

	ClientId     string
	ClientIdFile string

	UseClientSecret  bool
	ClientSecret     string
	ClientSecretFile string

	UseClientCert        bool
	ClientCertBase64     string
	ClientCertPfxFile    string
	ClientCertPassword   []byte
	SendCertificateChain bool

	UseOIDCToken bool
	OIDCToken    string

	UseOIDCTokenFile bool
	OIDCTokenFile    string

	UseOIDCTokenRequest    bool
	OIDCRequestToken       string
	OIDCRequestURL         string
	ADOServiceConnectionId string

	UseMSI         bool
	UseAzureCLI    bool
	UseAzureDevCLI bool

	// Common
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

func (Option) fromValueOrFile(name, v, file string) (*string, error) {
	if file == "" {
		return &v, nil
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q from file %q: %v", name, file, err)
	}
	fv := strings.TrimSpace(string(b))
	if v == "" {
		v = fv
		return &v, nil
	}
	if v != fv {
		return nil, fmt.Errorf("mismatch value of %q between the specified value and read value from %q", name, file)
	}
	return &v, nil
}

func (opt Option) getClientId() (*string, error) {
	return opt.fromValueOrFile("Client Id", opt.ClientId, opt.ClientIdFile)
}

func (opt Option) getTenantId() (*string, error) {
	return opt.fromValueOrFile("Tenant Id", opt.TenantId, opt.TenantIdFile)
}

func (opt Option) getClientSecret() (*string, error) {
	return opt.fromValueOrFile("Client Secret", opt.ClientSecret, opt.ClientSecretFile)
}

func (opt Option) getClientCert() ([]*x509.Certificate, crypto.PrivateKey, error) {
	var certData []byte
	if v := opt.ClientCertBase64; v != "" {
		vv, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to base64 decode client certificate")
		}
		certData = vv
	}
	if file := opt.ClientCertPfxFile; file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read client certificate at %q: %v", file, err)
		}

		if len(certData) == 0 {
			certData = b
		} else {
			if string(certData) != string(b) {
				return nil, nil, fmt.Errorf("mismatch value of client certificate between the specified value and read value from %q", file)
			}
		}
	}
	if len(certData) == 0 {
		return nil, nil, fmt.Errorf("no client certificate available")
	}
	certs, key, err := azidentity.ParseCertificates(certData, opt.ClientCertPassword)
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to parse client certificate": %v`, err)
	}

	return certs, key, nil
}

func (opt Option) buildOIDCTokenCredOpt() (entrauth.CredentialOption, error) {
	clientId, err := opt.getClientId()
	if err != nil {
		return nil, err
	}
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}
	return entrauth.AssertionPlainCredentialOption{
		ClientId:  *clientId,
		TenantId:  *tenantId,
		Assertion: opt.OIDCToken,

		ClientOptions:              opt.ClientOptions,
		AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
		DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
		Cache:                      opt.Cache,
	}, nil
}

func (opt Option) buildOIDCTokenFileCredOpt() (entrauth.CredentialOption, error) {
	clientId, err := opt.getClientId()
	if err != nil {
		return nil, err
	}
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}
	return entrauth.AssertionFileCredentialOption{
		ClientId:      *clientId,
		TenantId:      *tenantId,
		AssertionFile: opt.OIDCTokenFile,

		ClientOptions:              opt.ClientOptions,
		AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
		DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
		Cache:                      opt.Cache,
	}, nil
}

func (opt Option) buildOIDCTokenReqCredOpt() (entrauth.CredentialOption, error) {
	clientId, err := opt.getClientId()
	if err != nil {
		return nil, err
	}
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}

	// Based on whether the ADO service connection ID is specified, we choose to use either
	// the ADO flow or Github flow.
	if opt.ADOServiceConnectionId != "" {
		return entrauth.AssertionRequestCredentialOption{
			Type: entrauth.AssertionRequestTypeAzureDevOps,
			PlatformOption: entrauth.AssertionRequestAzureDevOpsCredentialOption{
				ClientId:            *clientId,
				TenantId:            *tenantId,
				ServiceConnectionId: opt.ADOServiceConnectionId,
				SystemAccessToken:   opt.OIDCRequestToken,

				ClientOptions:              opt.ClientOptions,
				AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
				DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
				Cache:                      opt.Cache,
			},
		}, nil
	} else {
		return entrauth.AssertionRequestCredentialOption{
			Type: entrauth.AssertionRequestTypeGithub,
			PlatformOption: entrauth.AssertionRequestGithubCredentialOption{
				ClientId:     *clientId,
				TenantId:     *tenantId,
				RequestUrl:   opt.OIDCRequestURL,
				RequestToken: opt.OIDCRequestToken,

				ClientOptions:              opt.ClientOptions,
				AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
				DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
				Cache:                      opt.Cache,
			},
		}, nil
	}
}

func (opt Option) buildClientSecretCredOpt() (entrauth.CredentialOption, error) {
	clientId, err := opt.getClientId()
	if err != nil {
		return nil, err
	}
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}
	secret, err := opt.getClientSecret()
	if err != nil {
		return nil, err
	}
	return entrauth.ClientSecretCredentialOption{
		TenantId:     *tenantId,
		ClientId:     *clientId,
		ClientSecret: *secret,

		ClientOptions:              opt.ClientOptions,
		AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
		DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
		Cache:                      opt.Cache,
	}, nil
}

func (opt Option) buildClientCertificateCredOpt() (entrauth.CredentialOption, error) {
	clientId, err := opt.getClientId()
	if err != nil {
		return nil, err
	}
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}
	certs, key, err := opt.getClientCert()
	if err != nil {
		return nil, err
	}
	return entrauth.ClientCertificateCredentialOption{
		TenantId: *tenantId,
		ClientId: *clientId,
		CertData: certs,
		CertKey:  key,

		ClientOptions:              opt.ClientOptions,
		AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
		DisableInstanceDiscovery:   opt.DisableInstanceDiscovery,
		Cache:                      opt.Cache,
	}, nil
}

func (opt Option) buildMSICredOpt() (entrauth.CredentialOption, error) {
	clientId, err := opt.getClientId()
	if err != nil {
		return nil, err
	}
	out := entrauth.ManagedIdentityCredentialOption{
		ClientOptions: opt.ClientOptions,
	}
	if *clientId != "" {
		out.ID = azidentity.ClientID(*clientId)
	}
	return out, nil
}

func (opt Option) buildAzureCLICredOpt() (entrauth.CredentialOption, error) {
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}
	return entrauth.AzureCLICredentialOption{
		TenantId:                   *tenantId,
		AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
	}, nil
}

func (opt Option) buildAzureDevCLICredOpt() (entrauth.CredentialOption, error) {
	tenantId, err := opt.getTenantId()
	if err != nil {
		return nil, err
	}
	return entrauth.AzureDevCLICredentialOption{
		TenantId:                   *tenantId,
		AdditionallyAllowedTenants: opt.AdditionallyAllowedTenants,
	}, nil
}

type Credential struct {
	cred *azidentity.ChainedTokenCredential
}

func NewCredential(opt Option) (cred azcore.TokenCredential, err error) {
	logger := log.New(io.Discard, "", 0)
	if opt.Logger != nil {
		logger = opt.Logger
	}

	var credOpts []entrauth.CredentialOption
	if opt.UseOIDCToken {
		if credOpt, err := opt.buildOIDCTokenCredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build oidc token cred option: %v", err)
		}
	}
	if opt.UseOIDCTokenFile {
		if credOpt, err := opt.buildOIDCTokenFileCredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build oidc token file cred option: %v", err)
		}
	}
	if opt.UseOIDCTokenRequest {
		if credOpt, err := opt.buildOIDCTokenReqCredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build oidc token request cred option: %v", err)
		}
	}
	if opt.UseClientSecret {
		if credOpt, err := opt.buildClientSecretCredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build client secret cred option: %v", err)
		}
	}
	if opt.UseClientCert {
		if credOpt, err := opt.buildClientCertificateCredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build client certificate cred option: %v", err)
		}
	}
	if opt.UseMSI {
		if credOpt, err := opt.buildMSICredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build MSI cred option: %v", err)
		}
	}
	if opt.UseAzureCLI {
		if credOpt, err := opt.buildAzureCLICredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build Azure CLI cred option: %v", err)
		}
	}
	if opt.UseAzureDevCLI {
		if credOpt, err := opt.buildAzureDevCLICredOpt(); err == nil {
			credOpts = append(credOpts, credOpt)
		} else {
			logger.Printf("failed to build Azure Dev CLI cred option: %v", err)
		}
	}

	return entrauth.NewCredential(credOpts, &entrauth.NewCredentialOption{
		Logger:             logger,
		ChainedTokenOption: azidentity.ChainedTokenCredentialOptions{RetrySources: true},
	})
}
