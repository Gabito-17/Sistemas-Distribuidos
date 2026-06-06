package telemetria

import (
	"fmt"
	"sync"
	"time"
)

// TODO 1: Definir el struct Telemetria que sera el servicio RPC.
// Debe contener un mapa protegido por sync.Mutex para almacenar
// la ultima lectura de cada sensor.
// Sugerencia: usar map[string]Lectura
//
// import "sync" cuando lo necesites

// Estructuras de datos para el protocolo RPC (Exportadas y con tags JSON)
type Lectura struct {
	SensorID  string    `json:"sensor_id"`
	Valor     float64   `json:"valor"`
	Timestamp time.Time `json:"timestamp"`
}

type RespuestaLectura struct {
	ID        int64 `json:"id"`
	Procesado bool  `json:"procesado"`
}

type ConsultaUltimaLectura struct {
	SensorID string `json:"sensor_id"`
}

type Telemetria struct {
	mu         sync.Mutex
	lecturas   map[string]Lectura
	contadorID int64
}

// NewTelemetria inicializa correctamente el servicio y su mapa interno
func NewTelemetria() *Telemetria {
	return &Telemetria{
		lecturas: make(map[string]Lectura),
	}
}

// TODO 2: Implementar el metodo RPC RegistrarLectura.
// Firma requerida por net/rpc:
//   func (t *Telemetria) RegistrarLectura(args Lectura, resp *RespuestaLectura) error
// Debe:
//   - Guardar la lectura en el mapa (protegiendo con mutex)
//   - Asignar un ID incremental a la respuesta
//   - Loguear la lectura recibida (import "fmt" y "time")
//   - Retornar nil en caso de exito
func (t *Telemetria) RegistrarLectura(args Lectura, resp *RespuestaLectura) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.contadorID++
	t.lecturas[args.SensorID] = args

	fmt.Printf("[%s] Lectura recibida - Sensor: %s, Valor: %.2f\n", 
		time.Now().Format("15:04:05"), args.SensorID, args.Valor)

	resp.ID = t.contadorID
	resp.Procesado = true

	return nil
}

// TODO 3: Implementar el metodo RPC ObtenerUltimaLectura.
// Firma requerida por net/rpc:
//   func (t *Telemetria) ObtenerUltimaLectura(args ConsultaUltimaLectura, resp *Lectura) error
// Debe:
//   - Buscar en el mapa la ultima lectura del SensorID solicitado
//   - Si no existe, retornar un error con fmt.Errorf
//   - Si existe, copiar el valor a resp y retornar nil
func (t *Telemetria) ObtenerUltimaLectura(args ConsultaUltimaLectura, resp *Lectura) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	lectura, existe := t.lecturas[args.SensorID]
	if !existe {
		return fmt.Errorf("no se encontraron lecturas para el sensor con ID: %s", args.SensorID)
	}

	*resp = lectura
	return nil
}