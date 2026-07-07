package elasticsearch

import (
	"bytes"
	"log"
)

func CreateIndex() {

	mapping := `
{
  "settings": {
    "analysis": {
      "analyzer": {
        "autocomplete": {
          "tokenizer": "standard",
          "filter": [
            "lowercase"
          ]
        }
      }
    }
  },
  "mappings": {
    "properties": {

      "name": {
        "type": "text",
        "analyzer": "autocomplete",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },

      "store": {
        "type": "keyword"
      },

      "price": {
        "type": "keyword"
      },

      "url": {
        "type": "keyword"
      },

      "image": {
        "type": "keyword"
      },

      "name_completion": {
        "type": "completion"
      }
    }
  }
}
`

	res, err := Client.Indices.Create(
		"products",
		Client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
	)

	if err != nil {
		log.Println(err)
		return
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Println("Index already exists")
		return
	}

	log.Println("✅ Elasticsearch index created")
}
