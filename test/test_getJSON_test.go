// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import (
	"log"
	"os"
	"testing"

	assets2036go "github.com/boschresearch/assets2036go/lib"
	"github.com/boschresearch/assets2036go/lib/constants"
)

func TestGetJSON(t *testing.T) {
	str, err := assets2036go.GetJSON("file:///./polling_service.json")

	if err != nil {
		t.Error(err)
	}

	log.Printf("getJSON returned %v", str)
}

func TestGetJSONLocalOverwrite(t *testing.T) {
	os.Setenv(constants.EnvVarAssets2036SubmodelsOverwrite, "file:///./")
	defer os.Unsetenv(constants.EnvVarAssets2036SubmodelsOverwrite)

	str, err := assets2036go.GetJSON("http://phantasy.com/part/polling_service.json")

	if err != nil {
		t.Error(err)
	}

	log.Printf("getJSON returned %v", str)
}
