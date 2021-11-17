// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/boschresearch/assets2036go/lib/constants"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	uuid "github.com/satori/go.uuid"
)

func (mgr *AssetMgr) findAssetBy(topicToSubscribe Topic, filter func(topic string, payload []byte) bool) (assetTopics []Topic, err error) {
	clientOptions := mqtt.NewClientOptions()
	url := "tcp://" + mgr.Host + ":" + strconv.FormatInt(int64(mgr.Port), 10)

	newUUID := uuid.NewV4()
	clientID := "assets2036go_TempClient" + newUUID.String()
	clientOptions.AddBroker(url)
	clientOptions.SetClientID(clientID)

	logg().Printf("Try connect to %v, clientID: %v", url, clientID)

	client := mqtt.NewClient(clientOptions)
	defer client.Disconnect(1000)

	assets := map[string]bool{}

	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var latestReceive time.Time = time.Now()

	client.Subscribe(topicToSubscribe.String(), 0, func(client mqtt.Client, message mqtt.Message) {
		if !message.Retained() {
			return
		}

		structuredTopic := TopicFromStr(message.Topic())

		if filter != nil {
			if !filter(message.Topic(), message.Payload()) {
				return
			}
		}

		assets[structuredTopic.AssetAbs()] = true

		latestReceive = time.Now()
	})

	// observe receiving of retained messages until no more coming
	for time.Since(latestReceive) < 1000*time.Millisecond {
		time.Sleep(1 * time.Millisecond)
	}

	var chanFinished chan byte = make(chan byte)

	go func(latestReceive time.Time, chanFinished chan byte) {
		for time.Since(latestReceive) < 1000*time.Millisecond {
			time.Sleep(1 * time.Millisecond)
		}

		chanFinished <- 1
	}(latestReceive, chanFinished)

	// wait for obserer routine to signal finish
	<-chanFinished

	var result []Topic = make([]Topic, len(assets))

	i := 0
	for assetAbs := range assets {
		result[i] = TopicFromStr(assetAbs)
		i++
	}

	return result, nil
}

// FindAssetsBySubmodelName returns a slice of strings containing all those assets in the given namespace
// supporting the subbmodel as defined in subodelName parameter submodelURL
func (mgr *AssetMgr) FindAssetsBySubmodelName(submodelName string) (assetNames []Topic, err error) {
	return mgr.findAssetBy(buildTopic("+", "+", submodelName, "#"), nil)
}

// FindAssetsBySubmodelURL returns ...
func (mgr *AssetMgr) FindAssetsBySubmodelURL(submodelURL string) (assetNames []Topic, err error) {

	filter := func(topic string, payload []byte) bool {
		var metaObject map[string]interface{}
		json.Unmarshal(payload, &metaObject)

		if metaObject == nil {
			return false
		}

		if metaObject[constants.PropertyNameMetaSubmodelURL] == nil {
			return false
		}

		if metaObject[constants.PropertyNameMetaSubmodelURL].(string) != submodelURL {
			return false
		}

		return true
	}

	return mgr.findAssetBy(buildTopic("+", "+", "+", constants.PropertyNameMeta), filter)
}

// FindAssetsInNamespace returns ...
func (mgr *AssetMgr) FindAssetsInNamespace(namespace string) (assetNames []Topic, err error) {
	return mgr.findAssetBy(buildTopic(namespace, "+", "+", constants.PropertyNameMeta), nil)
}

// FindAllEndpointAssets returns ...
func (mgr *AssetMgr) FindAllEndpointAssets() (assetNames []Topic, err error) {
	filter := func(topic string, payload []byte) bool {
		var metaObject map[string]interface{}
		json.Unmarshal(payload, &metaObject)

		if metaObject == nil {
			return false
		}

		if metaObject[constants.PropertyNameMetaSubmodelURL] == nil {
			return false
		}

		return true
	}

	return mgr.findAssetBy(buildTopic("+", "+", constants.SubmodelEndpointName, constants.PropertyNameMeta), filter)
}

// RemoveTraces removes all traces (i.e. retained messages) of the given asset from the broker
func (mgr *AssetMgr) RemoveTraces(namespace, name string) {
	clientOptions := mqtt.NewClientOptions()
	url := "tcp://" + mgr.Host + ":" + strconv.FormatInt(int64(mgr.Port), 10)

	newUUID := uuid.NewV4()
	clientID := "assets2036go_TempClient" + newUUID.String()
	clientOptions.AddBroker(url)
	clientOptions.SetClientID(clientID)

	logg().Printf("Try connect to %v, clientID: %v", url, clientID)

	client := mqtt.NewClient(clientOptions)
	defer client.Disconnect(1000)

	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var latestReceive time.Time = time.Now()

	topicToSubscribe := fmt.Sprintf("%v/%v/#", namespace, name)

	topicsToClean := []Topic{}

	client.Subscribe(topicToSubscribe, 0, func(client mqtt.Client, message mqtt.Message) {
		if !message.Retained() {
			return
		}

		structuredTopic := TopicFromStr(message.Topic())
		topicsToClean = append(topicsToClean, structuredTopic)

		latestReceive = time.Now()
	})

	// observe receiving of retained messages until no more coming
	for time.Since(latestReceive) < 1000*time.Millisecond {
		time.Sleep(1 * time.Millisecond)
	}

	for _, topic := range topicsToClean {
		client.Publish(topic.String(), 2, true, nil)
	}
}

// GetSubmodelURLs returns the urls of all submodels supported by given asset
func (mgr *AssetMgr) GetSubmodelURLs(assetNamespace, assetName string, waitForSubmodels time.Duration) []string {
	clientOptions := mqtt.NewClientOptions()
	url := "tcp://" + mgr.Host + ":" + strconv.FormatInt(int64(mgr.Port), 10)

	newUUID := uuid.NewV4()
	clientID := "assets2036go_TempClient" + newUUID.String()
	clientOptions.AddBroker(url)
	clientOptions.SetClientID(clientID)

	logg().Printf("Try connect to %v, clientID: %v", url, clientID)

	client := mqtt.NewClient(clientOptions)
	defer client.Disconnect(1000)

	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var latestReceive time.Time = time.Now()

	topicToSubscribe := buildTopic(assetNamespace, assetName, "+", constants.PropertyNameMeta)

	result := make([]string, 0, 3)

	client.Subscribe(topicToSubscribe.String(), 0, func(client mqtt.Client, message mqtt.Message) {
		if !message.Retained() {
			return
		}

		var jsonObj JSONObj
		err := json.Unmarshal(message.Payload(), &jsonObj)

		if err != nil {
			logg().Printf("GetSubmodelURLs: Unable to unmarshal into JSONObj: %v", err)
			return
		}

		result = append(result, jsonObj[constants.PropertyNameMetaSubmodelURL].(string))

		latestReceive = time.Now()
	})

	var chanFinished chan byte = make(chan byte)

	go func(latestReceive time.Time, chanFinished chan byte) {
		for time.Since(latestReceive) < waitForSubmodels {
			time.Sleep(1 * time.Millisecond)
		}

		chanFinished <- 1
	}(latestReceive, chanFinished)

	// wait for obserer routine to signal finish
	<-chanFinished

	return result
}

// CreateFullAssetProxy create a proxy only if the asset  it is existing already. All submodels are directly read from MQTT
func (mgr *AssetMgr) CreateFullAssetProxy(namespace, name string, waitForSubmodels time.Duration) (asset *Asset, err error) {
	urls := mgr.GetSubmodelURLs(namespace, name, waitForSubmodels)
	return mgr.CreateAssetProxy(namespace, name, urls...)
}

// AddSubmodelObserver lets you add a callback method which will inbform you, if some new asset
// supporting the given submodelURL is created
func (mgr *AssetMgr) AddSubmodelObserver(submodelName string, handler func(namespace, name, submodelName string)) {
	mgr.mqttClient.Subscribe(buildTopic("+", "+", submodelName, "#").String(), 2, func(c mqtt.Client, m mqtt.Message) {
		t := TopicFromStr(m.Topic())
		handler(t.Namespace(), t.Asset(), t.Submodel())
	})
}
