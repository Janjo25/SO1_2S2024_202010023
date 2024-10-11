package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "¡Hola, Olimpiadas de la USAC! Esta es la facultad de Agronomía 🌾🌿")
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := "8080"

	fmt.Printf("Servidor corriendo en el puerto %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
