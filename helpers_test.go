package apiary

import (
	"bytes"
	"errors"
	"gopkg.in/jarcoal/httpmock.v1"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

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
// Testing non-exported functions
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

func Test_Request(t *testing.T) {
	t.Run("Return error on .NewRequest error", func(t *testing.T) {
		a := NewApiary(ApiaryOptions{})
		_, _, err := a.request(";;;", "", map[string]string{}, nil)

		if err == nil {
			t.Error("Bad method should return error")
		}
	})

	t.Run("Return error on client request error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		responder := httpmock.NewErrorResponder(errors.New("Error"))
		httpmock.RegisterResponder("GET", ApiaryAPIURL, responder)

		a := NewApiary(ApiaryOptions{})
		_, _, err := a.request("GET", "", map[string]string{}, nil)

		if err == nil {
			t.Error("Bad client.Do should return error")
		}
	})
}
