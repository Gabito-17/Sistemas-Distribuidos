package coap

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
	"sync"

	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"
)

// ServidorCoAP expone recursos REST sobre UDP.
type ServidorCoAP struct {
	sim    *sensor.Simulador
	config nodo.Configuracion
	mu     sync.RWMutex
	modo   string
}

// NuevoServidor crea la instancia del servidor CoAP.
func NuevoServidor(sim *sensor.Simulador, config nodo.Configuracion) *ServidorCoAP {
	return &ServidorCoAP{
		sim:    sim,
		config: config,
		modo:   "automatico",
	}
}

// TODO 6: Iniciar arranca el servidor UDP en el puerto 5683.
func (s *ServidorCoAP) Iniciar() {
	// 6a. Crear router con mux.NewRouter()
	router := mux.NewRouter()

	// 6b. Registrar handler GET /temperatura
	router.Handle("/temperatura", mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		if r.Code() != codes.GET {
			w.SetResponse(codes.MethodNotAllowed, message.TextPlain, bytes.NewReader([]byte("Método no permitido")))
			return
		}

		type RespuestaTemperatura struct {
			NodoID      string  `json:"nodo_id"`
			Temperatura float64 `json:"temperatura"`
			Unidad      string  `json:"unidad"`
			Timestamp   string  `json:"timestamp"`
		}

		msg := RespuestaTemperatura{
			NodoID:      s.config.ID,
			Temperatura: s.sim.ObtenerUltima(),
			Unidad:      "C",
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("Error interno serializando JSON")))
			return
		}

		w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(payload))
	}))

	// 6c. Registrar handler PUT /config
	router.Handle("/config", mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		switch r.Code() {
		case codes.PUT:
			// Leer body del mensaje UDP CoAP
			
			payloadBytes, err := r.Message.ReadBody()
			if err != nil {
				w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("Error leyendo cuerpo del mensaje")))
				return
			}

			// Estructura para actualizar s.modo u otros parámetros de entrada
			var reqData map[string]interface{}
			if err := json.Unmarshal(payloadBytes, &reqData); err != nil {
				w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("JSON inválido")))
				return
			}

			s.mu.Lock()
			if nuevoModo, existe := reqData["modo"].(string); existe {
				s.modo = nuevoModo
			}
			s.mu.Unlock()

			w.SetResponse(codes.Changed, message.TextPlain, bytes.NewReader([]byte("Configuración actualizada")))

		case codes.GET:
			// 6d. Registrar handler GET /config que devuelva la configuración actual en JSON
			s.mu.RLock()
			modoActual := s.modo
			s.mu.RUnlock()

			configRespuesta := map[string]interface{}{
				"nodo_id":            s.config.ID,
				"edificio":           s.config.Edificio,
				"aula":               s.config.Aula,
				"modo":               modoActual,
				"intervalo_segundos": s.config.IntervaloSegundos.Seconds(),
			}

			payload, err := json.Marshal(configRespuesta)
			if err != nil {
				w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("Error interno")))
				return
			}

			w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(payload))

		default:
			w.SetResponse(codes.MethodNotAllowed, message.TextPlain, bytes.NewReader([]byte("Método no permitido")))
		}
	}))

	// 6e. Llamar coap.ListenAndServe("udp", ":5683", router)
	log.Fatal(coap.ListenAndServe("udp", ":5683", router))
}
