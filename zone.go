package cfdns

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Zone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ZonesResponse struct {
	Response
	Zones []Zone `json:"result"`
}

func (client *Client) ListZones() (*ZonesResponse, error) {
	var params map[string]string
	responseBody, err := client.get("zones", params)
	if err != nil {
		return nil, err
	}
	var zones ZonesResponse
	json.Unmarshal(responseBody, &zones)
	return &zones, nil
}

func (client *Client) GetZoneIdByName(name string) (string, error) {
	params := make(map[string]string)
	params["name"] = name
	responseBody, err := client.get("zones", params)
	if err != nil {
		return "", err
	}
	var zones ZonesResponse
	json.Unmarshal(responseBody, &zones)
	if !zones.Success {
		errorMsg := ""
		for _, jsonError := range zones.Errors {
			errorMsg = fmt.Sprintf("%s\n%d: %s. ", errorMsg, jsonError.Code, jsonError.Message)
		}
		return "", errors.New(errorMsg)
	}
	return zones.Zones[0].Id, nil
}
