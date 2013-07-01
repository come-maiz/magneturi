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

package magneturi

import (
	"testing"
)

func TestCompareParameters(t *testing.T) {
	scenarios := compareParametersScenarios
	for _, scenario := range scenarios {
		result := compareParameters(
			scenario.FirstParameters, scenario.SecondParameters)
		if result != scenario.ExpectedResult {
			t.Errorf(
				"Error on test %q: comparing %v and %v returns %t.",
				scenario.Name, scenario.FirstParameters,
				scenario.SecondParameters, result)
		}
	}
}

type compareParametersScenario struct {
	Name             string
	FirstParameters  []Parameter
	SecondParameters []Parameter
	ExpectedResult   bool
}

var compareParametersScenarios = []compareParametersScenario{
	{
		Name:             "Empty parameters",
		FirstParameters:  []Parameter{},
		SecondParameters: []Parameter{},
		ExpectedResult:   true,
	},
	{
		Name: "Multiple parameters",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param2"},
			Parameter{"pref", 0, "param3"},
		},
		SecondParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param2"},
			Parameter{"pref", 0, "param3"},
		},
		ExpectedResult: true,
	},
	{
		Name: "Parameters in different order",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param2"},
			Parameter{"pref", 0, "param3"},
		},
		SecondParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param3"},
			Parameter{"pref", 0, "param2"},
		},
		ExpectedResult: true,
	},
	{
		Name: "Missing parameter",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param2"},
		},
		SecondParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
		},
		ExpectedResult: false,
	},
	{
		Name: "Extra parameter",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param2"},
		},
		SecondParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
			Parameter{"pref", 0, "param3"},
			Parameter{"pref", 0, "param2"},
		},
		ExpectedResult: false,
	},
	{
		Name: "Wrong prefix",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
		},
		SecondParameters: []Parameter{
			Parameter{"wrong prefix", 0, "param1"},
		},
		ExpectedResult: false,
	},
	{
		Name: "Wrong index",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
		},
		SecondParameters: []Parameter{
			Parameter{"pref", 1, "param1"},
		},
		ExpectedResult: false,
	},
	{
		Name: "Wrong value",
		FirstParameters: []Parameter{
			Parameter{"pref", 0, "param1"},
		},
		SecondParameters: []Parameter{
			Parameter{"pref", 0, "wrong value"},
		},
		ExpectedResult: false,
	},
}

func TestCompareMagnetURIs(t *testing.T) {
	scenarios := compareMagnetURIsScenarios
	for _, scenario := range scenarios {
		result := scenario.FirstMagnetURI.Equal(scenario.SecondMagnetURI)
		if result != scenario.ExpectedResult {
			t.Errorf(
				"Error on test %q: comparing %v and %v returns %t.",
				scenario.Name, scenario.FirstMagnetURI,
				scenario.SecondMagnetURI, result)
		}
	}
}

type compareMagnetURIsScenario struct {
	Name            string
	FirstMagnetURI  MagnetURI
	SecondMagnetURI MagnetURI
	ExpectedResult  bool
}

var compareMagnetURIsScenarios = []compareMagnetURIsScenario{
	{
		Name:            "Empty Magnet URIs",
		FirstMagnetURI:  MagnetURI{},
		SecondMagnetURI: MagnetURI{},
		ExpectedResult:  true,
	},
	{
		Name: "Magnet URIs with all the parameters",
		FirstMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"xt", 0, "xt1"},
				Parameter{"xt", 0, "xt2"},
				Parameter{"dn", 0, "dn1"},
				Parameter{"dn", 0, "dn2"},
				Parameter{"kt", 0, "kt1"},
				Parameter{"kt", 0, "kt2"},
				Parameter{"mt", 0, "mt1"},
				Parameter{"mt", 0, "mt2"},
			},
		},
		SecondMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"xt", 0, "xt1"},
				Parameter{"xt", 0, "xt2"},
				Parameter{"dn", 0, "dn1"},
				Parameter{"dn", 0, "dn2"},
				Parameter{"kt", 0, "kt1"},
				Parameter{"kt", 0, "kt2"},
				Parameter{"mt", 0, "mt1"},
				Parameter{"mt", 0, "mt2"},
			},
		},
		ExpectedResult: true,
	},
	{
		Name: "Magnet URIs with wrong exact topics",
		FirstMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"xt", 0, "xt1"},
				Parameter{"xt", 0, "xt2"},
			},
		},
		SecondMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"xt", 0, "xt1"},
				Parameter{"xt", 0, "wrong parameter"},
			},
		},
		ExpectedResult: false,
	},
	{
		Name: "Magnet URIs with wrong display names",
		FirstMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"dn", 0, "dn1"},
				Parameter{"dn", 0, "dn2"},
			},
		},
		SecondMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"dn", 0, "dn1"},
				Parameter{"dn", 0, "wrong parameter"},
			},
		},
		ExpectedResult: false,
	},
	{
		Name: "Magnet URIs with wrong keyword topics",
		FirstMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"kt", 0, "kt1"},
				Parameter{"kt", 0, "kt2"},
			},
		},
		SecondMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"kt", 0, "kt1"},
				Parameter{"kt", 0, "wrong parameter"},
			},
		},
		ExpectedResult: false,
	},
	{
		Name: "Magnet URIs with wrong manifest topics",
		FirstMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"mt", 0, "mt1"},
				Parameter{"mt", 0, "mt2"},
			},
		},
		SecondMagnetURI: MagnetURI{
			Parameters: []Parameter{
				Parameter{"mt", 0, "mt1"},
				Parameter{"mt", 0, "wrong parameter"},
			},
		},
		ExpectedResult: false,
	},
}

func TestParseMagnetURIWithErrors(t *testing.T) {
	scenarios := parseMagnetURIWithErrorsScenarios
	for _, scenario := range scenarios {
		magnetURI, error := Parse(scenario.RawMagnetURI)
		if !magnetURI.Equal(MagnetURI{}) {
			t.Errorf(
				"Error on test %q: a non-empty Magnet URI was returned: %v.",
				scenario.Name, magnetURI)
		}
		if error == nil {
			t.Error("No error was returned on %q test.", scenario.Name)
		}
		if error.Error() != scenario.ExpectedError {
			t.Errorf(
				"Error on test %q: Expected error message: %q; got %q",
				scenario.Name, scenario.ExpectedError, error.Error())
		}
	}
}

type parseMagnetURIWithErrorsScenario struct {
	Name          string
	RawMagnetURI  string
	ExpectedError string
}

var parseMagnetURIWithErrorsScenarios = []parseMagnetURIWithErrorsScenario{
	{
		Name:         "URI without magnet schema prefix",
		RawMagnetURI: "I don't start with the magnet schema prefix.",
		ExpectedError: "The string doesn't start with the Magnet URI schema " +
			"prefix \"magnet:?\"",
	},
	{
		Name:          "URI without parameter prefix",
		RawMagnetURI:  "magnet:?parameterwithoutprefix",
		ExpectedError: "Parameter without prefix: \"parameterwithoutprefix\"",
	},
	{
		Name:          "URI with unknown parameter prefix",
		RawMagnetURI:  "magnet:?unknown=value",
		ExpectedError: "Unknown parameter prefix: \"unknown\"",
	},
}

func TestParseMagnetURI(t *testing.T) {
	scenarios := magnetURIConvertionScenarios
	for _, scenario := range scenarios {
		magnetURI, error := Parse(scenario.RawMagnetURI)
		if error != nil {
			t.Errorf("There was an error on test %q: %q",
				scenario.Name, error.Error())
		}
		if !magnetURI.Equal(scenario.URIStruct) {
			t.Errorf("Error on test %q: expected Magnet URI: %v; got %v",
				scenario.Name, scenario.URIStruct, magnetURI)
		}
	}
}

type magnetURIConvertionScenario struct {
	Name         string
	URIStruct    MagnetURI
	RawMagnetURI string
}

var magnetURIConvertionScenarios = []magnetURIConvertionScenario{
	// Overview examples taken from
	// http://magnet-uri.sourceforge.net/magnet-draft-overview.txt
	{
		Name: "Overview example 1",
		URIStruct: MagnetURI{
			Parameters: []Parameter{
				Parameter{
					"xt", 0, "urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C",
				},
			},
		},
		RawMagnetURI: "magnet:?xt=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C",
	},
	{
		Name: "Overview example 2",
		URIStruct: MagnetURI{
			Parameters: []Parameter{
				Parameter{
					"xt", 0, "urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C",
				},
				Parameter{
					"dn", 0, "Great+Speeches+-+Martin+Luther+King+Jr.+-+" +
						"I+Have+A+Dream.mp3",
				},
			},
		},
		RawMagnetURI: "magnet:?" +
			"xt=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C&" +
			"dn=Great+Speeches+-+Martin+Luther+King+Jr.+-+" +
			"I+Have+A+Dream.mp3",
	},
	{
		Name: "Overview example 3",
		URIStruct: MagnetURI{
			Parameters: []Parameter{
				Parameter{"kt", 0, "martin+luther+king+mp3"},
			},
		},
		RawMagnetURI: "magnet:?kt=martin+luther+king+mp3",
	},
	{
		Name: "Overview example 4",
		URIStruct: MagnetURI{
			Parameters: []Parameter{
				Parameter{
					"xt", 1, "urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C",
				},
				Parameter{
					"xt", 2, "urn:sha1:TXGCZQTH26NL6OUQAJJPFALHG2LTGBC7",
				},
			},
		},
		RawMagnetURI: "magnet:?" +
			"xt.1=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C&" +
			"xt.2=urn:sha1:TXGCZQTH26NL6OUQAJJPFALHG2LTGBC7",
	},
	{
		Name: "Overview example 5",
		URIStruct: MagnetURI{
			Parameters: []Parameter{
				Parameter{"mt", 0, "http://weblog.foo/all-my-favorites.rss"},
			},
		},
		RawMagnetURI: "magnet:?mt=http://weblog.foo/all-my-favorites.rss",
	},
}

func TestMagnetURIToStringWithoutParameters(t *testing.T) {
	magnetURI := MagnetURI{}
	magnetURIString, error := magnetURI.String()
	expectedErrorMessage := "The Magnet URI has no parameters."
	if magnetURIString != "" {
		t.Errorf("A Magnet URI string was returned: %q.", magnetURIString)
	}
	if error == nil {
		t.Error("No error was returned.")
	}
	if error.Error() != expectedErrorMessage {
		t.Errorf(
			"Expected error message: %q; got %q",
			expectedErrorMessage, error.Error())
	}
}

func TestMagnetURIToString(t *testing.T) {
	scenarios := magnetURIConvertionScenarios
	for _, scenario := range scenarios {
		magnetURIString, error := scenario.URIStruct.String()
		if error != nil {
			t.Errorf("There was an error on test %q: %q",
				scenario.Name, error.Error())
		}
		if magnetURIString != scenario.RawMagnetURI {
			t.Errorf("Error on test %q: expected Magnet URI: %q; got %q",
				scenario.Name, scenario.RawMagnetURI, magnetURIString)
		}
	}
}
