// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"
	"time"

	"github.com/boschresearch/assets2036go/lib/constants"
)

// SubmodelEvent encapsulates the event type submodel element.
// Struct is used for json deserialization and for runtime access.
type SubmodelEvent struct {
	Parameters map[string]interface{} `json:"parameters"`
	Handler    EventHandler           `json:"-"`

	assetElement
}

// Name returns name of the corresponding event
func (ev *SubmodelEvent) Name() string {
	return ev.Topic.SubmodelElement()
}

// Emit emits this event and thereby can be received by all
// subscribes consumers
func (ev *SubmodelEvent) Emit(params map[string]interface{}) {
	emission := map[string]interface{}{
		constants.PayloadPropEvTimestamp: time.Now().Format(time.RFC3339),
		constants.PayloadPropOpReqParams: params,
	}

	jsonBytes, _ := json.Marshal(emission)

	ev.assetMgr.publish(ev.Topic.String(), jsonBytes, false)
}

func (ev *SubmodelEvent) installHandlers(mgr *AssetMgr) {
	mgr.messageHandlers[ev.Topic.String()] = func(topic Topic, rawPayload []byte) {
		if ev.Handler != nil {
			var temp map[string]interface{} = make(map[string]interface{})
			json.Unmarshal(rawPayload, &temp)

			timestamp, _ := time.Parse(time.RFC3339, temp[constants.PayloadPropEvTimestamp].(string))

			var params map[string]interface{} = nil
			if temp[constants.PayloadPropOpReqParams] != nil {
				params = temp[constants.PayloadPropOpReqParams].(map[string]interface{})
			}

			ev.Handler(
				timestamp,
				params)
		}
	}
}

// EventHandler is the callback type for an event
type EventHandler func(timestamp time.Time, params map[string]interface{})
