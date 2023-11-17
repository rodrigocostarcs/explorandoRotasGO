package cidades

import (
	"errors"
	"fmt"
)

func BuscarCidades() map[int]string {
	Cidades := map[int]string{34: "VILHENA-RO", 27: "CACOAL-RO", 3: "JI-PARANÁ-RO", 6: "PRESIDENTE MEDICI-RO", 5: "PORTO VELHO-RO"}
	return Cidades
}

func ValidaDadosHelp(origem string, destino string, data string) error {

	switch origem {
	case "5", "6", "3", "27", "34":
		return nil
	default:
		return errors.New(fmt.Sprintf("O id de origem '%s' não está entre os valores especificados", origem)) // Retorna um erro indicando que o número não está entre os valores
	}

	switch destino {
	case "5", "6", "3", "27", "34":
		return nil
	default:
		return errors.New(fmt.Sprintf("O id de destino '%s' não está entre os valores especificados", destino)) // Retorna um erro indicando que o número não está entre os valores
	}

	if data != "" {
		return errors.New("A data não pode ser vazia")
	}

	return nil
}
