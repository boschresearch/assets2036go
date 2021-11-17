// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

// SubmodelOperationCallback is the signature of all functions, which can be used as callbacks in
// an asset implementation
type SubmodelOperationCallback func(req *SubmodelOperationRequest) (resp interface{}, err error)
