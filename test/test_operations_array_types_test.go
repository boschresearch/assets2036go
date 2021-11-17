// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import (
	"log"
	"testing"
	"time"

	assets2036go "github.com/boschresearch/assets2036go/lib"
)

func TestOperationsArrayTypes(t *testing.T) {
	t.Parallel()

	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, testNamespace, testAsset1Name, true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
	}

	mgr.SetHealthyCallback(func() bool {
		return true
	})

	if mgr == nil {
		t.Errorf("returned mgr was nil")
	}

	submodelURL := "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/raw/master/testmodel.json"

	asset, err := mgr.CreateAsset(testNamespace, testAsset1Name, submodelURL)

	if err != nil {
		log.Fatal(err)
	}

	asset.Submodel("testmodel").Operation("setArrayOfNumber").Callback = setArrayOfNumber
	asset.Submodel("testmodel").Operation("getArrayOfNumber").Callback = getArrayOfNumber

	// now create proxy to work with the asset
	proxy, err := mgr.CreateAssetProxy(testNamespace, testAsset1Name, submodelURL)

	if err != nil {
		t.Error(err)
	}

	arrayOfNumber := [][]float64{{0, 1, 2, 3}, {4, 5, 6, 7}}

	_, err = proxy.Submodel("testmodel").Operation("setArrayOfNumber").InvokeTimeout(map[string]interface{}{"param_1": arrayOfNumber}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 4; j++ {
			if arrayOfNumber[i][j] != float64(i*4+j) {
				t.Error("setArrayOfNumber failed")
			}
		}
	}

	resp, err := proxy.Submodel("testmodel").Operation("getArrayOfNumber").InvokeTimeout(map[string]interface{}{"param_1": arrayOfNumber}, 1000*time.Second)

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 4; j++ {
			if resp.([]interface{})[i].([]interface{})[j].(float64) != arrayOfNumber[i][j] {
				t.Error("setArrayOfNumber failed")
			}
		}
	}
}

var arrayOfNumber [][]float64

func setArrayOfNumber(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	arrayOfNumber = [][]float64{{0, 0, 0, 0}, {0, 0, 0, 0}}

	for i := 0; i < 2; i++ {
		iRow := req.Parameters["param_1"].([]interface{})[i]

		for j := 0; j < 4; j++ {
			arrayOfNumber[i][j] = iRow.([]interface{})[j].(float64)
		}
	}

	return nil, nil
}

func getArrayOfNumber(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {
	return arrayOfNumber, nil
}
