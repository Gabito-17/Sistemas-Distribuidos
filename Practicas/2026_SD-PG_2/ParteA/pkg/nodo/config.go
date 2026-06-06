package nodo

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
	
)

// Configuracion representa los parámetros del nodo IoT.
type Configuracion struct {
	ID                string
	Edificio          string
	Aula              string
	BrokerMQTT        string
	IntervaloSegundos time.Duration
}

// CargarConfiguracion lee variables de entorno o usa valores por defecto.
func CargarConfiguracion() Configuracion {
	id 				:= obtenerEnv("NODO_ID", "nodo-01")
	edificio 		:= obtenerEnv("NODO_EDIFICIO", "ingenieria")
	aula 			:= obtenerEnv("NODO_AULA", "lab3")
	broker 			:= obtenerEnv("MQTT_BROKER", "localhost:1883")
	intervaloStr 	:= obtenerEnv("INTERVALO_SEGUNDOS", "5")

	// TODO: validar que ID, Edificio y Aula no estén vacíos. 
	// Validar que el intervalo sea un número positivo. Si no salir con error.
	// Sugerencia: usar regexp para permitir solo letras, números y guiones.

	// Validación mediante expresiones regulares (sólo letras, números y guiones)
	validarRegexp := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	if id == "" || !validarRegexp.MatchString(id) {
		log.Fatalf("Error: NODO_ID inválido o vacío ('%s'). Solo se permiten letras, números y guiones.", id)
	}

	if edificio == "" || !validarRegexp.MatchString(edificio) {
		log.Fatalf("Error: NODO_EDIFICIO inválido o vacío ('%s'). Solo se permiten letras, números y guiones.", edificio)
	}

	if aula == "" || !validarRegexp.MatchString(aula) {
		log.Fatalf("Error: NODO_AULA inválido o vacío ('%s'). Solo se permiten letras, números y guiones.", aula)
	}

	intervaloInt, err := strconv.Atoi(intervaloStr)
	if err != nil || intervaloInt <= 0 {
		log.Fatalf("Error: INTERVALO_SEGUNDOS debe ser un número entero positivo. Valor recibido: '%s'", intervaloStr)
	}

	duracion := time.Duration(intervaloInt) * time.Second

	// duracion, err := time.ParseDuration(intervalo + "s")
	// if err != nil {
	// 	duracion = 5 * time.Second
	// }

	return Configuracion{
		ID:                id,
		Edificio:          edificio,
		Aula:              aula,
		BrokerMQTT:        broker,
		IntervaloSegundos: duracion,
	}
}

func obtenerEnv(clave, valorPorDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return valorPorDefecto
}
