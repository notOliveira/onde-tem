package http

import (
	"encoding/json"
	"net/http"

	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/usecase"
	"github.com/notOliveira/onde-tem/internal/infra/logger"
)

type EstablishmentHandler struct {
	createUC *usecase.CreateEstablishmentUseCase
	log      *logger.Logger
}

func NewEstablishmentHandler(createUC *usecase.CreateEstablishmentUseCase, log *logger.Logger) *EstablishmentHandler {
	return &EstablishmentHandler{
		createUC: createUC,
		log:      log,
	}
}

// Criamos um DTO específico para o Request HTTP
type createRequestPayload struct {
	Name    string   `json:"name"`
	Slug    string   `json:"slug"`
	Types   []string `json:"types"`
	Address struct {
		Street   string `json:"street"`
		Number   string `json:"number"`
		District string `json:"district"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		ZipCode  string `json:"zipCode"`
	} `json:"address"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

func (h *EstablishmentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Ler o JSON da requisição
	var payload createRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Montar o input para o Use Case (O handler faz a ponte HTTP -> Domínio)
	loc, _ := domain.NewLocation(payload.Location.Lat, payload.Location.Lon)
	addr, _ := domain.NewAddress(payload.Address.Street, payload.Address.Number, payload.Address.District, payload.Address.City, payload.Address.State, payload.Address.Country, payload.Address.ZipCode)

	input := usecase.CreateEstablishmentInput{
		Name:     payload.Name,
		Slug:     payload.Slug,
		Types:    payload.Types,
		Location: loc,
		Address:  addr,
	}

	// 3. Executar o Use Case
	est, err := h.createUC.Execute(r.Context(), input)
	if err != nil {
		h.log.Errorf("Error creating establishment: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Retornar Sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      est.ID().String(),
		"message": "Establishment created successfully",
	})
}
