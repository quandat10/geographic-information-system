package utils

type ErrorMsg struct {
	Message string `json:"message"`
}

type ResponseMsg struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
