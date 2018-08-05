package function

import (
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	Mutate([]byte) (string, error)
	SetNQuads(string, string, string) error
}

type Functions interface {
	CompaniesReadByName(string, string) []storage.Company
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCompanyCanNotBeCreated means that the company can't be added to database
	ErrCompanyCanNotBeCreated = errors.New("company can't be created")

	// ErrCompanyAlreadyExist means that the company is in the database already
	ErrCompanyAlreadyExist = errors.New("company already exist")
)

// ReadCompaniesByName is a method for get all nodes by categories name
func (executor *Executor) CreateCompany(company storage.Company, language string) (storage.Company, error) {

	existsCompanies := executor.Functions.CompaniesReadByName(company.Name, language)

	if len(existsCompanies) > 0 {
		ExecutorLogger.Printf("Company with name: %v exist: %v", company.Name, existsCompanies[0])
		return existsCompanies[0], ErrCompanyAlreadyExist
	}

	company.IsActive = true

	encodedCompany, err := json.Marshal(company)
	if err != nil {
		return company, ErrCompanyCanNotBeCreated
	}

	_, err = executor.Store.Mutate(encodedCompany)
	if err != nil {
		return company, ErrCompanyCanNotBeCreated
	}

	err = executor.Store.SetNQuads(company.ID, "companyName", "\""+company.Name+"\""+"@"+language)
	if err != nil {
		return company, ErrCompanyCanNotBeCreated
	}

	return storage.Company{}, nil
}
