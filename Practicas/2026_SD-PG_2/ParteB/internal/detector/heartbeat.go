package detector

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"sd-comunicacion/pkg/protocolo"
)

// Enviador se encarga de enviar heartbeats UDP periodicamente
type Enviador struct {
	destino   string
	intervalo time.Duration // TODO: usar time.Duration en vez de int64
	nodoID    string
	contador  int
}

// TODO 5: Implementar la funcion NuevaEnviador.
// Debe recibir destino (string), intervalo (time.Duration) y nodoID (string).
func NuevoEnviador(destino string, intervalo time.Duration, nodoID string) *Enviador {
	return &Enviador{
		destino:   destino,
		intervalo: intervalo,
		nodoID:    nodoID,
		contador:  0,
	}
}

// TODO 6: Implementar el metodo (e *Enviador) Iniciar().
// Debe enviar Heartbeat cada 'intervalo' por UDP al destino configurado.
func (e *Enviador) Iniciar() {
	raddr, err := net.ResolveUDPAddr("udp", e.destino)
	if err != nil {
		fmt.Printf("[Enviador] Error al resolver destino %s: %v\n", e.destino, err)
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Printf("[Enviador] Error al abrir socket UDP: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("[Enviador] Iniciado. Enviando heartbeats a %s cada %v...\n", e.destino, e.intervalo)

	ticker := time.NewTicker(e.intervalo)
	defer ticker.Stop()

	for range ticker.C {
		e.contador++
		hb := protocolo.Heartbeat{
			NodoID:    e.nodoID,
			Contador:  e.contador,
			Timestamp: time.Now().Unix(), // <--- SOLUCIONADO: Usa el campo correcto de tu protocolo
		}

		payload, err := json.Marshal(hb)
		if err != nil {
			fmt.Printf("[Enviador] Error al serializar JSON: %v\n", err)
			continue
		}

		_, err = conn.Write(payload)
		if err != nil {
			fmt.Printf("[Enviador] Error al enviar heartbeat: %v\n", err)
		}
	}
}

// Receptor escucha heartbeats y detecta si dejan de llegar.
type Receptor struct {
	puerto  string
	timeout time.Duration // TODO: usar time.Duration en vez de int64
	ultimo  time.Time
	estado  string
	activo  bool
	mu      sync.Mutex
}

// TODO 7: Implementar la funcion NuevoReceptor.
func NuevoReceptor(puerto string, timeout time.Duration) *Receptor {
	return &Receptor{
		puerto:  puerto,
		timeout: timeout,
		estado:  "suspect",
		activo:  false,
	}
}

// TODO 8: Implementar el metodo (r *Receptor) Escuchar().
func (r *Receptor) Escuchar() {
	laddr, err := net.ResolveUDPAddr("udp", r.puerto)
	if err != nil {
		fmt.Printf("[Receptor] Error al resolver puerto %s: %v\n", r.puerto, err)
		return
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Printf("[Receptor] Error al escuchar en puerto %s: %v\n", r.puerto, err)
		return
	}
	defer conn.Close()

	r.activo = true
	r.ultimo = time.Now()

	fmt.Printf("[Receptor] Escuchando heartbeats UDP en %s (Timeout: %v)...\n", r.puerto, r.timeout)

	go r.monitorearTimeout()

	buffer := make([]byte, 1024)
	for r.activo {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if !r.activo {
				break
			}
			fmt.Printf("[Receptor] Error de lectura: %v\n", err)
			continue
		}

		var hb protocolo.Heartbeat
		err = json.Unmarshal(buffer[:n], &hb)
		if err != nil {
			fmt.Printf("[Receptor] Error al decodificar JSON: %v\n", err)
			continue
		}

		r.mu.Lock()
		r.ultimo = time.Now()
		if r.estado != "alive" {
			fmt.Printf("[Receptor] >>> Estado cambio: %s -> ALIVE (Nodo: %s, HB #%d)\n", 
				r.estado, hb.NodoID, hb.Contador)
			r.estado = "alive"
		}
		r.mu.Unlock()
	}
}

func (r *Receptor) monitorearTimeout() {
	ticker := time.NewTicker(r.timeout / 3)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		if !r.activo {
			r.mu.Unlock()
			return
		}

		desdeUltimo := time.Since(r.ultimo)
		estadoAnterior := r.estado

		if desdeUltimo > 2*r.timeout {
			r.estado = "dead"
		} else if desdeUltimo > r.timeout {
			r.estado = "suspect"
		}

		if r.estado != estadoAnterior {
			fmt.Printf("[Receptor] !!! Alerta de tiempo: %s -> %s (Sin reportes desde hace %v)\n", 
				estadoAnterior, r.estado, desdeUltimo.Round(time.Millisecond))
		}
		r.mu.Unlock()
	}
}