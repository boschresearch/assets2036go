// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boschresearch/assets2036go/lib/constants"
	uuid "github.com/satori/go.uuid"
)

// SubmodelOperation encapsulates submodel operations.
// The struct is used for runtime access as well as for json
// deserialization from the submodel description from the submodel
// repository.
// To implement a submodel operation on an asset, set the #Callback
// member to a local function.
// To call a remotely implemented operation, call the #Invoke or #InvokeTimeout.
type SubmodelOperation struct {
	Parameters map[string]interface{}    `json:"parameters"`
	Response   map[string]interface{}    `json:"response"`
	Callback   SubmodelOperationCallback `json:"-"`

	assetElement
	// invokeLock sync.Mutex

	currentRequest  *SubmodelOperationRequest
	currentResponse *submodelOperationResponse
}

// Invoke runs calls the remote operation on the remote asset andn returns the remotely computed return values
// in the form of a map[string]interface{}
// If within default #timeout of 2s no answer is received, method return an error.
func (op *SubmodelOperation) Invoke(parameters map[string]interface{}) (interface{}, error) {
	return op.InvokeTimeout(parameters, 2*time.Second)
}

// Name returns the name of the corresponding submodeloperation
func (op *SubmodelOperation) Name() string {
	return op.Topic.SubmodelElement()
}

// InvokeTimeout runs calls the remote operation on the remote asset and returns the remotely computed return values
// in the form of a map[string]interface{}
// If within #timeout no answer is received, method return an error.
func (op *SubmodelOperation) InvokeTimeout(parameters map[string]interface{}, timeout time.Duration) (interface{}, error) {
	// TODO: Make this thread safe!

	// create request object
	op.currentRequest = op.createRequest(parameters)

	// when this func returns, current request is finished, so nil it
	defer func() {
		op.currentRequest = nil
	}()

	// nil old response object
	op.currentResponse = nil

	// send current request
	op.assetMgr.publish(op.currentRequest.topic().String(), op.currentRequest.getPayload(), false)

	// wait for response
	for started := time.Now(); op.currentResponse == nil; time.Sleep(1 * time.Millisecond) {

		if time.Since(started) >= timeout {
			return nil, fmt.Errorf("Operation %v timed out", op.Topic)
		}

		time.Sleep(1 * time.Millisecond)
	}

	// return response
	return op.currentResponse.Response, nil
}

func (op *SubmodelOperation) createRequest(parameters map[string]interface{}) *SubmodelOperationRequest {
	reqID := uuid.NewV4()

	req := &SubmodelOperationRequest{
		Parameters:        parameters,
		RequestID:         reqID.String(),
		SubmodelOperation: op,
	}

	return req
}

func (op *SubmodelOperation) installHandlersProxy(mgr *AssetMgr) {
	mgr.messageHandlers[buildTopic(op.Topic.String(), constants.TopicElementSubmodelOperationResp).String()] = func(topic Topic, rawPayload []byte) {

		var tempResp submodelOperationResponse
		err := json.Unmarshal(rawPayload, &tempResp)
		if err != nil {
			logg().Print(err)
			return
		}

		if op.currentRequest != nil && tempResp.RequestID == op.currentRequest.RequestID {
			op.currentResponse = &tempResp
		}
	}
}

func (op *SubmodelOperation) installHandlersAsset(mgr *AssetMgr) {
	mgr.messageHandlers[append(op.Topic, constants.TopicElementSubmodelOperationReq).String()] = func(topic Topic, rawPayload []byte) {
		err := json.Unmarshal(rawPayload, &op.currentRequest)

		if err != nil {
			logg().Println(err)
			return
		}

		op.currentRequest.SubmodelOperation = op

		if op.Callback != nil {
			respVal, err := op.Callback(op.currentRequest)

			resp := op.currentRequest.createResponse()
			resp.Response = respVal

			if err != nil {
				logg().Print(err)
				resp.Error = err.Error()
			}

			mgr.publish(resp.topic.String(), resp.getPayload(), false)
		}
	}
}
