package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	assets2036go "github.com/boschresearch/assets2036go/lib"
)

// var mqttHost = "192.168.100.3"
var mqttHost = "test.mosquitto.org"
var mqttPort = uint16(1883)

// Very simple example showing how to create assets, assetproxies,
// implementing and calling operations, setting and reading properties.
func main() {
	var assetNamespace string = "assets2036gotest"
	var assetName string = "example_1_asset"

	// create asset manager:
	mgr, _ := assets2036go.CreateAssetMgr(mqttHost, mqttPort, assetNamespace, "example_1_endpointasset", true)
	defer mgr.CloseConnections()

	// set url override if no repo is available
	os.Setenv("ASSETS2036_SUBMODELS_OVERWRITE", "file:///./models")

	// create an asset:
	example_1, err := mgr.CreateAsset(assetNamespace, assetName, "https://submodelrepo.bosch.com/testmodel.json")

	if err != nil {
		log.Fatal(err)
	}

	// implement a submodel operation
	example_1.Submodel("testmodel").Operation("getObject").Callback = func(req *assets2036go.SubmodelOperationRequest) (resp interface{}, err error) {

		log.Print("Handle operation getObject() request!")

		result := map[string]interface{}{
			"name":   "Walter",
			"age":    99,
			"weight": 70,
			"arrayOfInteger": []int{
				1, 2, 3, 4,
			},
		}

		return result, nil
	}

	// create an asset proxy:
	example_1_proxy, err := mgr.CreateFullAssetProxy(assetNamespace, assetName, 2*time.Second)

	if err != nil {
		log.Print(err)
		return
	}

	// call the remotly implemented operation of the asset with empty parameter set map[string]interface{}{}:
	result, err := example_1_proxy.Submodel("testmodel").Operation("getObject").Invoke(map[string]interface{}{})

	if err != nil {
		log.Print(err)
		return
	}

	// print result
	resultTyped := result.(map[string]interface{})
	log.Print(resultTyped)

	// set and get property:
	rand.Seed(time.Now().Unix())
	randomnumber := rand.Intn(1000)

	example_1.Submodel("testmodel").Property("number").SetValue(randomnumber)

	// wait for transmission
	time.Sleep(2 * time.Second)

	// read property from proxy:
	read_number := int(example_1_proxy.Submodel("testmodel").Property("number").Value().(float64))

	log.Printf("Number in asset: %v, number read from proxy: %v", randomnumber, read_number)
}
