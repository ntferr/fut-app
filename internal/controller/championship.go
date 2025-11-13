package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fut-app/internal/model"
	"github.com/fut-app/internal/service"
	"github.com/labstack/echo/v4"
)

type Champion struct {
	Service service.FootballAPI
}

func NewChampion(footbal service.Football) Champion {
	return Champion{
		Service: &footbal,
	}
}

func (cham *Champion) Championship(c echo.Context) error {
	competitionResponse, err := cham.Service.CompetitionList(c.Request().Context())
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			err,
		)
	}
	formatedCompetition := FilterCompetitions(*competitionResponse)
	return c.JSON(
		http.StatusOK,
		formatedCompetition,
	)
}

func FilterCompetitions(data model.CompetitionResponse) []model.FormattedCompetition {
	var result []model.FormattedCompetition

	for _, comp := range data.Competitions {
		formatted := model.FormattedCompetition{
			ID:        fmt.Sprintf("campeonato_%03d", comp.ID),
			Nome:      comp.Name,
			Temporada: extractYear(comp.CurrentSeason.StartDate),
		}
		result = append(result, formatted)
	}

	return result
}

func extractYear(date string) string {
	parts := strings.Split(date, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
