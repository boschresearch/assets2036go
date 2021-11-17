// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

// Asset is a base struct for AssetProxy and Asset
type Asset struct {
	Submodels map[string]*Submodel
	Namespace string
	Name      string
	assetElement
}

// Submodel returns the Submodel with the given name if existing, otherwise nil
func (asset *Asset) Submodel(name string) *Submodel {

	if submodel, ok := asset.Submodels[name]; ok {
		return submodel
	}

	logg().Printf("Submodel %v not found!", name)

	return nil
}
