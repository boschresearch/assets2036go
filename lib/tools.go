// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/boschresearch/assets2036go/lib/constants"
)

func GetJSON(url string) (string, error) {
	logg().Printf("getJSON(%v)", url)

	submodelsOverwrite := os.Getenv(constants.EnvVarAssets2036SubmodelsOverwrite)
	if submodelsOverwrite != "" {
		urlParts := strings.Split(url, "/")
		url = fmt.Sprintf("%v/%v", submodelsOverwrite, urlParts[len(urlParts)-1])

		logg().Printf("env var ASSETS2036_SUBMODELS_OVERWRITE set, so forward to: getJSON(%v)", url)
	}

	if url == constants.EnpointSubmodelURL {
		return constants.SubmodelEndpoint, nil
	}

	url = strings.TrimSpace(strings.ToLower(url))
	if url[0:8] == "file:///" {
		bytes, err := ioutil.ReadFile(url[8:])

		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	return string(body), nil
}

func keysFromStringByte(m map[string]byte) []string {
	result := make([]string, len(m))

	i := 0
	for k := range m {
		result[i] = k
		i++
	}

	return result
}

// FormatTimestamp return a standard conform string representation of the given timestamp
func FormatTimestamp(ts time.Time) string {
	return ts.Format(time.RFC3339)
}
