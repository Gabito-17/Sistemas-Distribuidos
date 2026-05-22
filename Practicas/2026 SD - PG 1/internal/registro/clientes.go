package registro

import (
	"net"
	"sync"
)

// RegistroClientes mantiene el listado de conexiones activas de forma segura
type RegistroClientes struct {
	// TODO 1: agregar un campo sync.RWMutex para proteger el mapa - Hecho
	clientes map[string]net.Conn
	//Lector-Escritor Mutex (Su función es garantizar que un solo hilo o proceso a la vez pueda acceder a un recurso compartido)
	mutex sync.RWMutex
}

// NuevoRegistro crea un registro vacío
func NuevoRegistro() *RegistroClientes {
	// TODO 2: inicializar el mapa de clientes - Hecho
	// Se crea un nuevo espacio en memoria para la estructura, devuelve un puntero "&" a esa estructura.
	return &RegistroClientes{clientes: make(map[string]net.Conn)}
}

// Agregar añade un cliente al registro
func (r *RegistroClientes) Agregar(nombre string, conexion net.Conn) {
	// TODO 3: bloquear para escritura, agregar al mapa, desbloquear

	//Se bloquea
	r.mutex.Lock()

	//Se desbloquea al finalizar la funcion
	defer r.mutex.Unlock()

	//Se agrega la conezion al mapa
	r.clientes[nombre] = conexion

}

// Eliminar remueve un cliente del registro
func (r *RegistroClientes) Eliminar(nombre string) {
	// TODO 4: bloquear para escritura, eliminar del mapa, desbloquear

	//Bloquear
	r.mutex.Lock()

	//Desbloquear al finalziar la funcion
	defer r.mutex.Unlock()

	//Quitar la conexion del mapa
	delete(r.clientes, nombre)

}

// ObtenerConexiones devuelve una copia de todas las conexiones activas
func (r *RegistroClientes) ObtenerConexiones() map[string]net.Conn {
	// TODO 5: bloquear para lectura, copiar conexiones a un slice, desbloquear

	//Bloquear lectura
	r.mutex.RLock()

	//Desbloquear lectura
	defer r.mutex.RUnlock()

	copia := make(map[string]net.Conn)
	for nombre, conexion := range r.clientes {
		copia[nombre] = conexion
	}

	return copia
}

// Cantidad devuelve el número de clientes conectados
func (r *RegistroClientes) Cantidad() int {
	// TODO 6: bloquear para lectura, retornar len del mapa, desbloquear

	//Bloquear lectura
	r.mutex.RLock()

	//desbloquear lectura
	defer r.mutex.RUnlock()

	//retornar len del mapa
	return len(r.clientes)
}

// Nombres devuelve un slice con los nombres de los clientes
func (r *RegistroClientes) Nombres() []string {
	// TODO 7: bloquear para lectura, copiar nombres a un slice, desbloquear
	//bloquear lectura
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	//Inicializar Slice snom
	snom := make([]string, 0, len(r.clientes))
	//Iterar rellenando el Slice
	for nombre := range r.clientes {
		snom = append(snom, nombre)
	}

	return snom
}
