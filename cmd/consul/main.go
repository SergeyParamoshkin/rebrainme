package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Создаем клиент Consul
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Error creating Consul client: %s", err)
	}

	// Читаем ключ из KV-хранилища
	kv := client.KV()
	pair, _, err := kv.Get("myapp/config", nil)
	if err != nil {
		log.Fatalf("Error reading key: %s", err)
	}

	if pair == nil {
		log.Fatalf("Key not found")
	}

	fmt.Printf("Config: %s\n", pair.Value)
}
