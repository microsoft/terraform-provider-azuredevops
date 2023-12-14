package azuredevops

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
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	creds, err := o.azIdentityFuncs.NewClientSecretCredential(o.tenantID, o.clientID, oidc_response.Value, nil)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	return creds.GetToken(ctx, opts)
}

func GetAuthTokenProvider(ctx context.Context, d *schema.ResourceData, azIdentityFuncs IdentityFuncsI) (func() (string, error), error) {
	// Personal Access Token
	if personal_access_token, ok := d.GetOk("personal_access_token"); ok {
		tokenFunction := func() (string, error) {
			auth := "_:" + personal_access_token.(string)
			return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)), nil
		}
		return tokenFunction, nil
	}

	// Azure Authentication Schemes
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	AzureDevOpsAppDefaultScope := "499b84ac-1321-427f-aa17-267ca6975798/.default"
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{AzureDevOpsAppDefaultScope},
	}

	var cred TokenGetter
	var err error

	if use_oidc, ok := d.GetOk("use_oidc"); ok && use_oidc.(bool) {
		if oidc_token, ok := d.GetOk("oidc_token"); ok {
			// Provided OIDC Token
			cred, err = azIdentityFuncs.NewClientSecretCredential(tenantID, clientID, oidc_token.(string), nil)
			if err != nil {
				return nil, err
			}
		} else if oidc_token_file_path, ok := d.GetOk("oidc_token_file_path"); ok {
			// OIDC Token From File
			fileBytes, err := os.ReadFile(oidc_token_file_path.(string))
			if err != nil {
				return nil, err
			}
			cred, err = azIdentityFuncs.NewClientSecretCredential(tenantID, clientID, strings.TrimSpace(string(fileBytes)), nil)
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
		} else {
			// OIDC Token from Terraform Cloud
			tfc_token_env_var := "TFC_WORKLOAD_IDENTITY_TOKEN"
			if oidc_tfc_tag, ok := d.GetOk("oidc_tfc_tag"); ok && oidc_tfc_tag.(string) != "" {
				tfc_token_env_var = "TFC_WORKLOAD_IDENTITY_TOKEN_" + oidc_tfc_tag.(string)
			}

			workloadIdentityToken := os.Getenv(tfc_token_env_var)
			if workloadIdentityToken == "" {
				return nil, errors.New("No OIDC token found in " + tfc_token_env_var + " environment variable.")
			}

			// Check if plan & apply phases use different service principals
			if clientIdPlan, ok := d.GetOk("client_id_plan"); ok {
				clientIdApply := d.Get("client_id_apply").(string)
				tenantIdPlan := d.Get("tenant_id_plan").(string)
				tenantIdApply := d.Get("tenant_id_apply").(string)

				// Parse which phase we're in from the OIDC token
				workloadIdentityTokenUnmarshalled := HCPWorkloadToken{}
				jwtParts := strings.Split(workloadIdentityToken, ".")
				if len(jwtParts) != 3 {
					return nil, errors.New("Unable to split TFC_WORKLOAD_IDENTITY_TOKEN jwt")
				}
				jwtClaims := jwtParts[1]
				if i := len(jwtClaims) % 4; i != 0 {
					jwtClaims += strings.Repeat("=", 4-i)
				}
				tokenClaims, err := base64.StdEncoding.DecodeString(jwtClaims)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(tokenClaims, &workloadIdentityTokenUnmarshalled)
				if err != nil {
					return nil, err
				}

				if strings.EqualFold(workloadIdentityTokenUnmarshalled.RunPhase, "apply") {
					clientID = clientIdApply
					tenantID = tenantIdApply
				} else if strings.EqualFold(workloadIdentityTokenUnmarshalled.RunPhase, "plan") {
					clientID = clientIdPlan.(string)
					tenantID = tenantIdPlan
				} else {
					return nil, errors.New(fmt.Sprintf("Unrecognized workspace run phase: %s", workloadIdentityTokenUnmarshalled.RunPhase))
				}
			} else if clientID == "" {
				return nil, errors.New(fmt.Sprintf("Either client_id or client_id_plan must be set when using Terraform Cloud Workload Identity Token authentication."))
			}

			cred, err = azIdentityFuncs.NewClientSecretCredential(tenantID, clientID, workloadIdentityToken, nil)
			if err != nil {
				return nil, err
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

	// Client Secret
	if client_secret, ok := d.GetOk("client_secret"); ok {
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantID, clientID, client_secret.(string), nil)
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
		return nil, errors.New(fmt.Sprintf("No valid credentials found."))
	}

	provider := newAzTokenProvider(cred, context.Background(), tokenOptions)
	return provider.GetToken, nil
}

type AzTokenProvider struct {
	ctx         context.Context
	cred        TokenGetter
	opts        policy.TokenRequestOptions
	cachedToken *azcore.AccessToken
}

func newAzTokenProvider(cred TokenGetter, ctx context.Context, opts policy.TokenRequestOptions) *AzTokenProvider {
	return &AzTokenProvider{
		cred:        cred,
		ctx:         ctx,
		opts:        opts,
		cachedToken: nil,
	}
}

func (provider *AzTokenProvider) GetToken() (string, error) {
	if provider.cachedToken == nil || provider.cachedToken.ExpiresOn.Before(time.Now().Local().Add(-5*time.Minute)) {
		cachedToken, err := provider.cred.GetToken(provider.ctx, provider.opts)
		provider.cachedToken = &cachedToken
		if err != nil {
			return "", err
		}
	}
	return "Bearer " + provider.cachedToken.Token, nil
}
