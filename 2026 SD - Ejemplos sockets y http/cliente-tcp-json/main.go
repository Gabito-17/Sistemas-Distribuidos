package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
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
	conexion, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("Error al conectar: %v", err)
	}
	defer conexion.Close()

	decodificador := json.NewDecoder(conexion)
	codificador := json.NewEncoder(conexion)

	// Goroutine para escuchar mensajes del servidor
	go func() {
		for {
			var mensaje Mensaje
			if err := decodificador.Decode(&mensaje); err != nil {
				log.Println("Desconectado del servidor")
				return
			}
			fmt.Printf("[%s] %s: %s\n", mensaje.Timestamp.Format("15:04:05"), mensaje.Emisor, mensaje.Contenido)
		}
	}()

	// Enviar identificación
	identificacion := Mensaje{
		Emisor:    "ClienteJSON",
		Contenido: "identificación",
		Tipo:      "sistema",
		Timestamp: time.Now(),
	}
	codificador.Encode(identificacion)

	// Leer mensajes del usuario y enviarlos
	lector := bufio.NewReader(os.Stdin)
	fmt.Println("Escribí mensajes y presioná Enter (Ctrl+C para salir):")

	for {
		texto, err := lector.ReadString('\n')
		if err != nil {
			log.Fatalf("Error al leer stdin: %v", err)
		}

		mensaje := Mensaje{
			Emisor:    "ClienteJSON",
			Contenido: texto,
			Tipo:      "broadcast",
			Timestamp: time.Now(),
		}

		if err := codificador.Encode(mensaje); err != nil {
			log.Fatalf("Error al enviar: %v", err)
		}
	}
}
