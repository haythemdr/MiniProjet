package services

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	reSpecialChars = regexp.MustCompile(`[^\p{L}\p{N}\s]+`)
	reSpaces       = regexp.MustCompile(`\s+`)
)

var replacements = map[string]string{

	"go":         "gb",
	"gigaoctet":  "gb",
	"gigaoctets": "gb",
	"gigabyte":   "gb",
	"gigabytes":  "gb",

	"to":         "tb",
	"teraoctet":  "tb",
	"teraoctets": "tb",
	"terabyte":   "tb",
	"terabytes":  "tb",

	"usb-c":  "usb c",
	"type-c": "usb c",

	"wi-fi": "wifi",

	"hewlett packard": "hp",
}

var marketingWords = map[string]bool{

	"nouveau":     true,
	"new":         true,
	"promo":       true,
	"promotion":   true,
	"officiel":    true,
	"original":    true,
	"authentique": true,
	"garantie":    true,
	"disponible":  true,
	"best":        true,
	"seller":      true,
}

var stopWords = map[string]bool{

	"de":   true,
	"du":   true,
	"des":  true,
	"la":   true,
	"le":   true,
	"les":  true,
	"pour": true,
	"avec": true,
	"et":   true,
	"the":  true,
	"for":  true,
	"and":  true,
}

func removeAccents(s string) string {

	t := norm.NFD.String(s)

	var b strings.Builder

	for _, r := range t {

		if unicode.Is(unicode.Mn, r) {
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

func NormalizeName(name string) string {

	name = strings.ToLower(name)

	name = removeAccents(name)

	name = strings.NewReplacer(
		"-", " ",
		"_", " ",
		"/", " ",
		"\\", " ",
		".", " ",
		",", " ",
		"(", " ",
		")", " ",
		"[", " ",
		"]", " ",
		"+", " ",
	).Replace(name)

	name = reSpecialChars.ReplaceAllString(name, " ")

	name = reSpaces.ReplaceAllString(name, " ")

	words := strings.Fields(name)

	tokens := make([]string, 0, len(words))
	seen := make(map[string]bool)

	for _, word := range words {

		if replacement, ok := replacements[word]; ok {
			word = replacement
		}

		if marketingWords[word] {
			continue
		}

		if stopWords[word] {
			continue
		}

		if seen[word] {
			continue
		}

		seen[word] = true
		tokens = append(tokens, word)
	}

	return strings.Join(tokens, " ")
}
