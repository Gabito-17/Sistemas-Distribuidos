package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	conexion, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("Error al conectar: %v", err)
	}
	defer conexion.Close()

	// Enviar mensaje
	fmt.Fprintf(conexion, "Hola servidor\n")

	// Leer respuesta
	lector := bufio.NewReader(conexion)
	respuesta, err := lector.ReadString('\n')
	if err != nil {
		log.Fatalf("Error al leer respuesta: %v", err)
	}
	fmt.Println("Respuesta:", respuesta)
}
