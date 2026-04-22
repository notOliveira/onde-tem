package http

import (
	"github.com/gin-gonic/gin"
	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/usecase"
	"net/http"
)

type EstablishmentHandler struct {
	createUC *usecase.CreateEstablishmentUseCase
}

func NewEstablishmentHandler(createUC *usecase.CreateEstablishmentUseCase) *EstablishmentHandler {
	return &EstablishmentHandler{createUC: createUC}
}

type createRequest struct {
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
		ZipCode  string `json:"zip_code"`
	} `json:"address"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

func (h *EstablishmentHandler) HandleCreate(c *gin.Context) {

	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
		return
	}

	var req createRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	address, err := domain.NewAddress(
		req.Address.Street,
		req.Address.Number,
		req.Address.District,
		req.Address.City,
		req.Address.State,
		req.Address.Country,
		req.Address.ZipCode,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := domain.NewLocation(
		req.Location.Lat,
		req.Location.Lon,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.CreateEstablishmentInput{
		Name:     req.Name,
		Slug:     req.Slug,
		Types:    req.Types,
		Location: location,
		Address:  address,
	}

	result, err := h.createUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      result.ID().String(),
		"name":    result.Name(),
		"slug":    result.Slug().String(),
		"message": "Establishment created successfully",
	})
}
