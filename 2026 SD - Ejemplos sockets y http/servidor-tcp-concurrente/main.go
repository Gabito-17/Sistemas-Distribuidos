package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	escuchador, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("Error al crear escuchador: %v", err)
	}
	defer escuchador.Close()

	fmt.Println("Servidor TCP concurrente escuchando en :4000")

	for {
		conexion, err := escuchador.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}
		go manejarConexion(conexion)
	}
}

func manejarConexion(conexion net.Conn) {
	defer conexion.Close()
	lector := bufio.NewReader(conexion)
	texto, err := lector.ReadString('\n')
	if err != nil {
		log.Printf("Error al leer: %v", err)
		return
	}
	fmt.Fprintf(conexion, "Eco: %s", texto)
}
