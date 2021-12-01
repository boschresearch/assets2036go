// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	assets2036go "github.com/boschresearch/assets2036go/lib"
)

func TestOperationsObjectTypes(t *testing.T) {
	t.Parallel()

	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, "arena2036", "gotest", true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
	}

	if mgr == nil {
		t.Errorf("returned mgr was nil")
	}

	submodelURL := fmt.Sprintf(submodelUrlTemplate, "testmodel.json")

	asset, err := mgr.CreateAsset("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Fatal(err)
	}

	asset.Submodel("testmodel").Operation("setObject").Callback = setObject
	asset.Submodel("testmodel").Operation("getObject").Callback = getObject

	proxy, err := mgr.CreateAssetProxy("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Error(err)
	}

	localperson := PersonWithArray{
		Name:           "Thomas",
		Age:            99,
		Weight:         70,
		Nice:           true,
		ArrayOfInteger: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	proxy.Submodel("testmodel").Operation("setObject").InvokeTimeout(map[string]interface{}{
		"person": localperson,
	}, 1000*time.Second)

	if localperson.Name != globalperson.Name {
		t.Error("setObject failed")
	}

	for i := 0; i < 10; i++ {
		if globalperson.ArrayOfInteger[i] != float64(i) {
			t.Error("setObject failed")
		}
	}
}

type PersonWithArray struct {
	Name           string    `json:"name"`
	Age            float64   `json:"age"`
	Weight         float64   `json:"weight"`
	Nice           bool      `json:"nice"`
	ArrayOfInteger []float64 `json:"arrayOfInteger"`
}

var globalperson PersonWithArray

func setObject(req *assets2036go.SubmodelOperationRequest) (response interface{}, err error) {
	jsonBytes, _ := json.Marshal(req.Parameters["person"])

	_ = json.Unmarshal(jsonBytes, &globalperson)

	return nil, nil
}

func getObject(req *assets2036go.SubmodelOperationRequest) (response interface{}, err error) {
	return globalperson, nil
}
