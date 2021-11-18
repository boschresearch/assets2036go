# assets2036go

## Table of content 
- [assets2036go](#assets2036go)
  - [Table of content](#table-of-content)
  - [Description](#description)
  - [Getting started](#getting-started)
  - [Dependencies](#dependencies)
  - [Remarks](#remarks)
  - [Authors](#authors)
  - [License](#license)
  - [Acknowledgments](#acknowledgments)

## Description

Use **assets2036go** to communicate with other assets following the assets2036 standard used on [Arena2036](https://www.arena2036.de/).

assets2036 is based on very lean MQTT and JSON conventions - you can even participate using only standard MQTT and JSON libs. This is a convenience library, simplifying participation with any go software. 

## Getting started

Simply add the line

    require github.com/boschresearch/assets2036go v0.3.2

to your go.mod file. 

Enter 

    go mod download

to get the lib sources. 

See ./example/example_1.go for details of usage. 

## Dependencies

- Eclipse Paho MQTT Go client	1.2.0	Eclipse Public License - v 1.0
- UUID package for Go language	1.2.0	MIT


## Remarks

## Authors

[Thomas Jung](https://github.com/thomasjosefjung)

## License

assets2036go is open-sourced under the Apache-2.0 license. See the LICENSE file for details.
For a list of other open source components assets2036go depends on, see [Dependencies](#dependencies). 

## Acknowledgments

Thanks to [Daniel Ewert](https://github.com/DaEwe/) for the inspiration, conceptual work and preliminary python library. 
