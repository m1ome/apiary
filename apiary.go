package apiary

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ApiaryAPIURL URL of public apiary.io API
const ApiaryAPIURL = "https://api.apiary.io/"

const (
	apiaryActionMe               = "me"
	apiaryActionGetApis          = "me/apis"
	apiaryActionGetTeamApis      = "me/teams/%s/apis"
	apiaryActionFetchBlueprint   = "blueprint/get/%s"
	apiaryActionPublishBlueprint = "blueprint/publish/%s"
)

// ApiaryMeResponse is a struct of answer to Me() call
//
// Description:
// ID - user id
// Name - user name
// URL - user API URL
// Teams - slice of
// * ID - team id
// * Name - team name
// * URL - team api url
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

// ApiaryApisResponse is a struct of answer to GetApis() all
type ApiaryApisResponse struct {
	Apis []ApiaryApiResponse `json:"apis"`
}

// ApiaryApiResponse is a helper struct of API response: GetApis(), GetTeamApis()
//
// Description:
// Name - API name
// DocumentationURL - URL of docs hosten on apiary.io
// Subdomain - short subdomain (3 level domain)
// Private - is this doc private
// Public - is this doc public
// Team - is this doc belongs to team
// Personal - this this doc personal
type ApiaryApiResponse struct {
	Name             string `json:"apiName"`
	DocumentationURL string `json:"apiDocumentationUrl"`
	Subdomain        string `json:"apiSubdomain"`
	Private          bool   `json:"apiIsPrivate"`
	Public           bool   `json:"apiIsPublic"`
	Team             bool   `json:"apiIsTeam"`
	Personal         bool   `json:"apiIsPersonal"`
}

// ApiaryFetchResponse is a struct of Fetch response
//
// Description:
// Error - is fetch return error
// Message - error message (when error -> false this would be "")
// Code - error code (when error -> false this would be "")
type ApiaryFetchResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Apiary basic API client
//
// Usage:
//package main
//
//import (
//"fmt"
//"log"
//"os"
//
//"github.com/m1ome/apiary"
//)
//
//func main() {
//	token := os.Getenv("APIARY_TOKEN")
//
//	api := NewApiary(ApiaryOptions{
//		Token: Token,
//	})
//
//	response, err := api.Me()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("ID: %d\n", response.ID)
//	fmt.Printf("Name: %s\n", response.Name)
//	fmt.Printf("URL: %s\n", response.URL)
//}
type Apiary struct {
	options ApiaryOptions
	client  *http.Client
}

// ApiaryOptions structure of possible API options
// Token - Your apiary.io token's to access API.
type ApiaryOptions struct {
	Token string
}

// NewApiary create new Apiary.io client
func NewApiary(opts ApiaryOptions) *Apiary {
	return &Apiary{
		options: opts,
		client:  &http.Client{},
	}
}

// Me retrieve user information
//
// Reference: http://docs.apiary.apiary.io/#reference/user-information/me/get-me
func (a *Apiary) Me() (me ApiaryMeResponse, err error) {
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
func (a *Apiary) GetApis() (apis *ApiaryApisResponse, err error) {
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
func (a *Apiary) GetTeamApis(team string) (apis *ApiaryApisResponse, err error) {
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

// PublishBlueprint publish blueprint in Apiary.io
//
// Reference: http://docs.apiary.apiary.io/#reference/blueprint/publish-blueprint/get-me
func (a *Apiary) PublishBlueprint(name string, content []byte) (published bool, err error) {
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

// FetchBlueprint fetches blueprint from Apiary.io
//
// Reference: Unknown
func (a *Apiary) FetchBlueprint(name string) (blueprint *ApiaryFetchResponse, err error) {
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
