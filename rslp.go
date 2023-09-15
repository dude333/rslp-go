package rslp

import (
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Palavra(term string) string {
	term = strings.ToLower(term)

	for _, rule := range defaultDictionary {
		if rule.MinWordSize > 0 && len(term) < rule.MinWordSize {
			continue
		}

		stem := applyRules(term, rule)

		if term != stem && (rule.Name == Substantivo || rule.Name == Verbo) {
			return removeAccents(stem)
		}

		term = stem
	}

	return removeAccents(term)
}

func applyRules(term string, rule Rule) string {
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
				return strings.TrimSuffix(term, subRule.Suffix) + subRule.Replace
			}
		}
	}
	return term
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

var stopwords = []string{
	"a",
	"as",
	"com",
	"da",
	"das",
	"de",
	"do",
	"dos",
	"e",
	"em",
	"na",
	"nas",
	"no",
	"nos",
	"não",
	"o",
	"os",
	"para",
	"que",
	"se",
	"um",
	"uma",
	"-",
	"_",
}

func Frase(doc string) string {
	s := strings.Map(removePunctuation, doc)
	words := strings.Fields(s)
	result := ""
	for _, word := range words {
		if slices.Contains(stopwords, word) {
			continue
		}
		result += Palavra(word) + " "
	}
	return strings.TrimSpace(result)
}
