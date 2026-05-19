# Verificar instalación
go version
# Comandos esenciales

* go mod init <nombre-proyecto> # Inicializar módulo
* go run <archivo.go> # Compilar y ejecutar
* go build # Compilar a binario
* go fmt # Formatear código
* go vet # Análisis estático
* go test # Ejecutar tests



## GOROUTINES


## CANALES

## CANALES: PATRONES DE USO (Productor-Consumidor)

* El productor (una funcion) solo envia.
* El consumidor (otra funcion) solo recibe (puede iterar hasta cerrarse).
### Instruccion Select
* Multiplexacion de operaciones de canales.
* Espera multiples canales simultaneamente. 