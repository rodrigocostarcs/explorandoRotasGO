package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rodrigocostarcs/explorandoRotasGO/cidades"
)

type Eucatur []struct {
	Key                         string  `json:"key"`
	ID                          int     `json:"id"`
	OriginSectionalID           int     `json:"origin_sectional_id"`
	OriginSectionalDescription  string  `json:"origin_sectional_description"`
	DestinySectionalID          int     `json:"destiny_sectional_id"`
	DestinySectionalDescription string  `json:"destiny_sectional_description"`
	TravelConnectionID          int     `json:"travel_connection_id"`
	DatetimeDeparture           string  `json:"datetime_departure"`
	DatetimeArrival             string  `json:"datetime_arrival"`
	Duration                    string  `json:"duration"`
	FreeSeats                   int     `json:"free_seats"`
	Price                       float64 `json:"price"`
	PricesStartingFrom          float64 `json:"prices_starting_from"`
	Points                      int     `json:"points"`
	WomanSpace                  bool    `json:"woman_space"`
	Tax                         int     `json:"tax"`
	Info                        string  `json:"info"`
	Items                       []struct {
		ID                int     `json:"id"`
		RoadItemID        int     `json:"road_item_id"`
		GatewayID         int     `json:"gateway_id"`
		GatewayType       string  `json:"gateway_type"`
		StationIDOrigin   int     `json:"station_id_origin"`
		StationIDDestiny  int     `json:"station_id_destiny"`
		ServiceTravel     string  `json:"service_travel"`
		ReferenceTravel   string  `json:"reference_travel"`
		DatetimeDeparture string  `json:"datetime_departure"`
		DatetimeArrival   string  `json:"datetime_arrival"`
		Duration          string  `json:"duration"`
		FreeSeats         int     `json:"free_seats"`
		Tariff            int     `json:"tariff"`
		Insurance         int     `json:"insurance"`
		Fee               int     `json:"fee"`
		TravelItemToll    int     `json:"travel_item_toll"`
		Price             float64 `json:"price"`
		PricePromotional  int     `json:"price_promotional"`
		Tax               int     `json:"tax"`
		Class             struct {
			ID        int    `json:"id"`
			LongName  string `json:"long_name"`
			ShortName string `json:"short_name"`
		} `json:"class"`
		Company struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
			Logo string `json:"logo"`
		} `json:"company"`
		LineCode  string `json:"line_code"`
		Direction string `json:"direction"`
		VehicleID int    `json:"vehicle_id"`
		Tariffs   []struct {
			Tariff   float64 `json:"tariff"`
			Cashback struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"cashback"`
			Insurance        int     `json:"insurance"`
			Toll             int     `json:"toll"`
			BoardingFee      float64 `json:"boarding_fee"`
			Ferry            int     `json:"ferry"`
			Additional       int     `json:"additional"`
			Calculation      float64 `json:"calculation"`
			Referential      float64 `json:"referential"`
			ResourceDiscount int     `json:"resource_discount"`
			Amount           float64 `json:"amount"`
			Points           int     `json:"points"`
			PricePromotional int     `json:"price_promotional"`
			PriceDiscount    int     `json:"price_discount"`
			TotalKm          float64 `json:"total_km"`
			ResourceKm       int     `json:"resource_km"`
			BeforeResource   int     `json:"before_resource"`
			ResourceType     string  `json:"resource_type"`
			Category         struct {
				MapaPoltronaID    int    `json:"MapaPoltrona_Id"`
				TipoVeiculoID     int    `json:"TipoVeiculoId"`
				Descricao         string `json:"Descricao"`
				DescricaoReduzida string `json:"DescricaoReduzida"`
				CategoriaSAT      string `json:"CategoriaSAT"`
				PoltronaInicial   int    `json:"PoltronaInicial"`
				PoltronaFinal     int    `json:"PoltronaFinal"`
				WomanSpace        bool   `json:"woman_space"`
			} `json:"category"`
		} `json:"tariffs"`
		ElderlyQuota struct {
			Type            string `json:"type"`
			PercentDiscount int    `json:"percent_discount"`
		} `json:"elderly_quota"`
	} `json:"items"`
}

func main() {
	http.HandleFunc("/listarcidades", pegarCidadesDisponveis)
	http.HandleFunc("/consultarHorarios", consultaHorariosDisponveis)
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

func consultaHorariosDisponveis(w http.ResponseWriter, r *http.Request) {

	origem := r.URL.Query().Get("origem")
	destino := r.URL.Query().Get("destino")
	data := r.URL.Query().Get("data")
	points := false

	err := cidades.ValidaDadosHelp(origem, destino, data)

	if err != nil {
		errorMsg := map[string]interface{}{
			"error": err.Error(), // Incluir a mensagem de erro no JSON
		}

		// Converter o mapa em bytes JSON
		jsonData, _ := json.Marshal(errorMsg)

		// Configurar o cabeçalho da resposta para indicar que o conteúdo é JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // Definir o código de status HTTP para Bad Request (400)

		// Escrever o JSON de erro na resposta
		_, _ = w.Write(jsonData)
		return
	}

	horarios, error := BuscaHorarios(origem, destino, data, points)

	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(cep) //pega a response e encoda ela e joga o CEP nela.

	result, err := json.Marshal(horarios)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(result)
}

func BuscaHorarios(origem string, destino string, data string, points bool) (*Eucatur, error) {
	pointsStr := strconv.FormatBool(points)
	resp, error := http.Get("https://api-v4.eucatur.com.br/road/travels/search?origin_sectional_id=" + origem + "&destiny_sectional_id=" + destino + "&departure_date=" + data + "&points=" + pointsStr)

	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	body, error := ioutil.ReadAll(resp.Body)

	if error != nil {
		return nil, error
	}

	var horarios Eucatur

	error = json.Unmarshal(body, &horarios)

	if error != nil {
		return nil, error
	}

	return &horarios, nil

}
