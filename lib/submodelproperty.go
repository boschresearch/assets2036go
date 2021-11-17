// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"
	"errors"
)

// HandlePropertyModified is the callback type to get informed when the property value changes
type HandlePropertyModified func(sp *SubmodelProperty, newValue interface{})

// SubmodelProperty encapsulates a submodel's property.
// To set the property's value on a locally implemented submodel,
// use #SetValue. To read the property's value of a remotely
// implemented sumodel, use #Value() or #AddModifiedListener to add
// a listener which will be called in the case of a remote property value
// modification.
type SubmodelProperty struct {
	value             interface{}
	Type              string `json:"type"`
	modifiedListeners []HandlePropertyModified
	assetElement
}

// AddModifiedListener lets you add a listener, whichh will be informed,
// when the remotely implemented property has modified value.
func (sp *SubmodelProperty) AddModifiedListener(listener HandlePropertyModified) {
	if sp.modifiedListeners == nil {
		sp.modifiedListeners = []HandlePropertyModified{}
	}

	sp.modifiedListeners = append(sp.modifiedListeners, listener)
}

// Name returns the name of the corresponding submodeloperation
func (sp *SubmodelProperty) Name() string {
	return sp.Topic.SubmodelElement()
}

// Value returns the current value of this Property Element
func (sp *SubmodelProperty) Value() interface{} {
	return sp.value
}

// SetValue sets the value and publishes it to all remote consumers.
func (sp *SubmodelProperty) SetValue(value interface{}) error {
	if sp.asset == nil {
		return errors.New("cannot set property value on an asset proxy - just on an asset")
	}

	sp.value = value

	// TODO: Check if the value actually has changed - only send when changed
	jsonPayloa, err := json.Marshal(value)

	if err != nil {
		logg().Fatal(err)
	}

	sp.assetMgr.publish(sp.Topic.String(), jsonPayloa, true)

	return nil
}

func (sp *SubmodelProperty) installHandlers(mgr *AssetMgr) {
	mgr.messageHandlers[sp.Topic.String()] = func(topic Topic, rawPaylod []byte) {
		var newValue interface{}

		err := json.Unmarshal(rawPaylod, &newValue)

		if err != nil {
			logg().Print(err)
			return
		}

		sp.value = newValue
		if sp.modifiedListeners != nil {
			for _, v := range sp.modifiedListeners {
				v(sp, newValue)
			}
		}
	}
}
