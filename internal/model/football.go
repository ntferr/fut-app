package model

type CompetitionResponse struct {
	Competitions []Competition `json:"competitions"`
}

type Competition struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	CurrentSeason Season `json:"currentSeason"`
}

type Season struct {
	ID              int    `json:"id"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	CurrentMatchday *int   `json:"currentMatchday"`
	Winner          *Team  `json:"winner"`
}

type Team struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	TLA         string `json:"tla"`
	Crest       string `json:"crest"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Founded     int    `json:"founded"`
	ClubColors  string `json:"clubColors"`
	Venue       string `json:"venue"`
	LastUpdated string `json:"lastUpdated"`
}

type FormattedCompetition struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Temporada string `json:"temporada"`
}
