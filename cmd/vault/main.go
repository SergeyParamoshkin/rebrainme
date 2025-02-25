package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
)

func main() {
	// Создаем клиент Vault
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Error creating Vault client: %s", err)
	}

	// Устанавливаем токен для аутентификации
	client.SetToken("s.YourVaultTokenHere")

	// Читаем секрет
	secret, err := client.Logical().Read("secret/data/myapp")
	if err != nil {
		log.Fatalf("Error reading secret: %s", err)
	}

	if secret == nil || secret.Data == nil {
		log.Fatalf("Secret not found")
	}

	// Извлекаем данные
	data := secret.Data["data"].(map[string]interface{})
	username := data["username"].(string)
	password := data["password"].(string)

	fmt.Printf("Username: %s, Password: %s\n", username, password)
}
