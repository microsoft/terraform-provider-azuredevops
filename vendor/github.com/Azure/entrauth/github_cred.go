package entrauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type GithubCredentialOption struct {
	azcore.ClientOptions
	AdditionallyAllowedTenants []string
	DisableInstanceDiscovery   bool
	Cache                      azidentity.Cache
}

type githubCredential struct {
	requestUrl   string
	requestToken string
	cred         *azidentity.ClientAssertionCredential
}

func NewGithubCredential(tenantId, clientId, requestUrl, requestToken string, options *GithubCredentialOption) (azcore.TokenCredential, error) {
	if tenantId == "" {
		return nil, errors.New("no tenant ID specified")
	}
	if clientId == "" {
		return nil, errors.New("no client ID specified")
	}
	if requestUrl == "" {
		return nil, errors.New("no request URL specified")
	}
	if requestToken == "" {
		return nil, errors.New("no request token specified")
	}

	if options == nil {
		options = &GithubCredentialOption{}
	}

	g := githubCredential{
		requestUrl:   requestUrl,
		requestToken: requestToken,
	}

	copt := azidentity.ClientAssertionCredentialOptions{
		AdditionallyAllowedTenants: options.AdditionallyAllowedTenants,
		Cache:                      options.Cache,
		ClientOptions:              options.ClientOptions,
		DisableInstanceDiscovery:   options.DisableInstanceDiscovery,
	}
	cred, err := azidentity.NewClientAssertionCredential(tenantId, clientId, g.getAssertion, &copt)
	if err != nil {
		return nil, err
	}
	g.cred = cred
	return &g, nil
}

func (g *githubCredential) getAssertion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.requestUrl, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("getAssertion: failed to build request")
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot parse URL query: %v", err)
	}

	if query.Get("audience") == "" {
		query.Set("audience", "api://AzureADTokenExchange")
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.requestToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot request token: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot parse response: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return "", fmt.Errorf("getAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenRes struct {
		Value *string `json:"value"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return "", fmt.Errorf("getAssertion: cannot unmarshal response: %v", err)
	}

	if tokenRes.Value == nil {
		return "", fmt.Errorf("getAssertion: nil JWT assertion received from Github")
	}

	return *tokenRes.Value, nil
}

func (g *githubCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	// NOTE: There is no trace available for this credential
	return g.cred.GetToken(ctx, options)
}
