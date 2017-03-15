package apiary

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"errors"
	"strings"
	"io"
)

//
// Data
//
var Token = os.Getenv("APIARY_TOKEN")
var Repository = os.Getenv("APIARY_REPO")
var Team = os.Getenv("APIARY_TEAM")

var ValidBlueprint = []byte(`FORMAT: 1A
HOST: http://api.example.com/

# Example API\n\nIntroduction.
# And update
`)

//
// Helpers
//
type fakeReader struct {
	io.Reader
}

func (r *fakeReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("OMG!")
}

//
// Test suite
//
func Test_ReadResponse(t *testing.T) {
	t.Run("Return empty response Error", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte(``))
		rc := ioutil.NopCloser(buf)

		response := &http.Response{
			ContentLength: 10,
			Body:          rc,
		}
		b, err := readResponse(response)

		if len(b) != 0 {
			t.Error("Something parsed from empty response")
		}

		if err.Error() != errors.New("Empty response").Error() {
			t.Error("Empty response should throw: [Empty response] error")
		}
	})

	t.Run("Bump error on Read from Buffer error", func(t *testing.T) {
		buf := strings.NewReader("Non empty response")
		reader := &fakeReader{buf}
		rc := ioutil.NopCloser(reader)

		response := &http.Response{
			ContentLength: 10,
			Body:          rc,
		}
		_, err := readResponse(response)

		if err.Error() != errors.New("OMG!").Error() {
			t.Error("Should throw: [OMG!] error")
		}
	})
}

func TestApiary_Me(t *testing.T) {
	t.Run("Retrieve data", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		r, err := a.Me()

		if r.ID == "" {
			t.Error("Empty ID")
		}

		if r.Name == "" {
			t.Error("Empty Name")
		}

		if r.URL == "" {
			t.Error("Empty URL")
		}

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
	})

	t.Run("Empty token", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: "",
		})

		_, err := a.Me()

		if err == nil {
			t.Error("Expected error returned on empty token")
		}
	})
}

func TestApiary_GetApis(t *testing.T) {
	t.Run("Retrieve data", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		r, err := a.GetApis()

		if r == nil || len(r.Apis) == 0 {
			t.Error("Empty apis returned")
		}

		for _, api := range r.Apis {
			if api.Name == "" {
				t.Error("Empty api name")
			}

			if api.DocumentationURL == "" {
				t.Error("Empty documentation URL")
			}

			if api.Subdomain == "" {
				t.Error("Empty sudbomain URL")
			}
		}

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
	})

	t.Run("Empty token", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: "",
		})

		_, err := a.GetApis()

		if err == nil {
			t.Error("Expected error returned on empty token")
		}
	})
}

func TestApiary_GetTeamApis(t *testing.T) {
	t.Run("Get invalid team", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		_, err := a.GetTeamApis("some_invalid_team_name")
		if err == nil {
			t.Error("Invalid team name should return error")
		}
	})

	t.Run("Get team", func(t *testing.T) {
		if Team == "" {
			t.Skip("Empty team token")
		}

		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		r, err := a.GetTeamApis(Team)

		if len(r.Apis) == 0 {
			t.Error("Empty team apis")
		}

		for _, api := range r.Apis {
			if api.Name == "" {
				t.Error("Empty api name")
			}

			if api.DocumentationURL == "" {
				t.Error("Empty documentation URL")
			}

			if api.Subdomain == "" {
				t.Error("Empty sudbomain URL")
			}
		}

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
	})
}

func TestApiary_FetchBlueprint(t *testing.T) {
	t.Run("Fetching blueprint", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		r, err := a.FetchBlueprint(Repository)

		if r.Error {
			t.Errorf("Error fetching repo: %s", r.Message)
		}

		if r.Message != "" {
			t.Errorf("Error fetching repo: %s", r.Message)
		}

		if r.Code == "" {
			t.Error("Empty repository code")
		}

		if r.Code != string(ValidBlueprint) {
			t.Error("Different blueprints")
		}

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
	})
}

func TestApiary_PublishBlueprint(t *testing.T) {
	t.Run("Publish blueprint", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		publish, err := a.PublishBlueprint(Repository, ValidBlueprint)

		if !publish {
			t.Error("Not published")
		}

		if err != nil {
			t.Error(fmt.Sprintf("Error: %s", err))
		}
	})

	t.Run("Publish same blueprint", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		publish, err := a.PublishBlueprint(Repository, ValidBlueprint)

		if !publish {
			t.Error("Not published")
		}

		if err != nil {
			t.Error(fmt.Sprintf("Error: %s", err))
		}
	})

	t.Run("Publish with wrong token", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: "",
		})

		publish, err := a.PublishBlueprint(Repository, ValidBlueprint)

		if publish {
			t.Error("Published")
		}

		if err == nil {
			t.Error("Wrong token should generate error")
		}
	})

	t.Run("Publish to repo with no rights", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: Token,
		})

		publish, err := a.PublishBlueprint("testingapiaryclitestingapiarycli", ValidBlueprint)

		if publish {
			t.Error("Published")
		}

		if err == nil {
			t.Error("Wrong repository should generate error")
		}
	})

	t.Run("Publish wrong content", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{
			Token: "",
		})

		publish, err := a.PublishBlueprint(Repository, []byte("some invalid data"))

		if publish {
			t.Error("Published")
		}

		if err == nil {
			t.Error("Wrong token should generate error")
		}
	})
}
