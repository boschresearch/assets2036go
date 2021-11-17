// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"

	"github.com/boschresearch/assets2036go/lib/constants"
)

// SubmodelOperationRequest encapsulates the call to a remote submodel operation.
// It is used for deserialization of the transferred message and in the callback
// methods implementing the submodel operation.
// It mainly carries the request #Parameters.
type SubmodelOperationRequest struct {
	RequestID         string                 `json:"req_id"`
	Parameters        map[string]interface{} `json:"params"`
	SubmodelOperation *SubmodelOperation
}

// Topic return the topic under which this request will be sent.
func (sor *SubmodelOperationRequest) topic() Topic {
	return append(sor.SubmodelOperation.Topic, constants.TopicElementSubmodelOperationReq)
}

func (sor *SubmodelOperationRequest) getPayload() []byte {

	payload := map[string]interface{}{
		constants.PayloadPropReqID:       sor.RequestID,
		constants.PayloadPropOpReqParams: sor.Parameters,
	}

	result, _ := json.Marshal(payload)

	return result
}

func (sor *SubmodelOperationRequest) createResponse() submodelOperationResponse {
	resp := submodelOperationResponse{
		RequestID: sor.RequestID,
	}

	resp.topic = append(sor.SubmodelOperation.Topic, constants.TopicElementSubmodelOperationResp)

	return resp
}
