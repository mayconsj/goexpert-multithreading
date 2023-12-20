package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mayconsj/goexpert-multithreading/internal/dto"
	"github.com/mayconsj/goexpert-multithreading/pkg/formatter"
	"github.com/mayconsj/goexpert-multithreading/pkg/validator"
	"github.com/spf13/viper"
)

const BRASIL_API_CEP_SERVICE_NAME string = "Brasil Api Cep Service"

type apiCepOutput struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type apiCepError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type BrasilApiCepUseCase struct {
	url string
}

func NewBrasilApiCepUseCase() *BrasilApiCepUseCase {
	return &BrasilApiCepUseCase{
		url: viper.GetString("URL_BRASIL_API_CEP"),
	}
}

func (uc *BrasilApiCepUseCase) Execute(rawCep string) (*dto.AddressApi, *dto.AddressApiError) {
	cep := formatter.SanitalizeCep(rawCep)

	if !validator.Cep(cep) {
		return nil, dto.NewAddressApiError(BRASIL_API_CEP_SERVICE_NAME, "CEP is invalid")
	}

	cep = formatter.Cep(cep)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	endpoint := fmt.Sprintf(uc.url + cep)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, dto.NewAddressApiError(BRASIL_API_CEP_SERVICE_NAME, err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, dto.NewAddressApiError(BRASIL_API_CEP_SERVICE_NAME, err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var cepErr apiCepError
		err = json.NewDecoder(res.Body).Decode(&cepErr)
		if err != nil {
			return nil, dto.NewAddressApiError(BRASIL_API_CEP_SERVICE_NAME, err.Error())
		}

		return nil, dto.NewAddressApiError(BRASIL_API_CEP_SERVICE_NAME, cepErr.Message)
	}

	var response apiCepOutput

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, dto.NewAddressApiError(BRASIL_API_CEP_SERVICE_NAME, err.Error())
	}

	address := dto.NewAddressApi(BRASIL_API_CEP_SERVICE_NAME, response.City, response.State)
	return address, nil
}
