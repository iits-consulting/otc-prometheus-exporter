package otcdoc

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseDocumentationPageFromHtmlBytes(t *testing.T) {

	testInput := "testdata/basic_ecs_metrics.html"
	inputBytes, err := os.ReadFile(testInput)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ParseDocumentationPageFromHtmlBytes(inputBytes, "ecs")

	expectedFile := "testdata/basic_ecs_metrics_expected.json"
	expectedBytes, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatal(err)
	}
	var expected DocumentationPage
	err = json.Unmarshal(expectedBytes, &expected)
	if err != nil {
		t.Fatal(err)
	}

	if actual.Namespace != expected.Namespace {
		t.Fatalf(
			"Expected namespace to be %s but got %s",
			expected.Namespace,
			actual.Namespace,
		)
	}

	if len(actual.Metrics) != len(expected.Metrics) {
		t.Fatalf(
			"Expected number of parsed metrics to be %d but got %d",
			len(expected.Metrics),
			len(actual.Metrics),
		)
	}

	for i := 0; i < len(actual.Metrics); i++ {
		if actual.Metrics[i] != expected.Metrics[i] {
			t.Errorf(
				"Metrics at position %d are different\nexpected %#v\ngot      %#v",
				i,
				expected.Metrics[i],
				actual.Metrics[i],
			)
		}
	}
}
