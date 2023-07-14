package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(w http.ResponseWriter, httpCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{}

	if message == "" {
		response = Response{
			Status:  "success",
			Message: "success",
			Data:    data,
		}
	} else {
		response = Response{
			Status:  "success",
			Message: message,
			Data:    data,
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}

}

func ResponseError(w http.ResponseWriter, httpCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	response := Response{
		Status:  "Error",
		Message: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}
