// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import "encoding/json"

func mapToStruct(theMap map[string]interface{}, theStruct interface{}) {
	jsonBytes, _ := json.Marshal(theMap)
	json.Unmarshal(jsonBytes, &theStruct)
}
