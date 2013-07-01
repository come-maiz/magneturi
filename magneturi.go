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
// See the schema overview at:
//     http://magnet-uri.sourceforge.net/magnet-draft-overview.txt
package magneturi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	magnetURISchemaPrefix = "magnet:?"
	exactTopicPrefix      = "xt"
	displayNamePrefix     = "dn"
	keywordTopicPrefix    = "kt"
	manifestTopicPrefix   = "mt"
)

// MagnetURI represents a uniform resource identifier following the magnet scheme.
type MagnetURI struct {
	Parameters []Parameter
}

// Parameter represents a parameter in a Magnet URI.
type Parameter struct {
	Prefix string
	Index  int // 0 means there is no index specified for the parameter.
	Value  string
}

// ExactTopics returns the list of exact topic parameters of the Magnet URI.
func (magnetURI *MagnetURI) ExactTopics() []Parameter {
	return magnetURI.parametersByPrefix(exactTopicPrefix)
}

func (magnetURI *MagnetURI) parametersByPrefix(prefix string) []Parameter {
	prefixParameters := make([]Parameter, 0, len(magnetURI.Parameters))
	for _, parameter := range magnetURI.Parameters {
		if parameter.Prefix == prefix {
			prefixParameters = append(prefixParameters, parameter)
		}
	}
	return prefixParameters
}

// DisplayNames returns the list of display name parameters of the Magnet URI.
func (magnetURI *MagnetURI) DisplayNames() []Parameter {
	return magnetURI.parametersByPrefix(displayNamePrefix)
}

// KeywordTopics returns the list of keyword topic parameters of the Magnet URI.
func (magnetURI *MagnetURI) KeywordTopics() []Parameter {
	return magnetURI.parametersByPrefix(keywordTopicPrefix)
}

// ManifestTopics returns the list of manifest topic parameters of the Magnet URI.
func (magnetURI *MagnetURI) ManifestTopics() []Parameter {
	return magnetURI.parametersByPrefix(manifestTopicPrefix)
}

// Equal returns true if the Magnet URIs are equal, false if not.
// The order of the parameters is not important.
func (magnetURI MagnetURI) Equal(x MagnetURI) bool {
	return compareParameters(magnetURI.Parameters, x.Parameters)
}

func compareParameters(first []Parameter, second []Parameter) bool {
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

func containsParameter(list []Parameter, parameter Parameter) bool {
	for _, element := range list {
		if parameter.Prefix == element.Prefix &&
			parameter.Index == element.Index &&
			parameter.Value == element.Value {
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

func parseParameters(parameters []string) (magnetURI MagnetURI, err error) {
	for _, parameter := range parameters {
		magnetURI, err = parseParameter(parameter, magnetURI)
		if err != nil {
			magnetURI = MagnetURI{}
		}
	}
	return
}

func parseParameter(parameter string, magnetURI MagnetURI) (MagnetURI, error) {
	parameterSplit := strings.SplitN(parameter, "=", 2)
	if len(parameterSplit) != 2 {
		return MagnetURI{}, errors.New(
			fmt.Sprintf("Parameter without prefix: %q", parameter))
	}
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

func splitPrefixIndex(prefix string) (string, int, error) {
	if strings.Contains(prefix, ".") {
		prefixSplit := strings.SplitN(prefix, ".", 2)
		index, err := strconv.Atoi(prefixSplit[1])
		if err != nil {
			return "", index, err
		}
		return prefixSplit[0], index, nil
	}
	return prefix, 0, nil
}

func addParameterToMagnetURI(prefix string, index int, value string, magnetURI MagnetURI) (MagnetURI, error) {
	if !isValidPrefix(prefix) {
		return MagnetURI{}, errors.New(
		    fmt.Sprintf("Unknown parameter prefix: %q", prefix))
	}
	var parameter = Parameter{prefix, index, value}
	magnetURI.Parameters = append(magnetURI.Parameters, parameter)
	return magnetURI, nil
}

func isValidPrefix(prefix string) bool {
	return prefix == exactTopicPrefix || prefix == displayNamePrefix ||
		prefix == keywordTopicPrefix || prefix == manifestTopicPrefix
}

// String reassembles the MagnetURI into a valid MagnetURI string.
func (magnetURI *MagnetURI) String() (string, error) {
	var s string
	var err error = nil
	if !magnetURI.hasParameters() {
		err = errors.New("The Magnet URI has no parameters.")
	} else {
		s = magnetURISchemaPrefix
		parameters := magnetURI.parameterStrings()
		s += strings.Join(parameters, "&")
	}
	return s, err
}

func (magnetURI *MagnetURI) hasParameters() bool {
	if len(magnetURI.ExactTopics()) != 0 ||
		len(magnetURI.DisplayNames()) != 0 ||
		len(magnetURI.KeywordTopics()) != 0 ||
		len(magnetURI.ManifestTopics()) != 0 {
		return true
	}
	return false
}

func (magnetURI *MagnetURI) parameterStrings() []string {
	parameterStrings := make([]string, 0, len(magnetURI.Parameters))
	for _, parameter := range magnetURI.Parameters {
		parameterStrings = append(parameterStrings, parameter.String())
	}
	return parameterStrings
}

// String reassembles the Parameter into a valid MagnetURI parameter string.
func (parameter *Parameter) String() string {
	if parameter.Index != 0 {
		return fmt.Sprintf(
			"%s.%d=%s", parameter.Prefix, parameter.Index, parameter.Value)
	}
	return fmt.Sprintf("%s=%s", parameter.Prefix, parameter.Value)
}
