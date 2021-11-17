// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

type assetElement struct {
	// Topic string
	Topic       Topic
	Description string `json:"description"`

	assetMgr   *AssetMgr
	asset      *Asset
	assetProxy *Asset
}

func (op *assetElement) injectAssetMgr(asssetMgr *AssetMgr) {
	op.assetMgr = asssetMgr
}

func (op *assetElement) injectAsset(asset *Asset) {
	op.asset = asset
}

func (op *assetElement) injectAssetProxy(assetProxy *Asset) {
	op.assetProxy = assetProxy
}

func (op *assetElement) setTopic(t string) {
	op.Topic = TopicFromStr(t)
}
