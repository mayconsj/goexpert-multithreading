package usecase_test

import (
	"log"
	"testing"

	"github.com/mayconsj/goexpert-multithreading/internal/usecase"
)

func TestShouldReturnAddressSuccessfulyFromViaCEPApi(t *testing.T) {
	uc := usecase.NewViaCepUseCase()
	address, err := uc.Execute("37540000")

	if err != nil {
		t.Errorf("fail to get address from api: %#v", err)
	}

	log.Printf("Api Name: %s", address.ApiName)
	log.Printf("City: %s", address.City)
	log.Printf("State: %s", address.State)
}

func TestShouldReturnAddressSuccessfulyFromApiCep(t *testing.T) {
	uc := usecase.NewBrasilApiCepUseCase()
	address, err := uc.Execute("37540000")

	if err != nil {
		t.Errorf("fail to get address from api: %#v", err)
	}

	log.Printf("Api Name: %s", address.ApiName)
	log.Printf("City: %s", address.City)
	log.Printf("State: %s", address.State)
}
