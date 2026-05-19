package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// Mensaje representa el protocolo JSON entre cliente y servidor
type Mensaje struct {
	Emisor    string    `json:"emisor"`
	Contenido string    `json:"contenido"`
	Tipo      string    `json:"tipo"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	escuchador, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("Error al crear escuchador: %v", err)
	}
	defer escuchador.Close()

	fmt.Println("Servidor TCP con JSON escuchando en :4000")

	for {
		conexion, err := escuchador.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}
		go manejarCliente(conexion)
	}
}

func manejarCliente(conexion net.Conn) {
	defer conexion.Close()

	decodificador := json.NewDecoder(conexion)
	codificador := json.NewEncoder(conexion)

	for {
		var mensaje Mensaje
		if err := decodificador.Decode(&mensaje); err != nil {
			log.Printf("Cliente desconectado o error: %v", err)
			return
		}

		log.Printf("Recibido de %s: %s", mensaje.Emisor, mensaje.Contenido)

		respuesta := Mensaje{
			Emisor:    "Servidor",
			Contenido: "Recibido: " + mensaje.Contenido,
			Tipo:      "respuesta",
			Timestamp: time.Now(),
		}

		if err := codificador.Encode(respuesta); err != nil {
			log.Printf("Error al enviar respuesta: %v", err)
			return
		}
	}
}
