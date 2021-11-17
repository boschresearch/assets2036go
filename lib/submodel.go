// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

// Submodel is the basis for all submodel instances. It holds the submodel's
// #Operations, #Properties and #Events. The struct is used for the deserialization
// of the JSON submodel descriptions from the submodel repository as well
// as for accessing the submodel elements at runtime.
type Submodel struct {
	Operations map[string]*SubmodelOperation `json:"operations"`
	Properties map[string]*SubmodelProperty  `json:"properties"`
	Events     map[string]*SubmodelEvent     `json:"events"`
	Name       string                        `json:"name"`
	Revision   string                        `json:"rev"`

	submodelURL    string
	submodelSchema map[string]interface{}

	assetElement
}

func (sm *Submodel) getElements() map[string]iSubmodelElement {
	result := map[string]iSubmodelElement{}

	for n, v := range sm.Operations {
		result[n] = v
	}

	for n, v := range sm.Properties {
		result[n] = v
	}

	for n, v := range sm.Events {
		result[n] = v
	}

	return result
}

// Operation returns the SubmodelOperation named name of submodel
// or nil, if no operation with given name exists
func (sm *Submodel) Operation(name string) *SubmodelOperation {
	if op, ok := sm.Operations[name]; ok {
		return op
	}

	logg().Printf("Operation %v not found", name)
	return nil
}

// Property returns the SubmodelProperty named name of submodel
// or nil, if no property with given name exists
func (sm *Submodel) Property(name string) *SubmodelProperty {
	if pr, ok := sm.Properties[name]; ok {
		return pr
	}

	logg().Printf("Property %v not found", name)

	return nil
}

// Event returns the SubmodelEvent named name of submodel
// or nil, if no event with given name exists
func (sm *Submodel) Event(name string) *SubmodelEvent {
	if ev, ok := sm.Events[name]; ok {
		return ev
	}

	logg().Printf("Event %v not found", name)

	return nil
}

func (sm *Submodel) injectAssetMgr(assetMgr *AssetMgr) {
	sm.assetMgr = assetMgr

	for _, ele := range sm.getElements() {
		ele.injectAssetMgr(assetMgr)
	}
}

func (sm *Submodel) setTopicsRecursively(namespace, assetname string) {
	for name, el := range sm.getElements() {
		el.setTopic(buildTopic(namespace, assetname, sm.Name, name).String())
	}
}
