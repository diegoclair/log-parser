package entity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// KillEvent represents a kill event in the game log.
type KillEvent struct {
	KillerID   int    // ID of the player who performed the kill.
	KilledID   int    // ID of the player who was killed.
	DeathCause string // Cause of the death.
}

var (
	TimeRegex        = regexp.MustCompile(`(\d+:\d+)`)
	KillRegex        = regexp.MustCompile(fmt.Sprintf(`%s Kill: (\d+) (\d+) \d+: (.+) killed (.+) by (.+)`, TimeRegex))
	UserChangedRegex = regexp.MustCompile(fmt.Sprintf(`%s ClientUserinfoChanged: (\d+) n\\(.+)\\t\\`, TimeRegex))
)

// IsKillByLine checks if a given line represents a kill event.
// If it is a kill event, it returns true and the corresponding KillEvent object.
// If it is not a kill event, it returns false.
func IsKillByLine(line string) (bool, KillEvent, error) {
	if strings.Contains(line, " Kill: ") {
		killData, err := getKillData(line)
		if err != nil {
			return false, KillEvent{}, err
		}

		return true, killData, err
	}

	return false, KillEvent{}, nil
}

// getKillData extracts the kill data from a kill event line.
func getKillData(line string) (KillEvent, error) {
	matches := KillRegex.FindStringSubmatch(line)

	if len(matches) != 7 {
		return KillEvent{}, fmt.Errorf("invalid number of matches in kill line: %s", line)
	}

	killerID, err := strconv.Atoi(matches[2])
	if err != nil {
		return KillEvent{}, err
	}

	KilledID, err := strconv.Atoi(matches[3])
	if err != nil {
		return KillEvent{}, err
	}

	return KillEvent{
		KillerID:   killerID,
		KilledID:   KilledID,
		DeathCause: matches[6],
	}, nil
}

// IsNewGameByLine checks if a given line represents the start of a new game.
// It returns true if the line contains "InitGame", indicating a new game has started.
func IsNewGameByLine(line string) bool {
	return strings.Contains(line, "InitGame")
}

// Player represents a player in the game.
type Player struct {
	ID   int    // ID of the player.
	Name string // Name of the player.
}

// IsUserInfoChangedByLine checks if a given line represents a player's information change event.
// If it is a player's information change event, it returns true and the corresponding Player object.
// If it is not a player's information change event, it returns false.
func IsUserInfoChangedByLine(line string) (bool, Player, error) {
	if strings.Contains(line, "ClientUserinfoChanged") {
		player, err := getPlayerName(line)
		if err != nil {
			return false, Player{}, err
		}

		return true, player, nil
	}

	return false, Player{}, nil
}

// getPlayerName extracts the player's name from a player's information change event line.
func getPlayerName(line string) (Player, error) {
	matches := UserChangedRegex.FindStringSubmatch(line)

	if len(matches) != 4 {
		return Player{}, fmt.Errorf("invalid number of matches in user changed line: %s", line)
	}

	id, err := strconv.Atoi(matches[2])
	if err != nil {
		return Player{}, err
	}

	return Player{
		ID:   id,
		Name: matches[3],
	}, nil
}
