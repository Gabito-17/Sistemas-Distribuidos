package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sd-broadcast/pkg/protocolo"
	"strings"
)

func main() {
	direccionServidor := os.Getenv("SERVIDOR")
	if direccionServidor == "" {
		direccionServidor = "localhost:4000"
	}

	nombre := os.Getenv("NOMBRE")
	if nombre == "" {
		fmt.Print("Ingrese su nombre: ")
		lector := bufio.NewReader(os.Stdin)
		nombreBytes, _, _ := lector.ReadLine()
		nombre = string(nombreBytes)
	}

	// TODO 20: conectar al servidor usando net.Dial("tcp", direccionServidor) - hecho
	conexion, err := net.Dial("tcp", direccionServidor)
	if err != nil {
		log.Fatalf("Error al conectar al servidor: %v", err)
	}
	defer conexion.Close()

	log.Printf("conectado al servidor %s", direccionServidor)

	// Manejar errores y usar defer conexion.Close() - Hecho

	// TODO 21: enviar mensaje de identificación con protocolo.Codificar - Hecho
	identificadorMsj := protocolo.NuevoMensaje(nombre, "", "identificacion")

	err = protocolo.Codificar(conexion, identificadorMsj)

	if err != nil {
		log.Fatalf("error al enviar mensaje: %v", err)
	}

	// mensaje de tipo "identificacion" con Emisor = nombre

	// TODO 22: iniciar una goroutine que escuche mensajes del servidor en paralelo
	// La goroutine debe usar protocolo.Decodificar en un bucle e imprimir los mensajes recibidos
	// Si hay error, imprimir y retornar (el servidor cerró la conexión)

	go recibirMensajes(conexion)

	// TODO 23: en el hilo principal, leer líneas de stdin y enviar mensajes de tipo "broadcast" - Hecho
	// Usar bufio.NewReader(os.Stdin) y ReadString('\n') - Hecho
	// Para cada línea, crear un Mensaje y enviarlo con protocolo.Codificar
	lector := bufio.NewReader(os.Stdin)

	for {
		texto, err := lector.ReadString('\n')
		if err != nil {
			log.Println("Error leyendo entrada")
			break
		}

		texto = strings.TrimSpace(texto)

		if texto == "" {
			continue
		}

		mensaje := protocolo.NuevoMensaje(
			nombre,
			texto,
			"broadcast",
		)

		err = protocolo.Codificar(conexion, mensaje)
		if err != nil {
			log.Println("Error enviando mensaje")
			break
		}
	}

	log.Println("Cliente finalizado")
}

// recibirMensajes lee continuamente desde la conexión e imprime en consola
func recibirMensajes(conexion net.Conn) {
	// TODO 24: implementar bucle infinito de protocolo.Decodificar - Hecho
	// Imprimir Emisor, Contenido y Timestamp de cada mensaje recibido - Hecho
	// Si Decode retorna error, imprimir "Desconectado del servidor" y retornar - Hecho
	for {

		mensaje, err := protocolo.Decodificar(conexion)

		if err != nil {
			log.Println("Desconectado del servidor")
			os.Exit(0)
		}

		fmt.Printf(
			"[%s] %s: %s\n",
			mensaje.Timestamp.Format("15:04:05"),
			mensaje.Emisor,
			mensaje.Contenido,
		)
	}
}
