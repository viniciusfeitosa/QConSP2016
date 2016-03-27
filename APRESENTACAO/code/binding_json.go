package main

import (
	"encoding/json"
	"log"

	pkgUser "gitlab.globoi.com/globoid/glive_models/user"
)

func enviaJSON() ([]byte, error) {
	user := pkgUser.NewUser()
	payload := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}
	output, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error to marshal json error: %s\n", err.Error())
		return []byte{}, err
	}
	return output, nil
}
