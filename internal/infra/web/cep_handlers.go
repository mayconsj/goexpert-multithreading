package web

import (
	"log"
	"net/http"
	"time"

	"github.com/mayconsj/goexpert-multithreading/internal/dto"
	"github.com/mayconsj/goexpert-multithreading/internal/usecase"
	"github.com/mayconsj/goexpert-multithreading/pkg/formatter"
)

type CepHandler struct {
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

func (h *CepHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cep = r.URL.Query().Get("cep")

	brasilApiCepUseCase := usecase.NewBrasilApiCepUseCase()
	viaCepUseCase := usecase.NewViaCepUseCase()

	addressChannel := make(chan *dto.AddressApi)
	addressErrChannel := make(chan *dto.AddressApiError)

	go func() {
		address, err := brasilApiCepUseCase.Execute(cep)
		if err != nil {
			addressErrChannel <- err
			return
		}
		addressChannel <- address
	}()

	go func() {
		address, err := viaCepUseCase.Execute(cep)
		if err != nil {
			addressErrChannel <- err
			return
		}
		addressChannel <- address
	}()

	select {
	case address := <-addressChannel:
		log.Println(formatter.JSON(address))
		w.Write([]byte(formatter.JSON(address)))

	case err := <-addressErrChannel:
		log.Println(formatter.JSON(err))
		w.Write([]byte(formatter.JSON(err)))

	case <-time.After(time.Second):
		log.Println("timeout exceeded")
		w.Write([]byte("timeout exceeded"))
	}

}
