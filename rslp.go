package rslp

import (
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RSLPStemmer(term string, only string) string {
	var removed bool
	var stem string
	term = strings.ToLower(term)
	only = strings.ToLower(only)

	for _, rule := range defaultDictionary {
		if rule.MinWordSize > 0 && len(term) < rule.MinWordSize {
			continue
		}

		// fmt.Printf(`--> "%s" != "%s -> %v\n"`, only, rule.Name, only != strings.ToLower(rule.Name))
		if only != "" && only != strings.ToLower(rule.Name) {
			continue
		}

		if (strings.ToLower(rule.Name) == "verb" || strings.ToLower(rule.Name) == "vowel") && removed {
			continue
		}

		if len(rule.Suffixes) == 0 || containsSuffix(term, rule.Suffixes) {
			for _, subRule := range rule.SubRules {
				if rule.WithSuffix {
					if containsSuffix(term, subRule.Exceptions) {
						continue
					}
				} else {
					if containsException(term, subRule.Exceptions) {
						continue
					}
				}

				if len(term)-len(subRule.Suffix) < subRule.StemSize {
					continue
				}

				if strings.HasSuffix(term, subRule.Suffix) {
					removed = true
					term = strings.TrimSuffix(term, subRule.Suffix)
					stem = term + subRule.Replace
					break
				}
			}
		}
	}

	if only != "" {
		return term
	}

	return stem
}

func containsSuffix(term string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(term, suffix) {
			return true
		}
	}
	return false
}

func containsException(term string, exceptions []string) bool {
	for _, exception := range exceptions {
		if term == exception {
			return true
		}
	}
	return false
}

func removePunctuation(r rune) rune {
	if strings.ContainsRune(".,:;", r) {
		return -1
	} else {
		return r
	}
}

func removeAccents(str string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, str)
	return result
}

func rslp(word string) string {
	if strings.HasSuffix(word, "s") {
		word = RSLPStemmer(word, "Plural")
		log.Println("plural: ", word)
	}

	if strings.HasSuffix(word, "a") {
		word = RSLPStemmer(word, "Feminine")
		log.Println("feminina: ", word)
	}

	word = RSLPStemmer(word, "Augmentative")
	log.Println("aumentativa: ", word)
	word = RSLPStemmer(word, "Adverb")
	log.Println("adverbio: ", word)

	wordAfterNoun := RSLPStemmer(word, "Noun")

	if wordAfterNoun == word {
		wordAfterVerb := RSLPStemmer(word, "Verb")
		if wordAfterVerb == word {
			word = RSLPStemmer(word, "Vowel")
		} else {
			word = wordAfterVerb
		}
	} else {
		word = wordAfterNoun
	}

	word = removeAccents(word)
	return word
}
