package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language        string
	CityName        string
	DatabaseGateway string
}

type ErrorResponse struct {
	Error string
	Data  ErrorData
}

type ErrorData struct {
	Error   string
	Request string
}

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal request error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: request.DatabaseGateway}}
	cities, err := executor.ReadCitiesByName(request.CityName, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadCitiesByName error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "ReadCitiesByName error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCities, err := json.Marshal(cities)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal cities error: %v. Error: %v", cities, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal cities error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedCities)
}
