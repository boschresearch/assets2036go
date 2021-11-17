// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"strings"

	"github.com/boschresearch/assets2036go/lib/constants"
)

// Topic is an assets2036 mqtt topic
type Topic []string

// TopicFromStr create a new topic object
func TopicFromStr(str string) Topic {
	return Topic(strings.Split(string(str), constants.TopicSeparator))
}

// TopicFromElements does stuff
func TopicFromElements(elements ...string) Topic {
	return Topic(elements)
}

// AssetAbs kjasdkajs
func (topic Topic) AssetAbs() string {
	return buildTopic(topic[0], topic[1]).String()
}

// Asset kjasdkajs
func (topic Topic) Asset() string {
	return topic[1]
}

// Namespace kjasdkajs
func (topic Topic) Namespace() string {
	return topic[0]
}

// Submodel asdas
func (topic Topic) Submodel() string {
	return topic[2]
}

// SubmodelElement ad asf
func (topic Topic) SubmodelElement() string {
	return topic[3]
}

func (topic Topic) String() string {
	return strings.Join(topic, constants.TopicSeparator)
}

func buildTopic(elements ...string) Topic {
	return Topic{strings.Join(elements, constants.TopicSeparator)}
}

// StringList returns simetins
func StringList(topics []Topic) string {
	result := ""
	for _, t := range topics {
		result += t.String()
		result += "\n"
	}

	return result
}
