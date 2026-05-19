package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Resolvemos la dirección del servidor UDP
	direccion, err := net.ResolveUDPAddr("udp", "localhost:4000")
	if err != nil {
		log.Fatalf("Error al resolver dirección: %v", err)
	}

	// Creamos un socket UDP local (el SO asigna puerto automáticamente)
	conexion, err := net.DialUDP("udp", nil, direccion)
	if err != nil {
		log.Fatalf("Error al conectar UDP: %v", err)
	}
	defer conexion.Close()

	// Enviamos un datagrama
	mensaje := []byte("Hola desde cliente UDP")
	_, err = conexion.Write(mensaje)
	if err != nil {
		log.Fatalf("Error al enviar: %v", err)
	}
	fmt.Println("Enviado:", string(mensaje))

	// Esperamos respuesta con timeout
	bufer := make([]byte, 1024)
	conexion.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conexion.Read(bufer)
	if err != nil {
		log.Fatalf("Error al recibir: %v", err)
	}
	fmt.Println("Recibido:", string(bufer[:n]))
}
