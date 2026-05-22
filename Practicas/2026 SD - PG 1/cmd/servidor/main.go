package main

import (
	"log"
	"net"
	"os"

	"sd-broadcast/internal/registro"
	"sd-broadcast/pkg/protocolo"
)

const puertoPorDefecto = "4000"

func main() {
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = puertoPorDefecto
	}

	escuchador, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		log.Fatalf("No se pudo iniciar el escuchador: %v", err)
	}
	defer escuchador.Close()

	log.Printf("Servidor de broadcast escuchando en :%s", puerto)

	// TODO 8: crear un RegistroClientes usando registro.NuevoRegistro() - Hecho

	registroClientes := registro.NuevoRegistro()

	// TODO 9: iniciar goroutine para descubrimiento UDP (bonus) - hecho
	go iniciarDescubrimientoUDP(puerto)

	for {
		conexion, err := escuchador.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}

		// TODO 10: en lugar de llamar directamente a manejarCliente,
		// lanzar una goroutine para atender la conexión concurrentemente - Hecho
		go manejarCliente(conexion, registroClientes)
	}
}

// Utiliza UDP para que los clientes "descubran" si hay servidor
func iniciarDescubrimientoUDP(puerto string) {
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		return
	}
	//conn es un socket UDP (recibe y envia paquetes UDP)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}

	//Cerrar conexion
	defer conn.Close()

	//crea un buffer de 1024 byes (UDP entrega bytes ¡NO STRINGS!)
	buffer := make([]byte, 1024)

	//blucle infinito ¿for{ }?
	for {
		//n es la cantidad de bytes ocupados en el buffer
		n, dirCliente, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		//:n se encarga de tomar solo los bytes ocupados del buffer
		mensaje := string(buffer[:n])

		//Si el mensaje recibido es DESCUBRIR_SERVIDOR
		//Envia SERVIDOR_DISPONIBLE + el puerto (Convertido a byte)
		if mensaje == "DESCUBRIR_SERVIDOR" {

			respuesta := "SEVIDOR_DISPONIBLE:" + puerto

			//dirCliente es la direccion del cliente que mando el paquete(ip:puerto)
			conn.WriteToUDP([]byte(respuesta), dirCliente)
		}
	}
}

func manejarCliente(conexion net.Conn, registroClientes *registro.RegistroClientes) {
	defer conexion.Close()

	// TODO 11: leer el primer mensaje de identificación del cliente
	// Usar protocolo.Decodificar para obtener el nombre del emisor - Hecho
	mensaje, err := protocolo.Decodificar(conexion)
	if err != nil {
		log.Printf("Error al decodificar el mensaje", err)
		return
	}

	nombreCliente := mensaje.Emisor

	log.Printf("Cliente conectado: %s desde %s", nombreCliente, conexion.RemoteAddr())

	// TODO 12: agregar el cliente al registro usando registroClientes.Agregar(nombreCliente, conexion) - Hecho

	registroClientes.Agregar(nombreCliente, conexion)

	// TODO 13: notificar a todos los demás clientes que "nombreCliente se unió"
	// Usar difundirMensaje excepto al emisor - Hecho

	difundirMensaje(registroClientes, protocolo.NuevoMensaje("Sistema", nombreCliente+" se conectó", "Sistema"), nombreCliente)

	// TODO 14: defer para eliminar al cliente del registro al desconectar - Hecho

	defer registroClientes.Eliminar(nombreCliente)

	// defer registroClientes.Eliminar(nombreCliente)
	// defer difundirMensaje(registroClientes, protocolo.NuevoMensaje("Sistema", nombreCliente+" se desconectó", "sistema"), nombreCliente)

	// TODO 15: bucle para leer mensajes del cliente y reenviarlos a todos los demás
	// Usar protocolo.Decodificar en un for {} - Hecho
	// Si el mensaje.Tipo es "broadcast", usar difundirMensaje - Hecho
	// Si hay error en Decode, salir del bucle (cliente desconectado) - Hecho

	for {
		mensaje, err := protocolo.Decodificar(conexion)
		if err != nil {
			//salir del bucle
			break
		}
		if mensaje.Tipo == "broadcast" {
			difundirMensaje(registroClientes, mensaje, nombreCliente)
		}
	}
	log.Printf("Cliente desconectado: %s", nombreCliente)
}

// difundirMensaje envía un mensaje a todos los clientes excepto al emisor indicado
func difundirMensaje(registroClientes *registro.RegistroClientes, mensaje protocolo.Mensaje, exceptoEmisor string) {
	// TODO 16: obtener todas las conexiones del registro - Hecho
	clientes := registroClientes.ObtenerConexiones()
	// TODO 17: iterar sobre las conexiones - Hecho
	for nombre, conexion := range clientes {
		// TODO 18: si el emisor de esa conexión no es exceptoEmisor, enviar el mensaje con protocolo.Codificar - Hecho
		if nombre != exceptoEmisor {
			err := protocolo.Codificar(conexion, mensaje)

			// TODO 19: si Codificar retorna error, ignorar (el cliente puede haberse desconectado abruptamente) - Hecho
			if err != nil {
				continue
			}
		}

	}
}
