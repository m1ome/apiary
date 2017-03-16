package apiary

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	ApiaryAPIURL = "https://api.apiary.io/"
)

const (
	apiaryActionMe               = "me"
	apiaryActionGetApis          = "me/apis"
	apiaryActionGetTeamApis      = "me/teams/%s/apis"
	apiaryActionFetchBlueprint   = "blueprint/get/%s"
	apiaryActionPublishBlueprint = "blueprint/publish/%s"
)

type ApiaryMeResponse struct {
	ID    string `json:"userId"`
	Name  string `json:"userName"`
	URL   string `json:"userApisUrl"`
	Teams []struct {
		ID   string `json:"teamId"`
		Name string `json:"teamName"`
		URL  string `json:"teamApisUrl"`
	}
}

type ApiaryApisResponse struct {
	Apis []ApiaryApiResponse `json:"apis"`
}

type ApiaryApiResponse struct {
	Name             string `json:"apiName"`
	DocumentationURL string `json:"apiDocumentationUrl"`
	Subdomain        string `json:"apiSubdomain"`
	Private          bool   `json:"apiIsPrivate"`
	Public           bool   `json:"apiIsPublic"`
	Team             bool   `json:"apiIsTeam"`
	Personal         bool   `json:"apiIsPersonal"`
}

type ApiaryFetchResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type apiary struct {
	options ApiaryOptions
	client  *http.Client
}

type ApiaryOptions struct {
	Token string
}

// NewApiary create new apiary.io client
func NewApiary(opts ApiaryOptions) *apiary {
	return &apiary{
		options: opts,
		client:  &http.Client{},
	}
}

// Me retrieve user information
//
// Reference: http://docs.apiary.apiary.io/#reference/user-information/me/get-me
func (a *apiary) Me() (me ApiaryMeResponse, err error) {
	data, response, err := a.sendRequest(apiaryActionMe)
	if err != nil {
		return
	}

	err = checkOk(response)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &me)
	if err != nil {
		return
	}

	return
}

// GetApis return list of user blueprints/APIs
//
// Reference: http://docs.apiary.apiary.io/#reference/api-list/user-api-list/get-me
func (a *apiary) GetApis() (apis *ApiaryApisResponse, err error) {
	data, response, err := a.sendRequest(apiaryActionGetApis)
	if err != nil {
		return
	}

	err = checkOk(response)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &apis)
	if err != nil {
		return
	}

	return
}

// GetTeamApis return list of team blueprints/APIs
//
// Reference: http://docs.apiary.apiary.io/#reference/api-list/team-api-list/get-me
func (a *apiary) GetTeamApis(team string) (apis *ApiaryApisResponse, err error) {
	uri := fmt.Sprintf(apiaryActionGetTeamApis, team)
	data, response, err := a.sendRequest(uri)
	if err != nil {
		return
	}

	err = checkOk(response)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &apis)
	if err != nil {
		return
	}

	return
}

// PublishBlueprint publish blueprint in apiary.io
//
// Reference: http://docs.apiary.apiary.io/#reference/blueprint/publish-blueprint/get-me
func (a *apiary) PublishBlueprint(name string, content []byte) (published bool, err error) {
	jsonData, err := json.Marshal(map[string]string{
		"code": string(content),
	})

	if err != nil {
		return
	}

	uri := fmt.Sprintf(apiaryActionPublishBlueprint, name)
	data, response, err := a.sendLegacyPostRequest(uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusCreated {
		var apiaryError struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		err = json.Unmarshal(data, &apiaryError)
		if err != nil {
			return
		}

		if apiaryError.Error {
			err = errors.New(fmt.Sprintf("Creation failed: %s", apiaryError.Message))
			return
		}
	}

	published = true

	return
}

// FetchBlueprint fetches blueprint from apiary.io
//
// Reference: Unknown
func (a *apiary) FetchBlueprint(name string) (blueprint *ApiaryFetchResponse, err error) {
	uri := fmt.Sprintf(apiaryActionFetchBlueprint, name)
	data, response, err := a.sendLegacyRequest(uri)
	if err != nil {
		return
	}

	err = checkOk(response)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &blueprint)
	if err != nil {
		return
	}

	return
}
