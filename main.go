package main

import (
	"encoding/json"
	"net/http"

	"github.com/rodrigocostarcs/explorandoRotasGO/cidades"
)

func main() {
	http.HandleFunc("/listarcidades", pegarCidadesDisponveis)
	http.ListenAndServe(":8080", nil)
}

func pegarCidadesDisponveis(w http.ResponseWriter, r *http.Request) {
	cidades := cidades.BuscarCidades()
	result, err := json.Marshal(cidades)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
