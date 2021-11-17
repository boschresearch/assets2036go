// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

type iSubmodelElement interface {
	injectAssetMgr(assetMgr *AssetMgr)

	injectAsset(asset *Asset)

	injectAssetProxy(assetProxy *Asset)

	setTopic(t string)

	Name() string
}
