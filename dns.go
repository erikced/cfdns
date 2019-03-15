package cfdns

import (
	"encoding/json"
	"fmt"
	"time"
)

type DnsRecord struct {
	Type       string    `json:"type,omitempty"`
	Name       string    `json:"name,omitempty"`
	Content    string    `json:"content,omitempty"`
	TTL        int       `json:"ttl,omitempty"`
	Proxied    bool      `json:"proxied,omitempty"`
	Id         string    `json:"id,omitempty"`
	Proxiable  bool      `json:"proxiable,omitempty"`
	Locked     bool      `json:"locked,omitempty"`
	ZoneId     string    `json:"zone_id,omitempty"`
	ZoneName   string    `json:"zone_name,omitempty"`
	CreatedOn  time.Time `json:"created_on,omitempty"`
	ModifiedOn time.Time `json:"modified_on,omitempty"`
}

type DnsRecordResponse struct {
	Response
	Record DnsRecord `json:"result"`
}

type DnsRecordsResponse struct {
	Response
	Records []DnsRecord `json:"result"`
}

type DnsRecordFilter struct {
	// Type specifies the type of record to be matched using the filter,
	// either of: A, AAAA, CNAME, TXT, SRV, LOC, MX, NS or SPF.
	Type string
	// Name contains the name of the record e.g. test.example.com.
	Name string
	// Content selects record with a specific content, e.g. "127.0.0.1"
	Content string
	// Match determines whether "all" or "any" of the properties above
	// should match. Defaults to "all" if empty.
	Match string
}

type updatedDnsRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     *int   `json:"ttl,omitempty"`
	Proxied *bool  `json:"proxied,omitempty"`
}

type newDnsRecord struct {
	updatedDnsRecord
	Priority *int `json:"priority,omitempty"`
}

func (client *Client) ListDnsRecords(zoneId string, filter DnsRecordFilter) (DnsRecordsResponse, error) {
	path := fmt.Sprintf("zones/%s/dns_records", zoneId)
	parameters := make(map[string]string)
	if len(filter.Type) > 0 {
		parameters["type"] = filter.Type
	}
	if len(filter.Name) > 0 {
		parameters["name"] = filter.Name
	}
	if len(filter.Content) > 0 {
		parameters["content"] = filter.Content
	}
	if len(filter.Match) > 0 {
		parameters["match"] = filter.Match
	}
	responseBody, err := client.get(path, parameters)
	if err != nil {
		return DnsRecordsResponse{}, err
	}
	var records DnsRecordsResponse
	json.Unmarshal(responseBody, &records)
	return records, nil
}

func (client *Client) CreateDnsRecord(zoneId, recordType, name, content string, ttl, priority *int, proxied *bool) (DnsRecordResponse, error) {
	record := newDnsRecord{updatedDnsRecord{recordType, name, content, ttl, proxied}, priority}
	data, err := json.Marshal(record)
	if err != nil {
		return DnsRecordResponse{}, err
	}
	path := fmt.Sprintf("zones/%s/dns_records", zoneId)
	responseBody, err := client.post(path, nil, data)
	if err != nil {
		return DnsRecordResponse{}, err
	}
	var newRecord DnsRecordResponse
	json.Unmarshal(responseBody, &newRecord)
	return newRecord, nil
}

func (client *Client) UpdateDnsRecord(zoneId, recordId, recordType, name, content string, ttl *int, proxied *bool) (DnsRecordResponse, error) {
	record := updatedDnsRecord{Type: recordType, Name: name, Content: content, TTL: ttl, Proxied: proxied}
	data, err := json.Marshal(record)
	if err != nil {
		return DnsRecordResponse{}, err
	}
	path := fmt.Sprintf("zones/%s/dns_records/%s", zoneId, recordId)
	responseBody, err := client.put(path, nil, data)
	if err != nil {
		return DnsRecordResponse{}, err
	}
	var updatedRecord DnsRecordResponse
	json.Unmarshal(responseBody, &updatedRecord)
	return updatedRecord, nil
}

func (client *Client) DeleteDnsRecord(zoneId string, recordId string) error {
	path := fmt.Sprintf("zones/%s/dns_records/%s", zoneId, recordId)
	return client.delete(path)
}
