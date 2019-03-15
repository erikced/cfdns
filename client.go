package cfdns

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const apiBaseAddress = "https://api.cloudflare.com/client/v4/"

type Response struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"messsages"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		TotalPages int `json:"total_pages"`
		Count      int `json:"count"`
		TotalCount int `json:"count"`
	} `json:"result_info"`
}

func (r *Response) Ok() bool {
	return r.Success
}

func (r *Response) FormatErrors() string {
	builder := strings.Builder{}
	for _, error := range r.Errors {
		builder.WriteString(fmt.Sprintf("%d: %s. ", error.Code, error.Message))
	}
	return builder.String()
}

type Client struct {
	apiKey     string
	email      string
	httpClient http.Client
}

func NewClient(email string, apiKey string) Client {
	return Client{email: email, apiKey: apiKey, httpClient: http.Client{}}
}

func (client *Client) get(path string, parameters map[string]string) ([]byte, error) {
	return client.request("GET", path, parameters, nil)
}

func (client *Client) post(path string, parameters map[string]string, data []byte) ([]byte, error) {
	return client.request("POST", path, parameters, data)
}

func (client *Client) put(path string, parameters map[string]string, data []byte) ([]byte, error) {
	return client.request("PUT", path, parameters, data)
}

func (client *Client) delete(path string) error {
	_, err := client.request("DELETE", path, make(map[string]string), nil)
	return err
}

func (client *Client) request(requestType, path string, parameters map[string]string, data []byte) ([]byte, error) {
	i := 0
	address := fmt.Sprint(apiBaseAddress, path)
	for key, value := range parameters {
		prefix := "?"
		if i > 0 {
			prefix = "&"
		}
		address = fmt.Sprint(address, prefix, key, "=", value)
		i++
	}
	req, err := http.NewRequest(requestType, address, bytes.NewBuffer(data))
	req.Header.Add("X-Auth-Key", client.apiKey)
	req.Header.Add("X-Auth-Email", client.email)
	req.Header.Add("Content-Type", "application/json")
	httpResponse, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	respBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
