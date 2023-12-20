package main

import (
	"fmt"

	"github.com/mayconsj/goexpert-multithreading/configs"
	"github.com/mayconsj/goexpert-multithreading/internal/infra/web"
	"github.com/mayconsj/goexpert-multithreading/internal/infra/web/webserver"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	cepHandler := web.NewCepHandler()

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webserver.AddHandler("/", cepHandler.Create)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start()

}
