package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language         string
	DatabaseGateway  string
	FunctionsGateway string
	City             storage.City
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

	executor := Executor{
		Store: &storage.Store{DatabaseGateway: request.DatabaseGateway},
		Functions: &FAASFunctions{
			DatabaseGateway:  request.DatabaseGateway,
			FunctionsGateway: request.FunctionsGateway}}

	createdCity, err := executor.CreateCity(request.City, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"CreateCity error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "CreateCity error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCity, err := json.Marshal(createdCity)
	if err != nil {
		warning := fmt.Sprintf(
			"Marshal city error: %v. Error: %v", createdCity, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Marshal createdCity error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedCity)
}
