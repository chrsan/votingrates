package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Only relevant JSON fields included.
type TableMetadataVariable struct {
	Text       string   `json:"text"`
	Values     []string `json:"values"`
	ValueTexts []string `json:"valueTexts"`
}

type TableMetadata struct {
	Variables []TableMetadataVariable `json:"variables"`
}

type TableResponseData struct {
	Key    []string `json:"key"`
	Values []string `json:"values"`
}

type TableResponse struct {
	Data []TableResponseData `json:"data"`
}

type API interface {
	VotingRatesMetadata() (TableMetadata, error)
	VotingRatesQuery() (TableResponse, error)
}

const VotingRatesURL = "http://api.scb.se/OV0104/v1/doris/en/ssd/ME/ME0104/ME0104D/ME0104T4"

type HTTP struct{}

func (h *HTTP) VotingRatesMetadata() (TableMetadata, error) {
	var tmd TableMetadata
	resp, err := http.Get(VotingRatesURL)
	if err != nil {
		return tmd, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tmd, fmt.Errorf("HTTP server responded with status: %d", resp.StatusCode)
	}
	if err := readJSON(resp.Body, &tmd); err != nil {
		return tmd, err
	}
	return tmd, nil
}

const VotingRatesQueryJSON = `
{
  "query": [
    {"code": "Region", "selection": {"filter": "all", "values": ["*"]}},
    {"code": "ContentsCode", "selection": {"filter": "item", "values": ["ME0104B8"]}}
  ],
  "response": {"format": "json"}
}
`

func (h *HTTP) VotingRatesQuery() (TableResponse, error) {
	var tr TableResponse
	req, err := http.NewRequest(http.MethodPost, VotingRatesURL, strings.NewReader(VotingRatesQueryJSON))
	if err != nil {
		return tr, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tr, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tr, fmt.Errorf("HTTP server responded with status: %d", resp.StatusCode)
	}
	if err := readJSON(resp.Body, &tr); err != nil {
		return tr, err
	}
	return tr, nil
}

func readJSON(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if len(data) >= 3 {
		// Remove possible Windows UTF-8 BOM
		if data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
			data = data[3:]
		}
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
