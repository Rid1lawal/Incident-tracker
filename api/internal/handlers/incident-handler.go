package handlers

import (
	"fmt"
	"net/http"

	"github.com/Reazy-ai/incident-tracker/internal/models"
	"github.com/Reazy-ai/incident-tracker/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type IncidentHandler struct {
	Repo *repositories.IncidentRepository
}

func NewIncidentHandler(
	repo *repositories.IncidentRepository,
) *IncidentHandler {
	return &IncidentHandler{
		Repo: repo,
	}
}

func (h *IncidentHandler) CreateIncident(c *gin.Context) {

	var req models.CreateIncidentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	incident, err := h.Repo.CreateIncident(req)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Error().
			Err(err).
			Msg("failed to create incident")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create incident",
		})

		return
	}

	c.JSON(http.StatusCreated, incident)
}

func (h *IncidentHandler) GetIncidents(c *gin.Context) {

	incidents, err := h.Repo.GetIncidents()
	if err != nil {

		log.Error().
			Err(err).
			Msg("failed to fetch incidents")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch incidents",
		})

		return
	}

	c.JSON(http.StatusOK, incidents)
}
