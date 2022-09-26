package utils

type Response struct {
    Status  int                    `json:"status"`
    Message string                 `json:"message"`
    Data    map[string]interface{} `json:"data"`
}

type RetrieveResponse struct {
	Status  int           `json:"status,omitempty"`
	Message string        `json:"message,omitempty"`
	Data    []interface{} `json:"data,omitempty"`
}

