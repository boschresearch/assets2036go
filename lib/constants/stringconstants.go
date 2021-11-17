// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package constants

const (
	// TopicSeparator dfg
	TopicSeparator = "/"
	// SubmodelEndpointName kg
	SubmodelEndpointName = "_endpoint"
	// PropertyNameMeta kjh
	PropertyNameMeta = "_meta"
	// PropertyNameMetaSubmodelURL  kgh kj
	PropertyNameMetaSubmodelURL = "submodel_url"
	// PropertyNameMetaSubmodelSchema kgu
	PropertyNameMetaSubmodelSchema = "submodel_definition"
	// PropertyNameMetaSource sdf
	PropertyNameMetaSource = "source"
	// PropertyNameHealthy sfd
	PropertyNameHealthy = "healthy"
	// PropertyNameOnline sef
	PropertyNameOnline = "online"
	// TopicElementSubmodelOperationReq is jhg
	TopicElementSubmodelOperationReq = "REQ"
	// TopicElementSubmodelOperationResp jhg
	TopicElementSubmodelOperationResp = "RESP"
	// PayloadPropReqID k
	PayloadPropReqID = "req_id"
	// PayloadPropOpReqParams jhg
	PayloadPropOpReqParams = "params"
	// PayloadPropOpResp kjh
	PayloadPropOpResp = "resp"
	// PayloadPropOpRespError kjh
	PayloadPropOpRespError = "error"
	// PayloadPropEvTimestamp kjh
	PayloadPropEvTimestamp = "timestamp"
	// EnpointSubmodelURL is the url to the endpoint submodel description
	EnpointSubmodelURL = "https://arena2036-infrastructure.saz.bosch-si.com/arena2036_public/assets2036_submodels/-/raw/master/_endpoint.json"
	// EndpointOperationNameShutdown is the name of the shutdown operation
	EndpointOperationNameShutdown = "shutdown"
	// EndpointOperationNamePing is the name of the ping operation
	EndpointOperationNamePing = "ping"
	// EndpointOperationNameRestart is the name of the restart operation
	EndpointOperationNameRestart = "restart"
	// EnvVarAssets2036SubmodelsOverwrite = "ASSETS2036_SUBMODELS_OVERWRITE"
	EnvVarAssets2036SubmodelsOverwrite = "ASSETS2036_SUBMODELS_OVERWRITE"
	// SubmodelEndpoint is the endpoint submodel directly
	SubmodelEndpoint = `
	{
		"name": "_endpoint",
		"revision": "0.0.3",
		"description": "Meldet Online- und Healthy-Status eines Assets",
		"properties": {
		  "online": {
			"description": "Asset ist online, der Adapter ist erreichbar",
			"type": "boolean"
		  },
		  "healthy": {
			"description": "Asset ist healthy, der Adapter ist bereit",
			"type": "boolean"
		  }
		},
		"events": {
		  "log": {
			"description": "Logging Event",
			"parameters": {
			  "entry": {
				"description": "Logging Text",
				"type": "string"
			  }
			}
		  }
		},
		"operations": {
		  "shutdown": {
			"description": "Asset ausschalten",
			"parameters": {}
		  },
		  "restart": {
			"description": "Asset neu starten",
			"parameters": {}
		  },
		  "ping": {
			"description": "Ping Asset",
			"parameters": {}
		  }
		}
	  }
	`
)
