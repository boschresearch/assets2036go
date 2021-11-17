// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import (
	"math/rand"
	"testing"
	"time"

	assets2036go "github.com/boschresearch/assets2036go/lib"
)

// TestEvents will test events
func TestEvents(t *testing.T) {
	t.Parallel()

	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, "arena2036", "gotest", true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
	}

	if mgr == nil {
		t.Errorf("returned mgr was nil")
	}

	submodelURL := "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/raw/master/testmodel.json"

	asset, _ := mgr.CreateAsset("arena2036", "gotest", submodelURL)
	proxy, _ := mgr.CreateAssetProxy("arena2036", "gotest", submodelURL)

	eventReceived := false
	proxy.Submodel("testmodel").Event("voidEvent").Handler = func(timestamp time.Time, params map[string]interface{}) {
		eventReceived = true
	}

	asset.Submodel("testmodel").Event("voidEvent").Emit(nil)

	started := time.Now()

	for time.Since(started) <= testTimeout && !eventReceived {
		time.Sleep(time.Millisecond)
	}

	if !eventReceived {
		t.Error("error in voidEvent")
	}

	eventReceived = false
	proxy.Submodel("testmodel").Event("boolEvent").Handler = func(timestamp time.Time, params map[string]interface{}) {
		eventReceived = params["param_1"].(bool)
	}

	asset.Submodel("testmodel").Event("boolEvent").Emit(map[string]interface{}{"param_1": true})

	started = time.Now()

	for time.Since(started) <= testTimeout && !eventReceived {
		time.Sleep(time.Millisecond)
	}

	if !eventReceived {
		t.Error("error in boolEvent")
	}

	eventReceived = false
	floatValue := rand.Float64()
	proxy.Submodel("testmodel").Event("numberEvent").Handler = func(timestamp time.Time, params map[string]interface{}) {
		eventReceived = params["param_1"].(float64) == floatValue
	}

	asset.Submodel("testmodel").Event("numberEvent").Emit(map[string]interface{}{"param_1": floatValue})

	started = time.Now()

	for time.Since(started) <= testTimeout && !eventReceived {
		time.Sleep(time.Millisecond)
	}

	if !eventReceived {
		t.Error("error in numberEvent")
	}

	eventReceived = false
	globalPerson := PersonWOArray{
		Name:   "Thomas",
		Age:    99,
		Weight: 70,
		Nice:   true,
	}

	proxy.Submodel("testmodel").Event("objectEvent").Handler = func(timestamp time.Time, params map[string]interface{}) {
		var localPerson PersonWOArray
		mapToStruct(params["param_1"].(map[string]interface{}), &localPerson)
		eventReceived = (localPerson == globalPerson)
	}

	asset.Submodel("testmodel").Event("objectEvent").Emit(map[string]interface{}{"param_1": globalPerson})

	started = time.Now()

	for time.Since(started) <= testTimeout && !eventReceived {
		time.Sleep(time.Millisecond)
	}

	if !eventReceived {
		t.Error("error in objectEvent")
	}
}
