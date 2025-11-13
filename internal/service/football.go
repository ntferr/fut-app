package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/fut-app/internal/model"
)

var (
	ErrToCreateRequest        = errors.New("failed to create request")
	ErrToDoRequest            = errors.New("failed to do request")
	ErrToReadResponse         = errors.New("failed to read response")
	ErrToComunicateFootbalAPP = errors.New("failed to comunicate with footbal app")
	ErrToUnmarshal            = errors.New("failed to unmarshal response")
)

type Football struct {
	URL string
}

func NewFootball(url string) Football {
	return Football{
		URL: url,
	}
}

type FootballAPI interface {
	CompetitionList(ctx context.Context) (*model.CompetitionResponse, error)
}

func (f *Football) CompetitionList(ctx context.Context) (*model.CompetitionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.URL+"/v4/competitions", nil)
	if err != nil {
		return nil, ErrToCreateRequest
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrToDoRequest
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrToReadResponse
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrToComunicateFootbalAPP
	}

	var competitionResponse model.CompetitionResponse
	if err := json.Unmarshal(body, &competitionResponse); err != nil {
		return nil, ErrToUnmarshal
	}

	return &competitionResponse, nil
}
