package opsec

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"go.albinodrought.com/neptunes-pride/internal/types"
)

func TestMerge(t *testing.T) {
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

	a := loadFile("aburrido.json")
	b := loadFile("burrito.json")

	merged := Merge(a, b)

	// from aburrido.json
	if merged.ScanningData.Fleets["36"].Strength != 1 {
		t.Errorf("Missing fleet #36 strength data, got %+v", merged.ScanningData.Fleets["36"])
	}
	if merged.ScanningData.Stars["20"].Strength != 4 {
		t.Errorf("Missing star #20 strength data, got %+v", merged.ScanningData.Stars["20"])
	}
	if merged.ScanningData.Stars["22"].Strength != 3 {
		t.Errorf("Missing star #22 strength data, got %+v", merged.ScanningData.Stars["22"])
	}
	if merged.ScanningData.Players["4"].Researching != "manufacturing" {
		t.Errorf("Missing player #4 research data, got %+v", merged.ScanningData.Players["4"])
	}

	// from burrito.json
	if merged.ScanningData.Fleets["2"].Strength != 1 {
		t.Errorf("Missing fleet #2 strength data, got %+v", merged.ScanningData.Fleets["2"])
	}
	if merged.ScanningData.Stars["5"].Strength != 21 {
		t.Errorf("Missing star #5 strength data, got %+v", merged.ScanningData.Stars["5"])
	}
	if merged.ScanningData.Stars["88"].Strength != 9 {
		t.Errorf("Missing star #88 strength data, got %+v", merged.ScanningData.Stars["88"])
	}
	if merged.ScanningData.Players["5"].Researching != "weapons" {
		t.Errorf("Missing player #5 research data, got %+v", merged.ScanningData.Players["5"])
	}

	data, err := json.Marshal(merged)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("merged.json", data, os.ModePerm)
}
