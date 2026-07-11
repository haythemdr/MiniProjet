package services

import "strings"

var stopWord = map[string]bool{
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
	"a":    true,
	"an":   true,
	"d":    true,
	"l":    true,
}

var lemmaWords = map[string]string{

	"chargeurs":   "chargeur",
	"cables":      "cable",
	"coques":      "coque",
	"films":       "film",
	"protections": "protection",
	"adaptateurs": "adaptateur",

	"ordinateurs": "ordinateur",
	"telephones":  "telephone",
	"smartphones": "smartphone",

	"ecrans":   "ecran",
	"claviers": "clavier",
	"souris":   "souris",

	"sacs":    "sac",
	"montres": "montre",

	"chaussures": "chaussure",

	"chemises": "chemise",

	"vestes": "veste",

	"robes": "robe",

	"pantalons": "pantalon",
}

func Lemmatize(word string) string {

	if lemma, ok := lemmaWords[word]; ok {
		return lemma
	}

	return word
}

func ReplaceSynonym(word string) string {

	if synonym, ok := Synonyms[word]; ok {
		return synonym
	}

	return word
}

func ProcessQuery(query string) string {

	query = NormalizeName(query)

	words := strings.Fields(query)

	result := make([]string, 0, len(words))

	for _, word := range words {

		word = Lemmatize(word)

		word = ReplaceSynonym(word)

		if stopWord[word] {
			continue
		}

		result = append(result, word)
	}

	return strings.Join(result, " ")
}
