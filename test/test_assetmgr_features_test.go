// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import (
	"log"
	"testing"

	assets2036go "github.com/boschresearch/assets2036go/lib"
)

func TestFindAssetsBySubmodelName(t *testing.T) {
	// t.Parallel()

	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, "arena2036", "gotest", true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
	}

	submodelURL := "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/-/raw/master/testmodel.json"

	testasset, err := mgr.CreateAsset("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Error(err)
	}

	testasset.Submodel("testmodel").Property("string").SetValue("Something")

	res, err := mgr.FindAssetsBySubmodelName("testmodel")

	if err != nil {
		t.Error(err)
	}

	log.Printf("Found assets with given submodel name: \n%v", assets2036go.StringList(res))

	if len(res) <= 0 {
		t.Error("No asset with submodel testmodel")
	}
}

func TestFindAssetsBySubmodelURL(t *testing.T) {
	t.Parallel()

	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, "arena2036", "gotest", true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
	}

	submodelURL := "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/-/raw/master/testmodel.json"

	testasset, err := mgr.CreateAsset("arena2036", "gotest", submodelURL)
	if err != nil {
		t.Error(err)
	}

	testasset.Submodel("testmodel").Property("string").SetValue("Something")

	res, err := mgr.FindAssetsBySubmodelURL(submodelURL)

	if err != nil {
		t.Error(err)
	}

	log.Printf("Found assets with given submodel url: \n%v", assets2036go.StringList(res))

	if len(res) <= 0 {
		t.Errorf("No asset found with submodelURL %v", submodelURL)
	}
}

func TestFindEndpointAssets(t *testing.T) {
	// t.Parallel()

	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, "arena2036", "gotest", true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
	}

	submodelURL := "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/-/raw/master/testmodel.json"

	testasset, err := mgr.CreateAsset("arena2036", "gotest", submodelURL)

	if err != nil {
		t.Error(err)
	}

	testasset.Submodel("testmodel").Property("string").SetValue("Something")

	res, err := mgr.FindAllEndpointAssets()

	if err != nil {
		t.Error(err)
	}

	log.Printf("Found endpoint assets: \n%v", assets2036go.StringList(res))

	if len(res) <= 0 {
		t.Errorf("No endpoint asset found with submodelURL %v", submodelURL)
	}
}
