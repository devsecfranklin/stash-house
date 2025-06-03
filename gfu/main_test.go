package main

import (
	"strings"
	"testing"
)

var buildInfoTextTestCases = []struct {
	testInput      string
	expectedOutput string
}{
	{
		//Plain line of text is converted to info text
		"A line of text",
		"iA line of text	false	null.host	1",
	},
}

var processFileTestCases = []struct {
	gophermap string
	plaintext string
}{
	{
		//Plain line of text is converted to info text
		"iA line of text	false	null.host	1\n",
		"A line of text\n",
	},

	{
		//Gopher submenu is not converted
		"1Floodgap Home	/home	gopher.floodgap.com	70\n",
		"1Floodgap Home	/home	gopher.floodgap.com	70\n",
	},
	{
		//HTML file is not converted
		"hhttp://tildegit.org/sloum/bombadillo	url:http://tildegit.org/sloum/bombadillo	colorfield.space	70\n",
		"hhttp://tildegit.org/sloum/bombadillo	url:http://tildegit.org/sloum/bombadillo	colorfield.space	70\n",
	},
}

var carriagereturnTestCases = []struct {
	testInput      string
	expectedOutput string
	deconstruct    bool
}{
	{
		//Plaintext to gophermap
		"A test line with a carriage return\r\n",
		"iA test line with a carriage return	false	null.host	1\n",
		false,
	},
	{
		//Gophermap to plaintext
		"iA test line with a cr	false	null.host	1\r\n",
		"A test line with a cr\n",
		true,
	},
	{
		//HTML file
		"hhttp://tildegit.org/sloum/bombadillo	url:http://tildegit.org/sloum/bombadillo	colorfield.space	70\r\n",
		"hhttp://tildegit.org/sloum/bombadillo	url:http://tildegit.org/sloum/bombadillo	colorfield.space	70\n",
		false,
	},
}

func TestBuildInfoText(t *testing.T) {
	for testNumber, testCase := range buildInfoTextTestCases {
		testOutput := buildInfoText(testCase.testInput)
		if testCase.expectedOutput != testOutput {
			t.Errorf(`buildInfoText test case %d failed. Expected "%s", got "%s"`,
				testNumber,
				testCase.expectedOutput,
				testOutput,
			)
			break
		}
	}
}

func TestProcessFileBuild(t *testing.T) {
	for testNum, testCase := range processFileTestCases {
		deconstructInfoTextLines = false
		testOutput, err := processFile(strings.NewReader(testCase.plaintext))
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		if testCase.gophermap != testOutput.String() {
			t.Errorf(`processFile build test case %d failed. Expected "%s", got "%s"`,
				testNum,
				testCase.plaintext,
				testOutput.String(),
			)
			break
		}
	}
}

func TestProcessFileDecon(t *testing.T) {
	for testNum, testCase := range processFileTestCases {
		deconstructInfoTextLines = true
		testOutput, err := processFile(strings.NewReader(testCase.gophermap))
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		if testCase.plaintext != testOutput.String() {
			t.Errorf(`processFile decon test case %d failed. Expected "%s", got "%s"`,
				testNum,
				testCase.plaintext,
				testOutput.String(),
			)
			break
		}
	}
}

func TestCarriageReturnsAreExpunged(t *testing.T) {
	for testNum, testCase := range carriagereturnTestCases {
		deconstructInfoTextLines = testCase.deconstruct
		testOutput, err := processFile(strings.NewReader(testCase.testInput))
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		if testCase.expectedOutput != testOutput.String() {
			t.Errorf(`CR test case %d failed. Expected "%s", got "%s"`,
				testNum,
				testCase.expectedOutput,
				testOutput.String(),
			)
			break
		}
	}
}
