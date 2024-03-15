package dto

// QuakeData represents the data structure for storing quake game information.
type QuakeData struct {
	TotalKills   int            // TotalKills represents the total number of kills in the game.
	Players      map[int]string // Players represents the mapping of player IDs to player names.
	Kills        map[int]int    // Kills represents the mapping of player IDs to their respective kill counts.
	KillsByMeans map[string]int // KillsByMeans represents the mapping of kill means to their respective counts.
}

func (q *QuakeData) Reset() {
	q.TotalKills = 0
	q.Players = make(map[int]string)
	q.Kills = make(map[int]int)
	q.KillsByMeans = make(map[string]int)
}

// ToReport converts the QuakeData into a Report object.
func (q *QuakeData) ToReport(gameName string) Report {
	report := Report{
		GameName:     gameName,
		TotalKills:   q.TotalKills,
		Players:      make([]string, 0),
		Kills:        make(map[string]int),
		KillsByMeans: q.KillsByMeans,
	}

	for _, player := range q.Players {
		report.Players = append(report.Players, player)
	}

	for playerID, playerKills := range q.Kills {
		report.Kills[q.Players[playerID]] = playerKills
	}

	return report
}

// Report represents the report structure for a Quake game.
type Report struct {
	GameName     string         `json:"-"`              // GameName represents the name of the game.
	TotalKills   int            `json:"total_kills"`    // TotalKills represents the total number of kills in the game.
	Players      []string       `json:"players"`        // Players represents the list of player names.
	Kills        map[string]int `json:"kills"`          // Kills represents the mapping of player names to their respective kill counts.
	KillsByMeans map[string]int `json:"kills_by_means"` // KillsByMeans represents the mapping of kill means to their respective counts.
}

// QuakeDataReport represents a collection of Quake game reports.
type QuakeDataReport map[string]Report
