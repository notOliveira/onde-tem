package http

import (
	"github.com/gin-gonic/gin"
	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/usecase"
	"net/http"
)

type EstablishmentHandler struct {
	createUC *usecase.CreateEstablishmentUseCase
	listUC   *usecase.ListEstablishmentsUseCase
}

func NewEstablishmentHandler(
	createUC *usecase.CreateEstablishmentUseCase,
	listUC *usecase.ListEstablishmentsUseCase) *EstablishmentHandler {
	return &EstablishmentHandler{
		createUC: createUC,
		listUC:   listUC,
	}
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

type listRequest struct {
	Types  []string `form:"types"`
	Search string   `form:"search"`
	Limit  int      `form:"limit,default=20"`
	Offset int      `form:"offset,default=0"`
}

func (h *EstablishmentHandler) HandleCreate(c *gin.Context) {

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

func (h *EstablishmentHandler) HandleList(c *gin.Context) {

	var req listRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters", "details": err.Error()})
		return
	}

	filter := &domain.EstablishmentFilter{
		Types:  req.Types,
		Search: req.Search,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	establishments, err := h.listUC.Execute(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"establishments": establishments,
	})

}
