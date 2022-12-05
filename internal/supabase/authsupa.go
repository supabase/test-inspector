package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	// AuthEndpoint is used to build the URL for the auth endpoint.
	AuthEndpoint = "auth/v1"
)

type authError struct {
	Message string `json:"message"`
}

// Auth is a struct for authentication in supabase project
// @property client - This is the client that will be used to make requests to Supabase.
type Auth struct {
	client *Client
}

// UserCredentials is a struct that contains two fields, Email and Password, both strings.
// @property {string} Email - The email address of the user.
// @property {string} Password - The password of the user.
type UserCredentials struct {
	Email    string
	Password string
}

// User is a struct that holds supabase user information
// @property {string} ID - The unique identifier for the user.
// @property {string} Aud - The role of the token.
// @property {string} Role - The role of the user.
// @property {string} Email - The email address of the user.
// @property InvitedAt - The time the user was invited to the application.
// @property ConfirmedAt - The date and time the user confirmed their email address.
// @property ConfirmationSentAt - The time the user was invited to the application.
// @property AppMetadata - This is where you can store any custom data you want to associate with
// the user.
// @property UserMetadata - This is a map of key-value pairs that you can use to store any
// additional information about the user.
// @property CreatedAt - The date and time the user was created.
// @property UpdatedAt - The date and time the user was last updated.
type User struct {
	ID                 string                    `json:"id"`
	Aud                string                    `json:"aud"`
	Role               string                    `json:"role"`
	Email              string                    `json:"email"`
	InvitedAt          time.Time                 `json:"invited_at"`
	ConfirmedAt        time.Time                 `json:"confirmed_at"`
	ConfirmationSentAt time.Time                 `json:"confirmation_sent_at"`
	AppMetadata        struct{ provider string } `json:"app_metadata"`
	UserMetadata       map[string]interface{}    `json:"user_metadata"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
}

// AuthenticatedDetails is a struct that holds the authentication details
//
// @property {string} AccessToken - The access token that you'll use to make authenticated requests to
// the API.
// @property {string} TokenType - The type of token.
// @property {int} ExpiresIn - The number of seconds until the access token expires.
// @property {string} RefreshToken - This is the token that you will use to refresh the access token
// when it expires.
// @property {User} User - The user object that was returned from the API.
type AuthenticatedDetails struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type authenticationError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// SignIn enters the user credentials and returns the current user if succeeded.
func (a *Auth) SignIn(ctx context.Context, credentials UserCredentials) (*AuthenticatedDetails, error) {
	reqBody, _ := json.Marshal(credentials)
	reqURL := fmt.Sprintf("%s/%s/token?grant_type=password", a.client.baseURL, AuthEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res := AuthenticatedDetails{}
	errRes := authenticationError{}
	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)
	if err != nil {
		return nil, err
	} else if hasCustomError {
		return nil, fmt.Errorf("%s: %s", errRes.Error, errRes.ErrorDescription)
	}

	return &res, nil
}

// RefreshUser refreshes the user token
func (a *Auth) RefreshUser(ctx context.Context, userToken string, refreshToken string) (*AuthenticatedDetails, error) {
	reqBody, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	reqURL := fmt.Sprintf("%s/%s/token?grant_type=refresh_token", a.client.baseURL, AuthEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	injectAuthorizationHeader(req, userToken)
	req.Header.Set("Content-Type", "application/json")
	res := AuthenticatedDetails{}
	errRes := authenticationError{}
	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)
	if err != nil {
		return nil, err
	} else if hasCustomError {
		return nil, fmt.Errorf("%s: %s", errRes.Error, errRes.ErrorDescription)
	}

	return &res, nil
}

// ProviderSignInOptions is a struct with three fields, `Provider`, `RedirectTo`, and `Scopes`.
//
// @property {string} Provider - The name of the provider you want to use.
// @property {string} RedirectTo - The URL to redirect to after the user signs in.
// @property {[]string} Scopes - A list of scopes that you want to request from the provider.
type ProviderSignInOptions struct {
	Provider   string   `url:"provider"`
	RedirectTo string   `url:"redirect_to"`
	Scopes     []string `url:"scopes"`
}

// ProviderSignInDetails is a struct with two fields, `URL` and `Provider`.
// @property {string} URL - The URL to redirect the user to after they sign in.
// @property {string} Provider - The name of the provider.
type ProviderSignInDetails struct {
	URL      string `json:"url"`
	Provider string `json:"provider"`
}

// SignInWithProvider returns a URL for signing in via OAuth
func (a *Auth) SignInWithProvider(opts ProviderSignInOptions) (*ProviderSignInDetails, error) {
	params, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	details := ProviderSignInDetails{
		URL:      fmt.Sprintf("%s/%s/authorize?%s", a.client.baseURL, AuthEndpoint, params.Encode()),
		Provider: opts.Provider,
	}
	return &details, nil
}

// User retrieves the user information based on the given token
func (a *Auth) User(ctx context.Context, userToken string) (*User, error) {
	reqURL := fmt.Sprintf("%s/%s/user", a.client.baseURL, AuthEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	injectAuthorizationHeader(req, userToken)
	res := User{}
	errRes := authError{}
	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)
	if err != nil {
		return nil, err
	} else if hasCustomError {
		return nil, fmt.Errorf("%s", errRes.Message)
	}

	return &res, nil
}

func injectAuthorizationHeader(req *http.Request, value string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", value))
}

func (c *Client) sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error) {
	req.Header.Set("apikey", c.apiKey)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return true, err
	}

	defer res.Body.Close()
	statusOK := res.StatusCode >= http.StatusOK && res.StatusCode < 300
	if !statusOK {
		if err = json.NewDecoder(res.Body).Decode(&errorValue); err == nil {
			return true, nil
		}

		return false, fmt.Errorf("unknown, status code: %d", res.StatusCode)
	} else if res.StatusCode != http.StatusNoContent {
		if err = json.NewDecoder(res.Body).Decode(&successValue); err != nil {
			return false, err
		}
	}

	return false, nil
}
