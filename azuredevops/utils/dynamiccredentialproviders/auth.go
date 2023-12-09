package dynamiccredentialproviders

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

func getGitHubOIDCToken(d *schema.ResourceData) (string, error) {
	requestUrl := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	client := &http.Client{}
	audience := "api://AzureADTokenExchange"

	if userAudience, ok := d.GetOk("oidc_github_actions_audience"); ok {
		audience = userAudience.(string)
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}
	query := parsedUrl.Query()
	query.Add("audience", audience)
	parsedUrl.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+requestToken)
	req.Header.Add("Accept", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	response_interface := GHIdTokenResponse{}
	err = json.NewDecoder(response.Body).Decode(&response_interface)
	if err != nil {
		return "", err
	}

	return response_interface.Value, nil
}

type TokenGetter interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type AzIdentityFuncs interface {
	NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (TokenGetter, error)
	NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (TokenGetter, error)
	NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenGetter, error)
}

type AzIdentityFuncsReal struct{}

func (a AzIdentityFuncsReal) NewClientAssertionCredential(tenantID, clientID string, getAssertion func(context.Context) (string, error), options *azidentity.ClientAssertionCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientAssertionCredential(tenantID, clientID, getAssertion, options)
}

func (a AzIdentityFuncsReal) NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *azidentity.ClientCertificateCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, options)
}

func (a AzIdentityFuncsReal) NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenGetter, error) {
	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}

func GetAuthToken(ctx context.Context, d *schema.ResourceData, azIdentityFuncs AzIdentityFuncs) (func() (string, error), error) {
	tenantId := d.Get("tenant_id").(string)
	clientId := d.Get("client_id").(string)
	AzureDevOpsAppDefaultScope := "499b84ac-1321-427f-aa17-267ca6975798/.default"
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{AzureDevOpsAppDefaultScope},
	}

	var cred TokenGetter
	var err error

	// OIDC Token
	if oidc_token, ok := d.GetOk("oidc_token"); ok {
		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return oidc_token.(string), nil }, nil)
		if err != nil {
			return nil, err
		}
	}

	// OIDC Token From File
	if oidc_token_path, ok := d.GetOk("oidc_token_path"); ok {
		fileBytes, err := os.ReadFile(oidc_token_path.(string))
		if err != nil {
			return nil, err
		}
		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return strings.TrimSpace(string(fileBytes)), nil }, nil)
		if err != nil {
			return nil, err
		}
	}

	// OIDC Token in a GitHub Action Workflow
	if oidc_github_actions, ok := d.GetOk("oidc_github_actions"); ok && oidc_github_actions.(bool) {
		gitHubToken, err := getGitHubOIDCToken(d)
		if err != nil {
			return nil, err
		}
		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return gitHubToken, nil }, nil)
		if err != nil {
			return nil, err
		}
	}

	// OIDC Token in a HashiCorp Vault run
	if oidc_hcp, ok := d.GetOk("oidc_hcp"); ok && oidc_hcp.(bool) {
		workloadIdentityToken := os.Getenv("TFC_WORKLOAD_IDENTITY_TOKEN")

		// Check if plan & apply phases use different service principals
		if clientIdPlan, ok := d.GetOk("client_id_plan"); ok {
			clientIdApply := d.Get("client_id_apply").(string)
			tenantIdPlan := d.Get("tenant_id_plan").(string)
			tenantIdApply := d.Get("tenant_id_apply").(string)

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
				clientId = clientIdApply
				tenantId = tenantIdApply
			} else if strings.EqualFold(workloadIdentityTokenUnmarshalled.RunPhase, "plan") {
				clientId = clientIdPlan.(string)
				tenantId = tenantIdPlan
			} else {
				return nil, errors.New(fmt.Sprintf("Unrecognized workspace run phase: %s", workloadIdentityTokenUnmarshalled.RunPhase))
			}
		} else if clientId == "" {
			return nil, errors.New(fmt.Sprintf("Either client_id or client_id_plan must be set when using Terraform Cloud Workload Identity Token authentication."))
		}

		cred, err = azIdentityFuncs.NewClientAssertionCredential(tenantId, clientId, func(context.Context) (string, error) { return workloadIdentityToken, nil }, nil)
		if err != nil {
			return nil, err
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

		cred, err = azIdentityFuncs.NewClientCertificateCredential(tenantId, clientId, certs, key, nil)
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
		cred, err = azIdentityFuncs.NewClientCertificateCredential(tenantId, clientId, certs, key, nil)
		if err != nil {
			return nil, err
		}
	}

	// Client Secret
	if client_secret, ok := d.GetOk("client_secret"); ok {
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantId, clientId, client_secret.(string), nil)
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
		cred, err = azIdentityFuncs.NewClientSecretCredential(tenantId, clientId, strings.TrimSpace(string(fileBytes)), nil)
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
	ctx context.Context
	cred TokenGetter
	opts policy.TokenRequestOptions
	cachedToken *azcore.AccessToken
}

func newAzTokenProvider(cred TokenGetter, ctx context.Context, opts policy.TokenRequestOptions) *AzTokenProvider {
	return &AzTokenProvider{
		cred: cred,
		ctx: ctx,
		opts: opts,
		cachedToken: nil,
	}
}

func (provider *AzTokenProvider) GetToken() (string, error) {
	if provider.cachedToken == nil || provider.cachedToken.ExpiresOn.Before(time.Now().Local().Add(-5 * time.Minute)) {
		cachedToken, err := provider.cred.GetToken(provider.ctx, provider.opts)
		provider.cachedToken = &cachedToken
		if err != nil {
			return "", err
		}
	}
	return "Bearer " + provider.cachedToken.Token, nil
}
