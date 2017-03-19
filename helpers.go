package apiary

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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
