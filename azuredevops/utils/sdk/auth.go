package sdk

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

type GHIdTokenResponse struct {
	Value string `json:"value"`
}

type HCPWorkloadToken struct {
	RunPhase string `json:"terraform_run_phase"`
}

type TokenGetter interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type IdentityFuncsI interface {
	NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (TokenGetter, error)
	NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (TokenGetter, error)
	NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenGetter, error)
	NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (TokenGetter, error)
}

type AzIdentityFuncsImpl struct{}

func (a AzIdentityFuncsImpl) NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientAssertionCredential(tenantID, clientID, getAssertion, options)
}

func (a AzIdentityFuncsImpl) NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, options)
}

func (a AzIdentityFuncsImpl) NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}

func (a AzIdentityFuncsImpl) NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (TokenGetter, error) {
	return azidentity.NewManagedIdentityCredential(options)
}

type OIDCCredentialProvder struct {
	audience        string
	clientID        string
	requestToken    string
	requestUrl      string
	tenantID        string
	azIdentityFuncs IdentityFuncsI
}

func (o *OIDCCredentialProvder) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	client := &http.Client{}

	// Assemble the URL with optional audience
	parsedUrl, err := url.Parse(o.requestUrl)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	query := parsedUrl.Query()
	if o.audience != "" {
		query.Add("audience", o.audience)
		parsedUrl.RawQuery = query.Encode()
	}

	// Configure the request
	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	req.Header.Add("Authorization", "Bearer "+o.requestToken)
	req.Header.Add("Accept", "application/json")

	// Make the request
	response, err := client.Do(req)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	// Parse the response
	defer response.Body.Close()
	oidc_response := GHIdTokenResponse{}
	err = json.NewDecoder(response.Body).Decode(&oidc_response)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	// Request the access token from Azure AD using the OIDC token
	creds, err := o.azIdentityFuncs.NewClientAssertionCredential(o.tenantID, o.clientID, AssertionProviderFromString(oidc_response.Value), nil)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	return creds.GetToken(ctx, opts)
}

func GetAuthProvider(ctx context.Context, d *schema.ResourceData, azIdentityFuncs IdentityFuncsI) (azuredevops.AuthProvider, error) {
	// Personal Access Token
	if personal_access_token, ok := d.GetOk("personal_access_token"); ok {
		return azuredevops.NewAuthProviderPAT(personal_access_token.(string)), nil
	}

	// Azure Authentication Schemes
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	var cred TokenGetter
	var err error

	if use_oidc, ok := d.GetOk("use_oidc"); ok && use_oidc.(bool) {
		if oidc_token, ok := d.GetOk("oidc_token"); ok {
			// Provided OIDC Token
			cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantID, clientID, AssertionProviderFromString(oidc_token.(string)), nil)
			if err != nil {
				return nil, err
			}
		} else if oidc_token_file_path, ok := d.GetOk("oidc_token_file_path"); ok {
			// OIDC Token From File
			fileBytes, err := os.ReadFile(oidc_token_file_path.(string))
			if err != nil {
				return nil, err
			}
			cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantID, clientID, AssertionProviderFromString(strings.TrimSpace(string(fileBytes))), nil)
			if err != nil {
				return nil, err
			}
		} else if oidc_request_url, ok := d.GetOk("oidc_request_url"); ok && oidc_request_url.(string) != "" {
			audience := "api://AzureADTokenExchange"
			if oidc_audience, ok := d.GetOk("oidc_audience"); ok && oidc_audience.(string) != "" {
				audience = oidc_audience.(string)
			}

			if _, ok = d.GetOk("oidc_request_token"); !ok {
				return nil, errors.New("No oidc_request_token token found.")
			}

			// OIDC Token from a REST request, ex: Github Action Workflow
			cred = &OIDCCredentialProvder{
				audience:        audience,
				requestUrl:      oidc_request_url.(string),
				requestToken:    d.Get("oidc_request_token").(string),
				tenantID:        tenantID,
				clientID:        clientID,
				azIdentityFuncs: azIdentityFuncs,
			}
		}
	}

	// Certificate from a file on disk
	if client_certificate_path, ok := d.GetOk("client_certificate_path"); ok {
		fileBytes, err := os.ReadFile(client_certificate_path.(string))
		if err != nil {
			return nil, err
		}

		certPassword := ([]byte)(nil)
		if password, ok := d.GetOk("client_certificate_password"); ok {
			certPassword = []byte(password.(string))
		}

		certs, key, err := azidentity.ParseCertificates(fileBytes, certPassword)
		if err != nil {
			return nil, err
		}

		cred, err = azIdentityFuncs.NewClientCertificateCredential(tenantID, clientID, certs, key, nil)
		if err != nil {
			return nil, err
		}
	}

	// Certificate from a base64 encoded string
	if client_certificate, ok := d.GetOk("client_certificate"); ok {
		cert_bytes, err := base64.StdEncoding.DecodeString(client_certificate.(string))
		if err != nil {
			return nil, err
		}
		certPassword := ([]byte)(nil)
		if password, ok := d.GetOk("client_certificate_password"); ok {
			certPassword = []byte(password.(string))
		}
		certs, key, err := azidentity.ParseCertificates(cert_bytes, certPassword)
		if err != nil {
			return nil, err
		}
		cred, err = azIdentityFuncs.NewClientCertificateCredential(tenantID, clientID, certs, key, nil)
		if err != nil {
			return nil, err
		}
	}

	// Client Secret from a file on disk
	if client_secret_path, ok := d.GetOk("client_secret_path"); ok {
		fileBytes, err := os.ReadFile(client_secret_path.(string))
		if err != nil {
			return nil, err
		}
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantID, clientID, strings.TrimSpace(string(fileBytes)), nil)
		if err != nil {
			return nil, err
		}
	}

	// Client Secret
	if client_secret, ok := d.GetOk("client_secret"); ok {
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantID, clientID, client_secret.(string), nil)
		if err != nil {
			return nil, err
		}
	}

	// Azure Managed Service Identity
	if use_msi, ok := d.GetOk("use_msi"); ok && use_msi.(bool) {
		options := &azidentity.ManagedIdentityCredentialOptions{}
		if client_id, ok := d.GetOk("client_id"); ok {
			options.ID = azidentity.ClientID(client_id.(string))
		}

		cred, err = azIdentityFuncs.NewManagedIdentityCredential(options)
		if err != nil {
			return nil, err
		}
	}

	if cred == nil {
		return nil, fmt.Errorf("No valid credentials found.")
	}

	AzureDevOpsAppDefaultScope := "499b84ac-1321-427f-aa17-267ca6975798/.default"
	ap := azuredevops.NewAuthProviderAAD(cred, policy.TokenRequestOptions{
		Scopes: []string{AzureDevOpsAppDefaultScope},
	})
	return ap, nil
}

func AssertionProviderFromString(assertion string) func(context.Context) (string, error) {
	return func(context.Context) (string, error) {
		return assertion, nil
	}
}
