# assets2036go

## Table of content 
- [assets2036go](#assets2036go)
  - [Table of content](#table-of-content)
  - [Description](#description)
  - [Getting started](#getting-started)
    - [Using the lib via sources](#using-the-lib-via-sources)
    - [Using the lib via golangs dependency system](#using-the-lib-via-golangs-dependency-system)
  - [Dependencies](#dependencies)
  - [Remarks](#remarks)
  - [Authors](#authors)
  - [License](#license)
  - [Acknowledgments](#acknowledgments)

## Description

Use **assets2036go** to your go-app communicate with other assets following the new leightweight asset administration shell standard using MQTT. 

## Getting started

### Using the lib via sources 

copy the directory assets2036go/assets2036go to any target directory you want to have it [myAssetsLibDir]. 

Then in the go.mod of your project add the line 
replace bosch.com/assets2036go v0.1.9 => [myAssetsLibDir]. 

As an example see the go.mod of [assets2036go/example/go.mod](./example/go.mod)


### Using the lib via golangs dependency system

Probably the best entry point is the example program example_1. You will learn how to create an asset, implement a submodel operation, set submodel properties, then create an asset proxy, call the operation or read the property. 

For more example you can have a look into the unit tests in the package assets2036gotest. 

For questions or problems: thomas.jung6@de.bosch.com

## Dependencies

- Eclipse Paho MQTT Go client	1.2.0	Eclipse Public License - v 1.0
- UUID package for Go language	1.2.0	MIT


## Remarks

## Authors

[Thomas Jung](mailto:thomas.jung6@de.bosch.com)

## License

assets2036go is open-sourced under the Apache-2.0 license. See the LICENSE file for details.
For a list of other open source components included in UUV Simulator, see the file 3rd-party-licenses.txt.

## Acknowledgments

Thanks to [Daniel Ewert](https://github.com/DaEwe/) for the inspiration, conceptual work and preliminary python library. 
