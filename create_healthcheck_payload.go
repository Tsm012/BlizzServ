package main

type CreateHealthcheckPayload struct {
	Endpoint string `json:"endpoint"`
}

type CreateHealthcheckResponse struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}
