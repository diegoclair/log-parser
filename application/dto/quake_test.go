package dto_test

import (
	"slices"
	"testing"

	"github.com/diegoclair/log-parser/application/dto"
)

func TestQuakeData_ToReport(t *testing.T) {
	quakeData := dto.QuakeData{
		TotalKills: 10,
		Players:    map[int]string{0: "Player1", 1: "Player2", 2: "Player3"},
		Kills: map[int]int{
			0: 5,
			1: 3,
			2: 2,
		},
		KillsByMeans: map[string]int{
			"MOD_TRIGGER_HURT": 8,
			"MOD_FALLING":      2,
		},
	}

	gameName := "game_test"
	want := dto.Report{
		GameName:   gameName,
		TotalKills: 10,
		Players:    []string{"Player1", "Player2", "Player3"},
		Kills: map[string]int{
			"Player1": 5,
			"Player2": 3,
			"Player3": 2,
		},
		KillsByMeans: map[string]int{
			"MOD_TRIGGER_HURT": 8,
			"MOD_FALLING":      2,
		},
	}

	got := quakeData.ToReport(gameName)

	if got.GameName != want.GameName {
		t.Errorf("ToReport() got.GameName = %v, want.GameName %v", got.GameName, want.GameName)
	}

	if got.TotalKills != want.TotalKills {
		t.Errorf("ToReport() got.TotalKills = %v, want.TotalKills %v", got.TotalKills, want.TotalKills)
	}

	if len(got.Players) != len(want.Players) {
		t.Errorf("ToReport() got.Players length = %v, want.Players length %v", len(got.Players), len(want.Players))
	}

	for _, player := range quakeData.Players {
		if !slices.Contains(got.Players, player) {
			t.Errorf("ToReport() got.Players does not contain %s", player)
		}
	}

	if len(got.Kills) != len(want.Kills) {
		t.Errorf("ToReport() got.Kills length = %v, want.Kills length %v", len(got.Kills), len(want.Kills))
	}

	for player, kills := range got.Kills {
		if kills != want.Kills[player] {
			t.Errorf("ToReport() got.Kills[%s] = %v, want.Kills[%s] %v", player, kills, player, want.Kills[player])
		}
	}

	if len(got.KillsByMeans) != len(want.KillsByMeans) {
		t.Errorf("ToReport() got.KillsByMeans length = %v, want.KillsByMeans length %v", len(got.KillsByMeans), len(want.KillsByMeans))
	}

	for means, count := range got.KillsByMeans {
		if count != want.KillsByMeans[means] {
			t.Errorf("ToReport() got.KillsByMeans[%s] = %v, want.KillsByMeans[%s] %v", means, count, means, want.KillsByMeans[means])
		}
	}
}
