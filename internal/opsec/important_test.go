package opsec

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"go.albinodrought.com/neptunes-pride/internal/types"
)

func TestFindThreats(t *testing.T) {
	loadFile := func(path string) *types.APIResponse {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		response := &types.APIResponse{}
		err = json.Unmarshal(data, response)
		if err != nil {
			panic(err)
		}

		return response
	}

	twoThreats := loadFile("two-threats.json")

	threats := FindThreats(twoThreats)
	if len(threats) != 6 {
		t.Fatalf("expected 6 threats but got %v", len(threats))
	}
}
