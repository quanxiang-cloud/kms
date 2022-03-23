package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kms/internal/enums"
	"net/http"
	"net/url"
	"strings"
)

const (
	fieldSplit    = "."
	defaultMethod = "GET"
)

// D data of http
type D struct {
	Method string
	URL    string

	Header map[string]string
	Body   map[string]interface{}
	Query  map[string]string

	Response *http.Response
	jsonBody map[string]interface{}
}

// R replace entity
type R struct {
	Old string
	New string
}

// Do request
func (d *D) Do(client *http.Client, r []*R) (*http.Response, error) {
	request, err := d.BuildRequest(r)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	// TODO: catch response code != 200
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failure, code: %d", response.StatusCode)
	}

	d.Response = response
	return d.Response, nil
}

// BuildRequest build
func (d *D) BuildRequest(r []*R) (*http.Request, error) {
	request := &http.Request{
		Header: http.Header{},
	}

	if err := d.buildBody(r, request); err != nil {
		return nil, err
	}

	if err := d.buildURL(r, request); err != nil {
		return nil, err
	}

	if err := d.buildHeader(r, request); err != nil {
		return nil, err
	}

	request.Method = defaultMethod
	if enums.HTTPMethodSet.Verify(d.Method) {
		request.Method = d.Method
	}

	return request, nil
}

func (d *D) buildBody(r []*R, req *http.Request) error {
	if d.Body != nil {
		b, err := json.Marshal(d.Body)
		if err != nil {
			return err
		}
		if len(r) > 0 {
			body := string(b)
			body = replace(body, r)
			b = []byte(body)
		}
		req.Body = io.NopCloser(bytes.NewReader(b))
	}
	return nil
}

func (d *D) buildURL(r []*R, req *http.Request) error {
	urlPath := d.URL
	if d.Query != nil {
		urlPath += "?"
		for k, v := range d.Query {
			urlPath += k + "=" + v + "&"
		}
		if len(r) > 0 {
			urlPath = replace(urlPath, r)
		}
	}
	authURL, err := url.ParseRequestURI(urlPath)
	if err != nil {
		return err
	}

	req.URL = authURL
	return nil
}

func (d *D) buildHeader(r []*R, req *http.Request) error {
	if d.Header != nil && len(r) > 0 {
		b, err := json.Marshal(d.Header)
		if err != nil {
			return err
		}
		sh := replace(string(b), r)

		var header map[string]string
		err = json.Unmarshal([]byte(sh), &header)
		if err != nil {
			return err
		}

		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	return nil
}

// GetJSONBody get field from response body
func (d *D) GetJSONBody(field string) (interface{}, error) {
	if d.jsonBody == nil {
		b, err := io.ReadAll(d.Response.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &d.jsonBody)
		if err != nil {
			return nil, err
		}
	}

	body := d.jsonBody
	f := strings.Split(field, fieldSplit)
	for i := 0; i < len(f)-1; i++ {
		if c, ok := body[f[i]].(map[string]interface{}); ok {
			body = c
		} else {
			return nil, fmt.Errorf("not found field in response")
		}
	}

	if ret, ok := body[f[len(f)-1]]; ok {
		return ret, nil
	}
	return nil, fmt.Errorf("not found field in response")
}

// GetHeader get header from response header
func (d *D) GetHeader(field string) string {
	return d.Response.Header.Get(field)
}

// GetCookies get cookie
func (d *D) GetCookies() string {
	builder := strings.Builder{}
	for _, v := range d.Response.Cookies() {
		s := fmt.Sprintf("%s=%s", v.Name, v.Value)
		if builder.String() != "" {
			builder.WriteString(fmt.Sprintf("; %s", s))
		} else {
			builder.WriteString(s)
		}
	}
	return builder.String()
}

func replace(src string, r []*R) string {
	for _, v := range r {
		src = strings.ReplaceAll(src, v.Old, v.New)
	}
	return src
}
