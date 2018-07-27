package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language        string
	CompanyName     string
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

	executor := Executor{Store: storage.New(request.DatabaseGateway)}
	companies, err := executor.ReadCompaniesByName(request.CompanyName, request.Language, request.DatabaseGateway)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadCompaniesByName error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "ReadCompaniesByName error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCompanies, err := json.Marshal(companies)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal companies error: %v. Error: %v", companies, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal companies error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedCompanies)
}
