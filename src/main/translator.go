package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	vowels          = "aeiou"
	forbiddenChars  = "’'"
	punctationMarks = "-:,.!?"
)

type Translator interface {
	Translate(string) (string, error)
}

func isValid(input string) (bool, error) {
	if len(input) > 0 && !strings.ContainsAny(input, forbiddenChars) {
		return true, nil
	}

	log.Println(fmt.Sprintf("Invalid input given. Translation not possible for '%s'", input))
	return false, errors.New(fmt.Sprintf("Invalid input given. Translation not possible for '%s'", input))
}

func normalizeWord(word string) string {
	return strings.Trim(word, punctationMarks)
}

type wordTranslator struct{}

func (translator *wordTranslator) Translate(word string) (string, error) {
	log.Println(fmt.Sprintf("Translating word '%s'...", word))

	if isValid, error := isValid(word); !isValid {
		return "", error
	}

	word = strings.ToLower(word)
	firstVowelAt := strings.IndexAny(word, vowels)

	switch {
	case firstVowelAt == 0:
		return fmt.Sprintf("g%s", word), nil
	case strings.HasPrefix(word, "xr"):
		return fmt.Sprintf("ge%s", word), nil
	case firstVowelAt == -1:
		return fmt.Sprintf("%sogo", word), nil
	case word[firstVowelAt-1] == 'q' && word[firstVowelAt] == 'u':
		return fmt.Sprintf("%s%suogo", word[firstVowelAt+1:], word[0:firstVowelAt]), nil
	default:
		return fmt.Sprintf("%s%sogo", word[firstVowelAt:], word[0:firstVowelAt]), nil
	}
}

type sentenceTranslator struct {
	inner wordTranslator
}

func (translator *sentenceTranslator) Translate(sentence string) (string, error) {
	log.Println(fmt.Sprintf("Translating sentence '%s'...", sentence))

	if isValid, error := isValid(sentence); !isValid {
		return "", error
	}

	var stringBuilder strings.Builder

	for _, word := range strings.Split(sentence, " ") {
		normalized := normalizeWord(word)

		if translated, error := translator.inner.Translate(normalized); error == nil {
			stringBuilder.WriteString(strings.Replace(word, normalized, translated, 1))
			stringBuilder.WriteString(" ")
		} else {
			return "", error
		}
	}

	return strings.Trim(stringBuilder.String(), " "), nil
}

func NewWordTranslator() Translator {
	return &wordTranslator{}
}

func NewSentenceTranslator() Translator {
	return &sentenceTranslator{
		inner: wordTranslator{},
	}
}
