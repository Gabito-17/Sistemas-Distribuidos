package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(escritor http.ResponseWriter, peticion *http.Request) {
		fmt.Fprintf(escritor, "Hola desde HTTP (sobre TCP)")
	})

	http.HandleFunc("/fin", func(escritor http.ResponseWriter, peticion *http.Request) {
		fmt.Fprintf(escritor, "FIN desde HTTP (sobre TCP)")
	})
	
	fmt.Println("Servidor HTTP escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
