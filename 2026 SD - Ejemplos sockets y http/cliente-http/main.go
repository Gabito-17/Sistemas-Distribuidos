package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	respuesta, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Fatalf("Error en la petición: %v", err)
	}
	defer respuesta.Body.Close()

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		log.Fatalf("Error al leer respuesta: %v", err)
	}
	fmt.Println("Respuesta:", string(cuerpo))

	respuesta, err = http.Get("http://localhost:8080/fin")
	if err != nil {
		log.Fatalf("Error en la petición: %v", err)
	}
	defer respuesta.Body.Close()

	cuerpo, err = io.ReadAll(respuesta.Body)
	if err != nil {
		log.Fatalf("Error al leer respuesta: %v", err)
	}
	fmt.Println("Respuesta:", string(cuerpo))

}
