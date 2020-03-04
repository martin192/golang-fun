package main

import "testing"

func TestThatGivenWordStartingWithVowel_WhenTranslateIsCalled_ThenCorrectResultIsReturned(t *testing.T) {
	translator := NewWordTranslator()

	translated, _ := translator.Translate("apple")
	if expected := "gapple"; translated != expected {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", translated, expected)
	}
}

func TestThatGivenWordStartingWithXR_WhenTranslateIsCalled_ThenCorrectResultIsReturned(t *testing.T) {
	translator := NewWordTranslator()

	translated, _ := translator.Translate("xray")
	if expected := "gexray"; translated != expected {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", translated, expected)
	}
}

func TestThatGivenWordStartingWithConsonants_WhenTranslateIsCalled_ThenCorrectResultIsReturned(t *testing.T) {
	translator := NewWordTranslator()

	translated, _ := translator.Translate("chair")
	if expected := "airchogo"; translated != expected {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", translated, expected)
	}
}

func TestThatGivenWordStartingWithConsonantsFollowedByQU_WhenTranslateIsCalled_ThenCorrectResultIsReturned(t *testing.T) {
	translator := NewWordTranslator()

	translated, _ := translator.Translate("square")
	if expected := "aresquogo"; translated != expected {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", translated, expected)
	}
}

func TestThatGivenSentence_WhenTranslateIsCalled_ThenCorrectResultIsReturned(t *testing.T) {
	translator := NewSentenceTranslator()

	translated, _ := translator.Translate("square, apple, xray, chair!")
	if expected := "aresquogo, gapple, gexray, airchogo!"; translated != expected {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", translated, expected)
	}
}
