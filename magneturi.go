// Copyright 2013.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package magneturi parses Magnet URIs.
// See the schema overview at
// http://magnet-uri.sourceforge.net/magnet-draft-overview.txt

package magneturi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	magnetURISchemaPrefix = "magnet:?"
	exactTopicPrefix = "xt"
	displayNamePrefix = "dn"
	keywordTopicPrefix = "kt"
	manifestTopicPrefix = "mt"
)

type MagnetURI struct {
	ExactTopics    []string // xt
	DisplayNames   []string // dn
	KeywordTopics  []string // kt
	ManifestTopics []string // mt
}

// Equal returns true if the Magnet URIs are equal, false if not.
// The order of the parameters is not important.
func (magnetURI MagnetURI) Equal(x MagnetURI) bool {
	return compareParameters(magnetURI.ExactTopics, x.ExactTopics) &&
		compareParameters(magnetURI.DisplayNames, x.DisplayNames) &&
		compareParameters(magnetURI.KeywordTopics, x.KeywordTopics) &&
		compareParameters(magnetURI.ManifestTopics, x.ManifestTopics)
}

func compareParameters(first []string, second []string) bool {
	if len(first) == len(second) {
		for _, parameter := range first {
			if !containsParameter(second, parameter) {
				return false
			}
		}
		return true
	}
	return false
}

func containsParameter(list []string, parameter string) bool {
	for _, element := range list {
		if parameter == element {
			return true
		}
	}
	return false
}

// Parse parses a raw Magnet URI string into a MagnetURI structure.
func Parse(rawMagnetURI string) (MagnetURI, error) {
	if strings.HasPrefix(rawMagnetURI, magnetURISchemaPrefix) {
		rawMagnetURIWithoutPrefix := strings.TrimPrefix(
			rawMagnetURI, magnetURISchemaPrefix)
		parameters := strings.Split(rawMagnetURIWithoutPrefix, "&")
		return parseParameters(parameters)
	}
	return MagnetURI{}, errors.New(
		fmt.Sprintf(
			"The string doesn't start with the Magnet URI schema prefix %q",
			magnetURISchemaPrefix))
}

func parseParameters(parameters []string) (MagnetURI, error) {
	magnetURI := MagnetURI{}
	var err error = nil
	for _, parameter := range parameters {
		magnetURI, err = parseParameter(parameter, magnetURI)
		if err != nil {
			magnetURI = MagnetURI{}
		}
	}
	return magnetURI, err
}

func parseParameter(
	parameter string, magnetURI MagnetURI) (MagnetURI, error) {
	parameterSplit := strings.SplitN(parameter, "=", 2)
	if len(parameterSplit) == 2 {
		prefix := parameterSplit[0]
		prefix, index, err := splitPrefixIndex(prefix)
		if err != nil {
			return MagnetURI{}, errors.New(
				fmt.Sprintf(
      				"Wrong parameter prefix: %q; %s", prefix, err.Error()))
		}
		value := parameterSplit[1]
		return addParameterToMagnetURI(prefix, index, value, magnetURI)
	}
	return MagnetURI{}, errors.New(
		fmt.Sprintf("Parameter without prefix: %q", parameter))
}

func splitPrefixIndex(prefix string) (string, int, error) {
	if strings.Contains(prefix, ".") {
		prefixSplit := strings.SplitN(prefix, ".", 2)
		index, err := strconv.Atoi(prefixSplit[1])
		return prefixSplit[0], index, err
	}
	return prefix, 1, nil
}

func addParameterToMagnetURI(
	prefix string, index int, value string, magnetURI MagnetURI) (
	MagnetURI, error) {
	switch {
		case prefix == exactTopicPrefix:
			magnetURI.ExactTopics = insertParameter(
				index, value, magnetURI.ExactTopics)
			return magnetURI, nil
		case prefix == displayNamePrefix:
			magnetURI.DisplayNames = append(
				magnetURI.DisplayNames, value)
			return magnetURI, nil
		case prefix == keywordTopicPrefix:
			magnetURI.KeywordTopics = append(
				magnetURI.KeywordTopics, value)
			return magnetURI, nil
		case prefix == manifestTopicPrefix:
			magnetURI.ManifestTopics = append(
				magnetURI.ManifestTopics, value)
			return magnetURI, nil
		default:
			return MagnetURI{}, errors.New(
				fmt.Sprintf("Unknown parameter prefix: %q", prefix))
	}
}

func insertParameter(
	index int, value string, parametersList []string) []string {
	if cap(parametersList) < index {
		newParameterList := make([]string, index, index)
		copy(newParameterList, parametersList)
		parametersList = newParameterList
	}
	parametersList[index-1] = value
	return parametersList
}

// String reassembles the MagnetURI into a valid MagnetURI string.
func (magnetURI *MagnetURI) String() (string, error) {
	var s string = ""
	var err error = nil
	if magnetURI.hasParameters() {
		s = magnetURISchemaPrefix
		parameters := getParametersWithPrefix(
			exactTopicPrefix, magnetURI.ExactTopics)
		parameters = append(parameters, getParametersWithPrefix(
			displayNamePrefix, magnetURI.DisplayNames)...)
		parameters = append(parameters, getParametersWithPrefix(
			keywordTopicPrefix, magnetURI.KeywordTopics)...)
		parameters = append(parameters, getParametersWithPrefix(
			manifestTopicPrefix, magnetURI.ManifestTopics)...)
		s += strings.Join(parameters, "&")
	} else {
		err = errors.New("The Magnet URI has no parameters.")
	}
	return s, err
}

func (magnetURI *MagnetURI) hasParameters() bool {
	if magnetURI.ExactTopics != nil || magnetURI.DisplayNames != nil ||
		magnetURI.KeywordTopics != nil || magnetURI.ManifestTopics != nil {
		return true
	}
	return false
}

func getParametersWithPrefix(prefix string, parameters []string) []string {
	numberOfTopics := len(parameters)
	parametersStrings := make([]string, 0, numberOfTopics)
	if numberOfTopics == 1 {
		parameterString := fmt.Sprintf("%s=%s", prefix, parameters[0])
		parametersStrings = append(parametersStrings, parameterString)
	} else {
		for index, parameter := range parameters {
			parameterString := fmt.Sprintf(
				"%s.%d=%s", prefix, index+1, parameter)
			parametersStrings = append(parametersStrings, parameterString)
		}
	}
	return parametersStrings
}
