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

func TestConnectionObservation(t *testing.T) {
	mgr, err := assets2036go.CreateAssetMgr(testHostname, testPort, "arena2036", "gotest", true)
	defer mgr.CloseConnections()

	if err != nil {
		t.Error(err)
		return
	}

	if mgr == nil {
		t.Errorf("returned mgr was nil")
	} else {
		mgr.ExitWhenConnObsvFails = true
	}

	submodelURL := "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/raw/master/testmodel.json"

	asset, err := mgr.CreateAsset("arena2036", "gotest", submodelURL)

	if err != nil {
		log.Fatal(err)
	}

	proxy, err := mgr.CreateAssetProxy("arena2036", "gotest", submodelURL)

	if err != nil {
		log.Fatal(err)
	}

	proxy.Submodel("testmodel").Event("objectEvent").Handler = func(timestamp time.Time, params map[string]interface{}) {
		log.Printf("Received event")
		time.Sleep(5 * time.Second)
	}

	// start dummy event loop
	go func() {
		for {
			asset.Submodel("testmodel").Event("objectEvent").Emit(map[string]interface{}{
				"name":   "Juergen",
				"age":    42,
				"weight": 72,
				"nice":   true,
			})

			time.Sleep(1000 * time.Millisecond)
		}
	}()

	<-time.After(20 * time.Second)
}
