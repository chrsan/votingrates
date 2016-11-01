package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCollectingRates(t *testing.T) {
	var api fileAPI
	want := []Rate{
		{[]string{"Vellinge"}, "1973", 95.4},
		{[]string{"Lomma"}, "1976", 96},
		{[]string{"Lomma", "Vellinge"}, "1979", 95.3},
		{[]string{"Lomma"}, "1982", 95.7},
		{[]string{"Danderyd"}, "1985", 94.5},
		{[]string{"Danderyd", "Vellinge"}, "1988", 92},
		{[]string{"Danderyd"}, "1991", 93},
		{[]string{"Vellinge"}, "1994", 92.6},
		{[]string{"Danderyd"}, "1998", 88.2},
		{[]string{"Lomma"}, "2002", 88.1},
		{[]string{"Lomma"}, "2006", 89.7},
		{[]string{"Lomma"}, "2010", 91.3},
		{[]string{"Lomma"}, "2014", 92.9},
	}
	m, err := Regions(&api)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Rates(&api, m)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("Want %d rates, got %d", len(want), len(got))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want %v, got %v", want, got)
	}
}

type fileAPI struct{}

func (f *fileAPI) VotingRatesMetadata() (TableMetadata, error) {
	var tmd TableMetadata
	if err := readJSONFile("testdata/metadata.json", &tmd); err != nil {
		return tmd, err
	}
	return tmd, nil
}

func (f *fileAPI) VotingRatesQuery() (TableResponse, error) {
	var tr TableResponse
	if err := readJSONFile("testdata/query.json", &tr); err != nil {
		return tr, err
	}
	return tr, nil
}

func readJSONFile(fn string, v interface{}) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	if err := readJSON(bytes.NewReader(data), v); err != nil {
		return err
	}
	return nil
}
