package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test-inspector/internal/supabase/tables"
	"test-inspector/internal/supabase/tables/launch"
	"test-inspector/internal/supabase/tables/result"
	"test-inspector/internal/supabase/tables/version"
	"test-inspector/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/supabase/postgrest-go"
)

const (
	// RestEndpoint is the endpoint for the PostgREST API
	RestEndpoint = "rest/v1"
)

// IClient is an interface to communicate with supabase project.
// @property GetVersion - This is used to get the version of the test that is being run.
// @property CreateLaunch - Creates a new launch in the database.
// @property CreateResult - This is the method that will be called when a test result is
// created.
// @property GetTemplate - This is the method that will be called to get the template for the test.
// @property GetFeatures - Returns a list of features that are available to be tested.
type IClient interface {
	GetVersion(id int32) (int32, error)
	CreateLaunch(l models.Launch) (int64, error)
	CreateResult(r models.SupaResult) error
	GetTemplate(versionID int64) ([]models.SupaResult, error)
	GetFeatures() ([]string, error)
}

// Client is a supabase client struct
// @property {string} baseURL - The base URL of the supabase project API.
// @property {string} apiKey - The API key that you will use to authenticate to use the API.
// @property httpClient - The http client used to make requests to the API.
// @property auth - Is used to make auth call to GoTrue API.
// @property DB - This is a pointer to a postgrest.Client. This is the client that will be used to
// make requests to the PostgREST API.
// @property {User} user - The user that is currently logged in.
type Client struct {
	baseURL string
	// apiKey can be a client API key or a service key
	apiKey     string
	httpClient *http.Client
	auth       *Auth
	DB         *postgrest.Client
	user       User
}

// CreateClient creates a new Supabase client
func CreateClient(baseURL string, supabaseKey string, user UserCredentials) (IClient, error) {
	parsedURL := fmt.Sprintf("%s/%s/", baseURL, RestEndpoint)
	client := &Client{
		baseURL: baseURL,
		apiKey:  supabaseKey,
		auth:    &Auth{},
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
		DB: postgrest.NewClient(parsedURL, "public", map[string]string{"apiKey": supabaseKey}),
	}
	client.auth.client = client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var details *AuthenticatedDetails
	var err error
	if user.Email != "" {
		details, err = client.auth.SignIn(ctx, user)
	} else {
		details = &AuthenticatedDetails{
			AccessToken: supabaseKey,
			User:        User{},
		}
	}
	if err != nil {
		return nil, err
	}
	client.DB = postgrest.NewClient(parsedURL, "public", map[string]string{
		"apiKey":        supabaseKey,
		"Authorization": fmt.Sprintf("Bearer %s", details.AccessToken),
	})
	client.user = details.User
	return client, nil
}

// GetVersion checking if the version exists in the database.
func (c *Client) GetVersion(id int32) (int32, error) {
	var ids []struct {
		ID int32 `json:"id"`
	}
	_, err := c.DB.
		From(tables.Versions.String()).
		Select(version.ID.String(), "1", false).
		Eq(version.ID.String(), strconv.Itoa(int(id))).
		ExecuteTo(&ids)
	if err != nil {
		return 0, err
	}
	if len(ids) != 1 {
		return 0, fmt.Errorf("version with ID: '%d' was not found", id)
	}
	return ids[0].ID, nil
}

// CreateLaunch adds a new launch in the database.
func (c *Client) CreateLaunch(l models.Launch) (int64, error) {
	var ids []struct {
		ID int64 `json:"id"`
	}
	if l.Name == "" {
		l.Name = uuid.NewString()
	}
	l.UserID = &c.user.ID
	_, err := c.DB.
		From(tables.Launches.String()).
		Select(launch.ID.String(), "1", false).
		Eq(launch.Name.String(), l.Name).
		ExecuteTo(&ids)
	if err != nil {
		return 0, err
	}
	if len(ids) == 1 {
		return 0, fmt.Errorf("launch with that name already exists, id=%d", ids[0].ID)
	}
	// remove template flag from previous template
	if l.IsTemplate {
		lastTemplates := []models.Launch{}
		// todo add filter by project
		_, err = c.DB.
			From(tables.Launches.String()).
			Select("*", "1", false).
			Eq(launch.IsTemplate.String(), "true").
			ExecuteTo(&lastTemplates)
		if err != nil {
			return 0, err
		}
		// todo add filter by project and fix all if smth wrong
		if len(lastTemplates) > 0 {
			prevTemplate := lastTemplates[0]
			prevTemplate.IsTemplate = false
			_, _, err = c.DB.
				From(tables.Launches.String()).
				Update(prevTemplate, launch.ID.String(), "1").
				Eq(launch.ID.String(), strconv.Itoa(int(*prevTemplate.ID))).
				Execute()
			if err != nil {
				return 0, err
			}
		}
	}
	_, err = c.DB.From(tables.Launches.String()).
		Insert(l, false, "", "representation", "exact").
		ExecuteTo(&ids)
	if err != nil {
		return 0, err
	}
	return ids[0].ID, nil
}

// CreateResult adds a new result in the database.
func (c *Client) CreateResult(r models.SupaResult) error {
	var ids []struct {
		ID uuid.UUID `json:"uuid"`
	}

	_, err := c.DB.
		From(tables.Results.String()).
		Select(result.ID.String(), "1", false).
		Eq(result.Name.String(), r.Name).
		Eq(result.Suite.String(), r.Suite).
		Eq(result.ParentSuite.String(), r.ParentSuite).
		Eq(result.Feature.String(), r.Feature).
		Eq(result.LaunchID.String(), strconv.Itoa(int(r.LaunchID))).
		ExecuteTo(&ids)
	if err != nil {
		return err
	}
	if len(ids) == 1 {
		fmt.Printf("result with that name already exists, id=%d", ids[0].ID)
		return nil
	}

	_, err = c.DB.From(tables.Results.String()).
		Insert(r, false, "", "representation", "exact").
		ExecuteTo(&ids)
	if err != nil {
		return err
	}

	if len(ids) != 1 {
		return fmt.Errorf("result %s was not inserted, smth went wrong", r.ID.String())
	}
	r.ID = ids[0].ID
	return nil
}

// GetTemplate getting the template results from the database for the project of the given version.
func (c *Client) GetTemplate(versionID int64) ([]models.SupaResult, error) {
	var rpcBody struct {
		Version int64 `json:"version"`
	}
	rpcBody.Version = versionID
	templatesRaw := c.DB.Rpc("template", "", rpcBody)
	if len(templatesRaw) == 0 {
		return nil, fmt.Errorf("no template results found")
	}

	var templates []models.SupaResult
	err := json.Unmarshal([]byte(templatesRaw), &templates)
	if err != nil {
		return nil, fmt.Errorf("error parsing template results: %+v", err)
	}
	if len(templates) == 0 {
		return nil, fmt.Errorf("no template results found")
	}

	return templates, nil
}

// GetFeatures getting all features from the database for the latest reference run.
func (c *Client) GetFeatures() ([]string, error) {
	featuresRaw := c.DB.Rpc("features", "", nil)
	if len(featuresRaw) == 0 {
		return nil, fmt.Errorf("no features found")
	}

	var features []string
	err := json.Unmarshal([]byte(featuresRaw), &features)
	if err != nil {
		return nil, fmt.Errorf("error parsing features: %+v", err)
	}
	if len(features) == 0 {
		return nil, fmt.Errorf("no features found")
	}

	return features, nil
}
