package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
	Token      string
}

type CovidStatus struct {
	Odp       int `json:"odp"`
	Pdp       int `json:"pdp"`
	Positif   int `json:"positif"`
	Sembuh    int `json:"sembuh"`
	Meninggal int `json:"meninggal"`
}

type CheckRequest struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type CheckResponse struct {
	Kelurahan   string         `json:"kelurahan"`
	Kecamatan   string         `json:"kecamatan"`
	Kabkot      string         `json:"kabkot"`
	Provinsi    string         `json:"provinsi"`
	CovidStatus []*CovidStatus `json:"covidStatus"`
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", c.Token)
	return req, nil
}
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// Check covid status by lon, lat
func (c *Client) Check(cr *CheckRequest) (*CheckResponse, error) {
	req, err := c.newRequest("POST", "/api/covid/kelurahan/data", cr)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Status string         `json:"status"`
		Data   *CheckResponse `json:"data"`
	}
	_, err = c.do(req, &resp)
	return resp.Data, err
}

// NewClient return new api client
func NewClient(baseUrl, token string) (*Client, error) {

	if baseUrl == "" {
		return nil, errors.New("baseURL shouldn't be empty")
	}

	if token == "" {
		return nil, errors.New("token shouldn't be empty")
	}

	u, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	c := &Client{
		BaseURL:    u,
		Token:      token,
		httpClient: http.DefaultClient,
	}

	return c, err
}
