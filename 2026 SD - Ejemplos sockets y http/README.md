# Ejemplos de Sockets TCP/UDP y HTTP en Go

Cada carpeta contiene un programa autónomo que se ejecuta con `go run main.go`.

## Requisitos

- Go 1.22 o superior

## Ejemplos

### 1. Servidor TCP Secuencial
Atiende un cliente a la vez.

```bash
cd servidor-tcp-secuencial
go run main.go
```

En otra terminal:
```bash
cd cliente-tcp
go run main.go
```

---

### 2. Servidor TCP Concurrente
Una goroutine por cliente. Permite atender múltiples clientes simultáneamente.

```bash
cd servidor-tcp-concurrente
go run main.go
```

En varias terminales:
```bash
cd cliente-tcp
go run main.go
```

---

### 3. Cliente TCP
Se conecta al servidor en `localhost:4000`, envía un mensaje y muestra la respuesta.

```bash
cd cliente-tcp
go run main.go
```

Requiere que el servidor (ejemplo 1 o 2) esté corriendo.

---

### 4. Servidor UDP
Recibe y responde datagramas sin establecer conexión.

```bash
cd servidor-udp
go run main.go
```

---

### 5. Cliente UDP
Envía un datagrama al servidor UDP y espera la respuesta.

```bash
cd cliente-udp
go run main.go
```

Requiere que el servidor UDP (ejemplo 4) esté corriendo.

---

### 6. Servidor TCP con JSON
Servidor concurrente que comunica mediante mensajes estructurados JSON.

```bash
cd servidor-tcp-json
go run main.go
```

---

### 7. Cliente TCP con JSON
Se conecta al servidor JSON, envía mensajes desde `stdin` y recibe respuestas en paralelo.

```bash
cd cliente-tcp-json
go run main.go
```

Requiere que el servidor JSON (ejemplo 6) esté corriendo.

---

### 8. Servidor HTTP
Abstracción sobre TCP usando el paquete `net/http`.

```bash
cd servidor-http
go run main.go
```

Abrir en navegador: `http://localhost:8080`

---

### 9. Cliente HTTP
Realiza una petición GET al servidor HTTP.

```bash
cd cliente-http
go run main.go
```

Requiere que el servidor HTTP (ejemplo 8) esté corriendo.

