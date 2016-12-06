package baremetal

// APIKey is returned for operations that create or modify user API keys.
import (
	"net/http"
	"time"
)

// See https://docs.us-az-phoenix-1.oracleiaas.com/api/#/en/identity/20160918/ApiKey/
type APIKey struct {
	KeyID        string    `json:"keyId"`
	KeyValue     string    `json:"keyValue"`
	Fingerprint  string    `json:"fingerprint"`
	UserID       string    `json:"userId"`
	TimeCreated  time.Time `json:"timeCreated"`
	TimeModified time.Time `json:"timeModified"`
	State        string    `json:"lifecycleState"`
}

// ListAPIKeyResponses contains a list of API keys
type ListAPIKeyResponses struct {
	ResourceContainer
	Keys []APIKey
}

func (l *ListAPIKeyResponses) GetList() interface{} {
	return &l.Keys
}

// Deletes an API key belonging to a user.
//
// See https://docs.us-az-phoenix-1.oracleiaas.com/api/#/en/identity/20160918/ApiKey/DeleteApiKey
func (c *Client) DeleteAPIKey(userID, fingerprint string, opts *IfMatchOptions) (e error) {
	details := &requestDetails{
		ids:      urlParts{userID, apiKeys, fingerprint},
		name:     resourceUsers,
		optional: opts,
	}

	return c.identityApi.deleteRequest(details)
}

// ListAPIKeys returns information about a user's API keys.
//
// See https://docs.us-az-phoenix-1.oracleiaas.com/api/#/en/identity/20160918/ApiKey/ListApiKeys
func (c *Client) ListAPIKeys(userID string) (response *ListAPIKeyResponses, e error) {
	details := &requestDetails{
		ids:  urlParts{userID, apiKeys, "/"},
		name: resourceUsers,
	}

	var getResp *requestResponse
	if getResp, e = c.identityApi.getRequest(details); e != nil {
		return
	}

	response = &ListAPIKeyResponses{}
	e = getResp.unmarshal(response)
	return
}

// UploadAPIKey - add an API signing key for user. The key must be an RSA public
// key in pem format.
//
// See https://docs.us-az-phoenix-1.oracleiaas.com/api/#/en/identity/20160918/ApiKey/UploadApiKey
func (c *Client) UploadAPIKey(userID, key string, opts *RetryTokenOptions) (apiKey *APIKey, e error) {
	required := struct {
		Key string `json:"key" url:"-"`
	}{
		Key: key,
	}

	details := &requestDetails{
		ids:      urlParts{userID, apiKeys, "/"},
		name:     resourceUsers,
		optional: opts,
		required: required,
	}

	var resp *requestResponse
	if resp, e = c.identityApi.request(http.MethodPost, details); e != nil {
		return
	}

	apiKey = &APIKey{}
	e = resp.unmarshal(apiKey)
	return
}
