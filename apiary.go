package apiary

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	ApiaryAPIURL = "https://api.apiary.io/"
)

const (
	ApiaryActionMe               = "me"
	ApiaryActionGetApis          = "me/apis"
	ApiaryActionGetTeamApis      = "me/teams/%s/apis"
	ApiaryActionFetchBlueprint   = "blueprint/get/%s"
	ApiaryActionPublishBlueprint = "blueprint/publish/%s"
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

type Apiary struct {
	options ApiaryOptions
	client  *http.Client
}

type ApiaryOptions struct {
	Token string
}

func NewApiary(opts ApiaryOptions) *Apiary {
	return &Apiary{
		options: opts,
		client:  &http.Client{},
	}
}

func checkOk(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Bad response code: %s", response.Status))
	}

	return nil
}

func readResponse(response *http.Response) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, response.ContentLength))
	n, err := buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, errors.New("Empty response")
	}

	return buf.Bytes(), nil
}

func bearerToken(token string) string {
	buf := bytes.NewBuffer(make([]byte, 0, len(token)+7))
	buf.Write([]byte(`bearer `))
	buf.Write([]byte(token))

	return buf.String()
}
func bearerTokenLegacy(token string) string {
	buf := bytes.NewBuffer(make([]byte, 0, len(token)+6))
	buf.Write([]byte(`Token `))
	buf.Write([]byte(token))

	return buf.String()
}

func (a *Apiary) request(method string, path string, headers map[string]string, body io.Reader) (response []byte, res *http.Response, err error) {
	url := ApiaryAPIURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err = a.client.Do(req)
	if err != nil {
		return
	}

	response, err = readResponse(res)
	return
}

func (a *Apiary) sendRequest(path string) (data []byte, response *http.Response, err error) {
	headers := make(map[string]string)
	headers["Authorization"] = bearerToken(a.options.Token)
	data, response, err = a.request("GET", path, headers, nil)
	return
}

func (a *Apiary) sendLegacyRequest(path string) (data []byte, response *http.Response, err error) {
	headers := make(map[string]string)
	headers["Authentication"] = bearerTokenLegacy(a.options.Token)
	data, response, err = a.request("GET", path, headers, nil)
	return
}

func (a *Apiary) sendLegacyPostRequest(path string, body io.Reader) (data []byte, response *http.Response, err error) {
	headers := make(map[string]string)
	headers["Authentication"] = bearerTokenLegacy(a.options.Token)
	headers["Content-Type"] = "application/json; charset=utf-8"
	data, response, err = a.request("POST", path, headers, body)
	return
}


// Me retrieve user information
// Reference: http://docs.apiary.apiary.io/#reference/user-information/me/get-me
func (a *Apiary) Me() (me ApiaryMeResponse, err error) {
	data, response, err := a.sendRequest(ApiaryActionMe)
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
// Reference: http://docs.apiary.apiary.io/#reference/api-list/user-api-list/get-me
func (a *Apiary) GetApis() (apis *ApiaryApisResponse, err error) {
	data, response, err := a.sendRequest(ApiaryActionGetApis)
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
// Reference: http://docs.apiary.apiary.io/#reference/api-list/team-api-list/get-me
func (a *Apiary) GetTeamApis(team string) (apis *ApiaryApisResponse, err error) {
	uri := fmt.Sprintf(ApiaryActionGetTeamApis, team)
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
// Reference: http://docs.apiary.apiary.io/#reference/blueprint/publish-blueprint/get-me
func (a *Apiary) PublishBlueprint(name string, content []byte) (published bool, err error) {
	jsonData, err := json.Marshal(map[string]string{
		"code": string(content),
	})

	if err != nil {
		return
	}

	uri := fmt.Sprintf(ApiaryActionPublishBlueprint, name)
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
// Reference: Unknown
func (a *Apiary) FetchBlueprint(name string) (blueprint *ApiaryFetchResponse, err error) {
	uri := fmt.Sprintf(ApiaryActionFetchBlueprint, name)
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
