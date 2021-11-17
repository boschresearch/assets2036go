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

// TestPropertiesSimpleTypes kjh
func TestProperties(t *testing.T) {
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

	asset, err := mgr.CreateAsset("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Error(err)
	}

	proxy, err := mgr.CreateAssetProxy("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Error(err)
	}

	// set test value
	strValue := "Ein schöner Name"
	err = asset.Submodel("testmodel").Property("string").SetValue(strValue)

	if err != nil {
		t.Fatal(err)
	}

	if !checkPropertiesEquality(
		strValue,
		proxy.Submodel("testmodel").Property("string").Value,
		testTimeout) {
		t.Error("SetValue(string) failed")
	}

	intValue := rand.Int31()
	err = asset.Submodel("testmodel").Property("integer").SetValue(intValue)

	if err != nil {
		t.Fatal(err)
	}

	if !checkPropertiesEquality(
		float64(intValue),
		proxy.Submodel("testmodel").Property("integer").Value,
		testTimeout) {
		t.Error("SetValue(integer) failed")
	}

	floatValue := rand.Float64()
	err = asset.Submodel("testmodel").Property("number").SetValue(floatValue)

	if err != nil {
		t.Fatal(err)
	}

	if !checkPropertiesEquality(
		floatValue,
		proxy.Submodel("testmodel").Property("number").Value,
		testTimeout) {
		t.Error("SetValue(number) failed")
	}

	// test object property
	localperson := PersonWOArray{
		Name:   "Thomas",
		Age:    99,
		Weight: 70,
		Nice:   true,
	}

	err = asset.Submodel("testmodel").Property("person").SetValue(localperson)

	if err != nil {
		t.Fatal(err)
	}

	// time.Sleep(100 * time.Millisecond)

	var receivedPerson PersonWOArray

	for started := time.Now(); time.Since(started) <= testTimeout; {
		if proxy.Submodel("testmodel").Property("person").Value() == nil {
			continue
		}

		mapToStruct(
			proxy.Submodel("testmodel").Property("person").Value().(map[string]interface{}),
			&receivedPerson)

		if localperson == receivedPerson {
			break
		}
	}

	if localperson != receivedPerson {
		t.Error("SetValue(object) failed")
	}

	var receivedValue int32

	// check modified event
	proxy.Submodel("testmodel").Property("integer").AddModifiedListener(func(p *assets2036go.SubmodelProperty, newValue interface{}) {
		receivedValue = int32(newValue.(float64))
	})

	var testValue int32 = rand.Int31()
	asset.Submodel("testmodel").Property(("integer")).SetValue(testValue)

	if !checkPropertiesEquality(testValue, func() interface{} { return receivedValue }, 10*time.Second) {
		t.Error("property modified failed")
	}
}

type PersonWOArray struct {
	Name   string  `json:"name"`
	Age    float64 `json:"age"`
	Weight float64 `json:"weight"`
	Nice   bool    `json:"nice"`
}
