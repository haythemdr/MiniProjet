package elasticsearch

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"tunisianet-scraper/internal/models"
)

func IndexProduct(product models.Product) error {

	log.Println("📦 Indexing:", product.Name)

	body, err := json.Marshal(map[string]interface{}{
		"name":  product.Name,
		"price": product.Price,
		"store": product.Store,
		"url":   product.URL,
		"image": product.Image,
		"name_completion": map[string]interface{}{
			"input": []string{product.Name},
		},
	})

	if err != nil {
		return err
	}

	// Use URL as the document ID
	hash := sha1.Sum([]byte(product.URL))
	docID := hex.EncodeToString(hash[:])

	log.Println("➡️ Sending document to Elasticsearch...")

	res, err := Client.Index(
		"products",
		bytes.NewReader(body),
		Client.Index.WithDocumentID(docID),
	)

	if err != nil {
		return fmt.Errorf("client error: %w", err)
	}

	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)

	log.Println("HTTP Status :", res.Status())
	log.Println("Response    :", string(bodyBytes))

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", string(bodyBytes))
	}

	log.Println("✅ Indexed :", product.Name)

	return nil
}
