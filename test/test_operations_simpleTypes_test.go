// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import (
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	assets2036go "github.com/boschresearch/assets2036go/lib"
)

func TestOperationsSimpleTypes(t *testing.T) {
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
		log.Fatal(err)
	}

	asset.Submodel("testmodel").Operation("getInteger").Callback = getInteger
	asset.Submodel("testmodel").Operation("setInteger").Callback = setInteger

	asset.Submodel("testmodel").Operation("getNumber").Callback = getNumber
	asset.Submodel("testmodel").Operation("setNumber").Callback = setNumber

	asset.Submodel("testmodel").Operation("getBool").Callback = getBool
	asset.Submodel("testmodel").Operation("setBool").Callback = setBool

	asset.Submodel("testmodel").Operation("getString").Callback = getString
	asset.Submodel("testmodel").Operation("setString").Callback = setString

	// now create proxy to work with the asset
	proxy, err := mgr.CreateAssetProxy("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Error(err)
	}

	_, err = proxy.Submodel("testmodel").Operation("setString").InvokeTimeout(map[string]interface{}{
		"param_1": "333",
	}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	if testString != "333" {
		t.Error("setString without effect")
	}

	_, err = proxy.Submodel("testmodel").Operation("setInteger").InvokeTimeout(map[string]interface{}{
		"param_1": 333,
	}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	if testInteger != 333 {
		t.Error("setInteger without effect")
	}

	result, err := proxy.Submodel("testmodel").Operation("getInteger").InvokeTimeout(map[string]interface{}{}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	if int(result.(float64)) != 333 {
		t.Error("getInteger failed")
	}

	_, err = proxy.Submodel("testmodel").Operation("setNumber").InvokeTimeout(map[string]interface{}{
		"param_1": 333.333,
	}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	if testNumber != 333.333 {
		t.Error("setNumber without effect")
	}

	result, err = proxy.Submodel("testmodel").Operation("getNumber").InvokeTimeout(map[string]interface{}{}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	if result.(float64) != 333.333 {
		t.Error("getNumber failed")
	}
}

var testInteger = 99

func setInteger(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	testInteger = int(req.Parameters["param_1"].(float64))
	return nil, nil
}

func getInteger(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	return testInteger, nil
}

var testBool = false

func setBool(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	testBool = req.Parameters["param_1"].(bool)
	return nil, nil
}

func getBool(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	return testBool, nil
}

var testString string = "Teststring aus assets2036go"

func setString(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	testString = req.Parameters["param_1"].(string)
	return nil, nil
}

func getString(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	return testString, nil
}

var testNumber float64 = math.Pi

func setNumber(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	testNumber = req.Parameters["param_1"].(float64)
	return nil, nil
}

func getNumber(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	return testNumber, nil
}
