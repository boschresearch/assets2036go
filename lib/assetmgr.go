// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boschresearch/assets2036go/lib/constants"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
)

type messageHandler func(topic Topic, rawPayload []byte)

// JSONObj is just an alias for the typical JSON obj representation map[string]interface{}
type JSONObj map[string]interface{}

// The AssetMgr is the basic access point to the assets2036go library. Use it
// to create assets and assetproxys for your assets environment
type AssetMgr struct {
	Host                  string
	Port                  uint16
	EndpointName          string
	DefaultNamespace      string
	EndpointAsset         *Asset
	ExitWhenConnObsvFails bool
	ConnObsvTimeout       time.Duration

	mqttClient          mqtt.Client
	operationResponses  map[string]*submodelOperationResponse
	messageHandlers     map[string]messageHandler
	chanStopEndpObsv    chan bool
	chanStopHealthyObsv chan bool
	healthyCallback     func() bool
	subscribedTopics    map[string]byte
}

// CreateAssetMgr creates an instance of AssetMgr and trys to connect it to the
// given MQTT Server on host:port.
func CreateAssetMgr(host string, port uint16, defaultNamespace, endpointName string, observeConnection bool) (*AssetMgr, error) {
	mgr := &AssetMgr{
		Host:                  host,
		Port:                  port,
		EndpointName:          endpointName,
		DefaultNamespace:      defaultNamespace,
		ExitWhenConnObsvFails: true,
		ConnObsvTimeout:       20 * time.Second,
		operationResponses:    map[string]*submodelOperationResponse{},
		messageHandlers:       map[string]messageHandler{},
		subscribedTopics:      map[string]byte{},
	}

	// connect to mqtt broker
	err := mgr.connect(host, port)
	if err != nil {
		return nil, err
	}

	// create stop signal channel for observation goroutine
	mgr.chanStopEndpObsv = make(chan bool)

	// start the goroutine to observe the endpoint
	mgr.observeMqttConnection(mgr.chanStopEndpObsv)

	return mgr, nil
}

// SetHealthyCallback lets you set a callback function to get some custom healthy information. This
// healthy infvormation will be publish in the AssetMgr's endpoint asset's submodel _endpoint.
func (mgr *AssetMgr) SetHealthyCallback(cb func() bool) {
	if cb != nil {

		// stop existing goroutine, if there
		if mgr.healthyCallback != nil {
			mgr.chanStopHealthyObsv <- true
		}

		mgr.healthyCallback = cb
		mgr.chanStopHealthyObsv = make(chan bool)

		go mgr.observeHealthy(mgr.chanStopHealthyObsv)
	}
}

// GetEndpoint return the submodel _endpoint of the automatically created endpoint asset
func (mgr *AssetMgr) GetEndpoint() *Submodel {
	if mgr.EndpointAsset == nil {
		epa, err := mgr.createEndpointAsset(mgr.DefaultNamespace, mgr.EndpointName)
		mgr.EndpointAsset = epa

		if err != nil {
			logg().Printf("Error: Create Endpoint Asset failed: %v", err)
			return nil
		}
	}
	return mgr.EndpointAsset.Submodel(constants.SubmodelEndpointName)
}

// SetShutdownCallback lets you set your local method when your assetmgrendpoint is forced to shutdown
func (mgr *AssetMgr) SetShutdownCallback(cb func()) {
	mgr.GetEndpoint().Operation(constants.EndpointOperationNameShutdown).Callback = func(req *SubmodelOperationRequest) (resp interface{}, err error) {
		cb()
		return nil, nil
	}
}

func (mgr *AssetMgr) connect(host string, port uint16) error {
	clientOptions := mqtt.NewClientOptions()
	url := "tcp://" + host + ":" + strconv.FormatInt(int64(port), 10)

	newUUID := uuid.NewV4()
	clientID := "assets2036go_" + newUUID.String() + "_" + mgr.EndpointName
	clientOptions.AddBroker(url)
	clientOptions.SetClientID(clientID)

	jsonFalse, _ := json.Marshal(false)
	clientOptions.SetWill(
		buildTopic(
			mgr.DefaultNamespace,
			mgr.EndpointName,
			constants.SubmodelEndpointName,
			constants.PropertyNameOnline).String(), string(jsonFalse), 2, true)
	//	clientOptions.SetKeepAlive(10 * time.Second)

	logg().Printf("Try connect to %v, clientID: %v", url, clientID)

	mgr.mqttClient = mqtt.NewClient(clientOptions)

	token := mgr.mqttClient.Connect()

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("unable to connect to MQTT broker at %v:%v", host, port)
	}

	return nil
}

func (mgr *AssetMgr) createAssetBaseInt(namespace, name string, submodelUrls ...string) (asset *Asset, err error) {

	// preconditions:
	if mgr.mqttClient == nil {
		return nil, errors.New("mqttClient not available. Run 'connect()' first")
	}

	if namespace == "" {
		return nil, errors.New("need a non empty domain string")
	}

	if name == "" {
		return nil, errors.New("need a non empty name string")
	}

	asset = &Asset{
		assetElement: assetElement{
			assetMgr: mgr,
		},
		Submodels: map[string]*Submodel{},
		Name:      name,
		Namespace: namespace,
	}

	for _, element := range submodelUrls {
		jsonString, err := GetJSON(element)
		if err != nil {
			logg().Printf("Submodel at URL \"%v\" not found", element)
			continue
		}

		// Parse submodel and create properties, events, operations
		var submodel Submodel
		json.Unmarshal([]byte(jsonString), &submodel)

		// parse it again into a generic map --> needed for the endpoint property submodel_definition
		json.Unmarshal([]byte(jsonString), &submodel.submodelSchema)

		if submodel.Properties == nil {
			submodel.Properties = make(map[string]*SubmodelProperty)
		}

		if submodel.Events == nil {
			submodel.Events = make(map[string]*SubmodelEvent)
		}

		if submodel.Operations == nil {
			submodel.Operations = make(map[string]*SubmodelOperation)
		}

		if submodel.Name == "" {
			logg().Printf(
				"WARNING: Name file of submodel at url %v is empty. Probably something is wrong!", element)
		}

		if len(submodel.Operations) == 0 && len(submodel.Properties) == 0 && len(submodel.Events) == 0 {
			logg().Printf(
				"Properties, Operations and Events of submodel %v of asset %v at url %v are empty. Something wrong?",
				submodel.Name,
				name,
				element)
		}

		asset.Submodels[submodel.Name] = &submodel

		// build up hierarchy for topic generation
		submodel.injectAssetMgr(mgr)

		// set topic properties --> unmarshalling sets the map keys in an operation, property or event, but not the name property in the operation instance itself
		submodel.setTopicsRecursively(namespace, name)

		submodel.submodelURL = element
	}

	return asset, nil
}

// func (mgr *AssetMgr) resubscribeAll() {
// 	for k := range mgr.subscribedTopics {
// 		mgr.mqttClient.Unsubscribe(k)
// 	}

// 	mgr.mqttClient.SubscribeMultiple(mgr.subscribedTopics, mgr.messageReceived)
// }

func (mgr *AssetMgr) createAssetInt(namespace, name string, submodelUrls ...string) (*Asset, error) {
	asset, err := mgr.createAssetBaseInt(namespace, name, submodelUrls...)
	if err != nil {
		return nil, err
	}

	for _, submodel := range asset.Submodels {

		// create meta property
		submodel.Properties[constants.PropertyNameMeta] = &SubmodelProperty{
			assetElement: assetElement{
				Topic:    TopicFromElements(namespace, name, submodel.Name, constants.PropertyNameMeta),
				assetMgr: mgr,
			},
		}

		// store relevant topics to subscribe later
		for _, operation := range submodel.Operations {
			operation.asset = asset
			mgr.subscribedTopics[buildTopic(operation.Topic.String(), constants.TopicElementSubmodelOperationReq).String()] = 2
			operation.installHandlersAsset(mgr)
		}

		for _, property := range submodel.Properties {
			property.asset = asset
		}

		for _, ev := range submodel.Events {
			ev.asset = asset
		}

		asset.asset = asset

		// publish submodel's meta information:
		submodel.Properties[constants.PropertyNameMeta].SetValue(map[string]interface{}{
			constants.PropertyNameMetaSubmodelURL:    submodel.submodelURL,
			constants.PropertyNameMetaSubmodelSchema: submodel.submodelSchema,
			constants.PropertyNameMetaSource:         mgr.EndpointName,
		})
	}

	if len(mgr.subscribedTopics) == 0 {
		logg().Printf("Ups")
	}

	logg().Printf("Subscribe to \n%v", strings.Join(keysFromStringByte(mgr.subscribedTopics), "\n"))

	for k := range mgr.subscribedTopics {
		mgr.mqttClient.Unsubscribe(k)
	}

	mgr.mqttClient.SubscribeMultiple(mgr.subscribedTopics, mgr.messageReceived)

	return asset, nil
}

// CreateAsset creates an asset which is implemented on this host. If you want to _use_ an existing
// asset, use CreateAssetProxy because in that case a proxy is what you want.
func (mgr *AssetMgr) CreateAsset(namespace, name string, submodelUrls ...string) (*Asset, error) {

	// To make sure it is created:
	_ = mgr.GetEndpoint()

	return mgr.createAssetInt(namespace, name, submodelUrls...)
}

// CreateAssetProxy creates a local proxy to a remotely managed asset
// to make use of its properties, operations and events without bothering about
// where it is implemented
func (mgr *AssetMgr) CreateAssetProxy(namespace, name string, submodelUrls ...string) (*Asset, error) {

	assetBase, err := mgr.createAssetBaseInt(namespace, name, submodelUrls...)
	if err != nil {
		return nil, err
	}

	// topicsToSubscribe := map[string]byte{}

	for _, submodel := range assetBase.Submodels {
		// store relevant topics to subscribe later
		for _, operation := range submodel.Operations {
			operation.assetProxy = assetBase
			mgr.subscribedTopics[buildTopic(operation.Topic.String(), constants.TopicElementSubmodelOperationResp).String()] = 2
			operation.installHandlersProxy(mgr)
		}

		for _, property := range submodel.Properties {
			property.assetProxy = assetBase
			mgr.subscribedTopics[property.Topic.String()] = 2
			property.installHandlers(mgr)
		}

		for _, ev := range submodel.Events {
			ev.assetProxy = assetBase
			mgr.subscribedTopics[ev.Topic.String()] = 2
			ev.installHandlers(mgr)
		}
	}

	logg().Printf("Subscribe to \n%v", strings.Join(keysFromStringByte(mgr.subscribedTopics), "\n"))

	for k := range mgr.subscribedTopics {
		mgr.mqttClient.Unsubscribe(k)
	}

	mgr.mqttClient.SubscribeMultiple(mgr.subscribedTopics, mgr.messageReceived)

	return assetBase, nil
}

func (mgr *AssetMgr) createEndpointAsset(namespace, endpointName string) (*Asset, error) {
	asset, err := mgr.createAssetInt(namespace, endpointName, constants.EnpointSubmodelURL)

	if err != nil {
		return nil, err
	}

	ep := asset.Submodel(constants.SubmodelEndpointName)

	ep.Operation(constants.EndpointOperationNameShutdown).Callback = func(req *SubmodelOperationRequest) (resp interface{}, err error) {
		defer os.Exit(0)
		return nil, nil
	}

	ep.Operation(constants.EndpointOperationNamePing).Callback = func(req *SubmodelOperationRequest) (resp interface{}, err error) {
		return nil, nil
	}

	ep.Property(constants.PropertyNameOnline).SetValue(true)

	return asset, nil
}

func (mgr *AssetMgr) publish(topic string, message []byte, retained bool) (err error) {
	logg().Printf("publish: %v to %v", string(message), topic)

	token := mgr.mqttClient.Publish(topic, 2, retained, message) // tbd: Hier kann man wohl auch auf en versand warten! --> Gut für die Verbindungssicherheit
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (mgr *AssetMgr) messageReceived(client mqtt.Client, message mqtt.Message) {
	logg().Printf("received: %v at %v", string(message.Payload()), message.Topic())

	handler := mgr.messageHandlers[message.Topic()]
	if handler != nil {
		go handler(TopicFromStr(message.Topic()), message.Payload())
	}
}

func (mgr *AssetMgr) observeMqttConnection(stopSignal chan bool) {
	proxy, err := mgr.CreateAssetProxy(
		mgr.DefaultNamespace,
		mgr.EndpointName,
		constants.EnpointSubmodelURL,
	)

	if err != nil {
		logg().Fatal("Unable to create asset proxy for endpoint observance. Connection observation inactive!")
		return
	}

	go func(mgr *AssetMgr, endpoint *Submodel) {
		for {
			select {
			case <-mgr.chanStopEndpObsv:
				logg().Println("Got exit signal. Connection observation stops!")
				return
			case <-time.After(10 * time.Second):
				_, err := endpoint.Operation(constants.EndpointOperationNamePing).InvokeTimeout(map[string]interface{}{}, 5*time.Second)
				if err != nil {
					logg().Printf("Ping in connection observation returned error: %v", err)
					if mgr.ExitWhenConnObsvFails {
						os.Exit(1)
					}
				} else {
					logg().Println("Ping returned. Connection observation successful!")
				}
			}
		}
	}(mgr, proxy.Submodel(constants.SubmodelEndpointName))
}

func (mgr *AssetMgr) observeHealthy(chanStop chan bool) {
	for {
		select {
		case <-chanStop:
			return
		case <-time.After(20 * time.Second):
			if mgr.EndpointAsset != nil {
				if mgr.healthyCallback == nil {
					mgr.EndpointAsset.Submodel(constants.SubmodelEndpointName).Property(constants.PropertyNameHealthy).SetValue(true)
				} else {
					mgr.EndpointAsset.Submodel(constants.SubmodelEndpointName).Property(constants.PropertyNameHealthy).SetValue(mgr.healthyCallback())
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}

// CloseConnections kjh
func (mgr *AssetMgr) CloseConnections() {
	mgr.chanStopEndpObsv <- true

	if mgr.chanStopHealthyObsv != nil {
		mgr.chanStopHealthyObsv <- true
	}

	if mgr.EndpointAsset != nil {
		mgr.EndpointAsset.Submodel(constants.SubmodelEndpointName).Property(constants.PropertyNameOnline).SetValue(false)
	}

	mgr.mqttClient.Disconnect(1000)
}
