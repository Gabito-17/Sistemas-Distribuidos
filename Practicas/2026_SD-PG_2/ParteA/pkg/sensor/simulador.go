package sensor

import (
	"math/rand"
	"sync"
)

// Simulador genera lecturas de temperatura de forma thread-safe.
type Simulador struct {
	mu            sync.RWMutex
	ultimaLectura float64
}

// NuevoSimulador crea un simulador con una lectura inicial.
func NuevoSimulador() *Simulador {
	return &Simulador{
		ultimaLectura: 22.0 + rand.Float64()*5.0, // entre 22.0 y 27.0
	}
}

// Leer devuelve una nueva lectura simulada y la almacena.
func (s *Simulador) Leer() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: generar una temperatura realista variando +/- 0.5 grados respecto a la última lectura.

	// Variación realista entre -0.5 y +0.5
	variacion := (rand.Float64() * 1.0) - 0.5
	nuevaLectura := s.ultimaLectura + variacion

	// Usar rand.Float64() y mantener la nueva lectura dentro de un rango razonable (15°C - 35°C).
	if nuevaLectura < 15.0 {
		nuevaLectura = 15.0
	} else if nuevaLectura > 35.0 {
		nuevaLectura = 35.0
	}
	
	s.ultimaLectura = nuevaLectura
	return s.ultimaLectura
}

// ObtenerUltima devuelve la última lectura sin generar una nueva.
func (s *Simulador) ObtenerUltima() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ultimaLectura
}
