package services

import (
	"encoding/json"
	"log"
	"os"
)

var Synonyms map[string]string

func LoadSynonyms(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read synonyms file: %v", err)
	}

	if err := json.Unmarshal(data, &Synonyms); err != nil {
		log.Fatalf("Cannot parse synonyms file: %v", err)
	}

	log.Printf("✅ %d synonyms loaded", len(Synonyms))
}
