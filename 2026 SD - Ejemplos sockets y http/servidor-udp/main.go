package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	direccion, err := net.ResolveUDPAddr("udp", ":4000")
	if err != nil {
		log.Fatalf("Error al resolver dirección: %v", err)
	}

	conexion, err := net.ListenUDP("udp", direccion)
	if err != nil {
		log.Fatalf("Error al crear socket UDP: %v", err)
	}
	defer conexion.Close()

	fmt.Println("Servidor UDP escuchando en :4000")

	bufer := make([]byte, 1024)
	for {
		n, addr, err := conexion.ReadFromUDP(bufer)
		if err != nil {
			log.Printf("Error al leer: %v", err)
			continue
		}
		fmt.Printf("Recibido de %v: %s\n", addr, string(bufer[:n]))
		conexion.WriteToUDP([]byte("Eco UDP"), addr)
	}
}
