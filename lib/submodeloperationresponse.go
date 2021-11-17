// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"

	"github.com/boschresearch/assets2036go/lib/constants"
)

type submodelOperationResponse struct {
	RequestID string      `json:"req_id"`
	Response  interface{} `json:"resp"`
	Error     string      `json:"error"`
	topic     Topic
}

func (sor *submodelOperationResponse) getPayload() []byte {

	payload := map[string]interface{}{
		constants.PayloadPropReqID:       sor.RequestID,
		constants.PayloadPropOpResp:      sor.Response,
		constants.PayloadPropOpRespError: sor.Error,
	}

	result, _ := json.Marshal(payload)

	return result
}
