package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"
)

// Cliente encapsula la conexión MQTT del nodo.
type Cliente struct {
	config   nodo.Configuracion
	interno  mqtt.Client
	opciones *mqtt.ClientOptions
}

// TODO 1: NuevoCliente crea la configuración inicial del cliente MQTT.
// Debe:
//   1a. Construir el tópico del testamento: nodo/{id}/estado
//   1b. Configurar el mensaje del testamento como {"estado":"offline"} con QoS 1 y retained=true.
//   1c. Configurar ClientID único, timeout de conexión y reconexión automática.
//
// Sugerencia: usar mqtt.NewClientOptions().AddBroker(...).SetClientID(...).SetWill(...)
func NuevoCliente(config nodo.Configuracion) (*Cliente, error) {
	// COMPLETAR

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.BrokerMQTT)
	opts.SetClientID(fmt.Sprintf("%s-client", config.ID))
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetAutoReconnect(true)

	// 1a. Construir el tópico del testamento
	topicoLWT := fmt.Sprintf("nodo/%s/estado", config.ID)

	// 1b. Configurar el mensaje del testamento como {"estado":"offline"} con QoS 1 y retained=true.
	payloadLWT := `{"estado":"offline"}`
	opts.SetWill(topicoLWT, payloadLWT, 1, true)

	clienteInterno := mqtt.NewClient(opts)

	return &Cliente{
		config:   config,
		interno:  clienteInterno,
		opciones: opts,
	}, nil
}

// TODO 2: Conectar establece la sesión con el broker.
// Tras conectar, debe publicar un mensaje retenido {"estado":"online"} en nodo/{id}/estado.
func (c *Cliente) Conectar() error {
	// COMPLETAR
	token := c.interno.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Tras conectar, publicar mensaje retenido {"estado":"online"} en nodo/{id}/estado.
	topicoEstado := fmt.Sprintf("nodo/%s/estado", c.config.ID)
	payloadOnline := `{"estado":"online"}`
	tokenPub := c.interno.Publish(topicoEstado, 1, true, payloadOnline)
	tokenPub.Wait()

	return nil
}

// TODO 3: PublicarLecturas envía periódicamente las lecturas del sensor.
// Debe:
//   3a. Construir el tópico: campus/{edificio}/{aula}/sensor/temperatura
//   3b. En un ticker cada config.IntervaloSegundos, llamar sim.Leer(), serializar a JSON y publicar con QoS 1.
//   El JSON debe tener: {"nodo_id": ..., "temperatura": ..., "unidad":"C", "timestamp":"..."}
func (c *Cliente) PublicarLecturas(sim *sensor.Simulador, config nodo.Configuracion) {
	// COMPLETAR

	// 3a. Construir el tópico
	topicoLecturas := fmt.Sprintf("campus/%s/%s/sensor/temperatura", config.Edificio, config.Aula)
	ticker := time.NewTicker(config.IntervaloSegundos)
	defer ticker.Stop()

	type MensajeTemperatura struct {
		NodoID      string  `json:"nodo_id"`
		Temperatura float64 `json:"temperatura"`
		Unidad      string  `json:"unidad"`
		Timestamp   string  `json:"timestamp"`
	}

	for range ticker.C {
		if !c.interno.IsConnected() {
			log.Println("Cliente MQTT desconectado. Esperando reconexión para publicar...")
			continue
		}

		temp := sim.Leer()
		msg := MensajeTemperatura{
			NodoID:      config.ID,
			Temperatura: temp,
			Unidad:      "C",
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error al serializar lectura MQTT: %v", err)
			continue
		}

		token := c.interno.Publish(topicoLecturas, 1, false, payload)
		token.Wait()
		log.Printf("[MQTT] Publicado en %s -> %s", topicoLecturas, string(payload))
	}
}

// TODO 4: SuscribirComandos se une al tópico de actuadores y procesa mensajes.
// Debe:
//   4a. Suscribirse al tópico campus/{edificio}/{aula}/actuador/cmd con QoS 1.
//   4b. En el callback, deserializar el JSON, imprimir el comando recibido y simular la ejecución.
//   Ejemplo de payload esperado: {"accion":"encender_alarma", "origen":"dashboard"}
func (c *Cliente) SuscribirComandos(config nodo.Configuracion) error {
	// COMPLETAR
	topicoCmd := fmt.Sprintf("campus/%s/%s/actuador/cmd", config.Edificio, config.Aula)

	type Comando struct {
		Accion string `json:"accion"`
		Origen string `json:"origen"`
	}

	callback := func(client mqtt.Client, msg mqtt.Message) {
		var cmd Comando
		if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
			log.Printf("Error al deserializar comando recibido: %v. Payload: %s", err, string(msg.Payload()))
			return
		}

		log.Printf("[MQTT] Comando Recibido desde %s: ¡Ejecutando acción '%s'!", cmd.Origen, cmd.Accion)
		
		// Simulación de la ejecución de la acción
		switch cmd.Accion {
		case "encender_alarma":
			log.Println("⚠️ [ACTUADOR] ¡¡ALERTA!! Alarma encendida en el aula.")
		case "apagar_alarma":
			log.Println("✅ [ACTUADOR] Alarma apagada. Estado normalizado.")
		default:
			log.Printf("❓ [ACTUADOR] Acción desconocida: %s", cmd.Accion)
		}
	}

	token := c.interno.Subscribe(topicoCmd, 1, callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("Suscrito con éxito al tópico: %s", topicoCmd)
	return nil
}

// TODO 5: Desconectar cierra limpiamente la sesión MQTT.
// Sugerencia: publicar estado offline retenido antes de desconectar.
func (c *Cliente) Desconectar() {
	// COMPLETAR

	topicoEstado := fmt.Sprintf("nodo/%s/estado", c.config.ID)
	payloadOffline := `{"estado":"offline"}`
	token := c.interno.Publish(topicoEstado, 1, true, payloadOffline)
	token.Wait()

	c.interno.Disconnect(250)
	log.Println("Conexión MQTT cerrada limpiamente.")
}
