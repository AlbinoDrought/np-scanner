package opsec

import (
	"testing"
)

func TestFindThreats(t *testing.T) {
	t.Skip("needs more sample logs from NP4")

	// loadFile := func(path string) *types.APIResponse {
	// 	data, err := ioutil.ReadFile(path)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	response := &types.APIResponse{}
	// 	err = json.Unmarshal(data, response)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	return response
	// }

	// twoThreats := loadFile("two-threats.json")

	// threats := FindThreats(twoThreats)
	// if len(threats) != 6 {
	// 	t.Fatalf("expected 6 threats but got %v", len(threats))
	// }
}
