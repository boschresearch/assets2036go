// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036gotest

import "time"

const (
	// testHostname = "192.168.100.3"
	testHostname   = "broker.hivemq.com"
	testPort       = 1883
	testTimeout    = 2 * time.Second
	testNamespace  = "testNamespace"
	testAsset1Name = "testAsset1"
	testAsset2Name = "testAsset2"
)

// func checkPropertiesEquality(sollWert interface{}, istWert func() interface{}, timeout time.Duration) bool {
func checkPropertiesEquality(sollWert interface{}, istWertGetter func() interface{}, timeout time.Duration) bool {
	for started := time.Now(); time.Since(started) <= timeout; {
		if sollWert == istWertGetter() {
			return true
		}

		time.Sleep(1 * time.Millisecond)
	}

	return false
}
