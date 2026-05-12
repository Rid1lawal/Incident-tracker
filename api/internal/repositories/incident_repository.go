package repositories

import (
	"context"
	"time"

	"github.com/Reazy-ai/incident-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository struct {
	DB *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{
		DB: db,
	}
}

func (r *IncidentRepository) CreateIncident(
	req models.CreateIncidentRequest,
) (*models.Incident, error) {

	query := `
		INSERT INTO incidents (title, description, severity)
		VALUES ($1, $2, $3)
		RETURNING id, status, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var incident models.Incident

	err := r.DB.QueryRow(
		ctx,
		query,
		req.Title,
		req.Description,
		req.Severity,
	).Scan(
		&incident.ID,
		&incident.Status,
		&incident.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	incident.Title = req.Title
	incident.Description = req.Description
	incident.Severity = req.Severity

	return &incident, nil
}

func (r *IncidentRepository) GetIncidents() ([]models.Incident, error) {

	query := `
		SELECT id, title, description, severity, status, created_at
		FROM incidents
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var incidents []models.Incident

	for rows.Next() {
		var incident models.Incident

		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.Severity,
			&incident.Status,
			&incident.CreatedAt,
		)

		if err != nil {
			continue
		}

		incidents = append(incidents, incident)
	}

	return incidents, nil
}
