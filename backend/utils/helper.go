package utils

import (
	"encoding/json"
)

type ErrorMsg struct {
	Message string `json:"message"`
}

type ResponseMsg struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func AIToAB(arr []interface{}) []byte {
	result := make([]byte, 0, len(arr))
	for _, val := range arr {
		userBytes, err := json.Marshal(val)
		if err != nil {
			// Handle error
		}
		result = append(result, userBytes...)
	}

	return result
}
