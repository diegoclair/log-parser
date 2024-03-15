package service

import (
	"context"
	"fmt"

	"github.com/diegoclair/log-parser/application"
	"github.com/diegoclair/log-parser/application/contract"
	"github.com/diegoclair/log-parser/application/dto"
	"github.com/diegoclair/log-parser/domain/entity"
)

// quakeService is a struct that implements the contract.QuakeService interface.
type quakeService struct {
	svc *service
}

// newQuakeService creates a new instance of the quakeService struct.
func newQuakeService(svc *service) contract.QuakeService {
	return &quakeService{
		svc: svc,
	}
}

// StartExtractingData is a method that starts extracting data from the rawLinesChan channel and writes the report to the writerChan channel.
func (s *quakeService) StartExtractingData(ctx context.Context, rawLinesChan <-chan string, writerChan chan<- dto.Report) {
	var gameData dto.QuakeData
	gameCount := 0

	for line := range rawLinesChan {
		processed := processNewGameEvent(line, gameCount, &gameData, writerChan)
		if processed {
			gameCount++
			continue
		}

		// if gameCount still 0 here, then we don't have a game yet
		if gameCount == 0 {
			continue
		}

		processed, err := s.processUserChangedEvent(ctx, line, &gameData)
		if err != nil || processed {
			continue
		}

		// should use return values to do continue if this function is not the last in the chain
		_, _ = s.processKillEvent(ctx, line, &gameData)
	}

	s.sendLastGameReport(&gameData, gameCount, writerChan)

	close(writerChan)
}

func (s *quakeService) sendLastGameReport(gameData *dto.QuakeData, gameCount int, writerChan chan<- dto.Report) {
	// if there are no players, we don't consider it a game
	if len(gameData.Players) == 0 {
		return
	}

	writerChan <- gameData.ToReport(generateGameName(gameCount))
}

// processNewGameEvent is a method that processes the new game event from the given line and updates the gameData.
// It also writes the gameData to the writerChan channel if the gameCount is greater than 0.
// It returns true if the line is a new game event.
func processNewGameEvent(line string, gameCount int, gameData *dto.QuakeData, writerChan chan<- dto.Report) (processed bool) {
	if entity.IsNewGameByLine(line) {
		if gameCount > 0 {
			writerChan <- gameData.ToReport(generateGameName(gameCount))
		}

		// reset game data for the new game stats
		gameData.Reset()
		return true
	}

	return false
}

// processUserChangedEvent is a method that processes the user changed event from the given line and updates the gameData.
// It returns true if the line is a user changed event.
func (s *quakeService) processUserChangedEvent(ctx context.Context, line string, gameData *dto.QuakeData) (processed bool, err error) {
	hasPlayer, player, err := entity.IsUserInfoChangedByLine(line)
	if err != nil {
		s.svc.log.Error(ctx, fmt.Sprintf("Error to process user changed event: %v", err))
		return false, err
	}

	if hasPlayer {
		gameData.Players[player.ID] = player.Name
		return true, nil
	}

	return false, nil
}

// processKillEvent is a method that processes the kill event from the given line and updates the gameData.
// It returns true if the line is a kill event.
func (s *quakeService) processKillEvent(ctx context.Context, line string, gameData *dto.QuakeData) (processed bool, err error) {
	hasKill, killInfo, err := entity.IsKillByLine(line)
	if err != nil {
		s.svc.log.Error(ctx, fmt.Sprintf("Error to process kill event: %v", err))
		return false, err
	}

	if hasKill {
		gameData.TotalKills++
		gameData.KillsByMeans[killInfo.DeathCause]++

		if killInfo.KillerID == application.WorldPlayerID {
			gameData.Kills[killInfo.KilledID]--
			return true, nil
		}

		// do not count as kill for an user if the killer is the same as the killed
		if killInfo.KillerID == killInfo.KilledID {
			return true, nil
		}

		gameData.Kills[killInfo.KillerID]++
		return true, nil
	}

	return false, nil
}

// generateGameName generates a game name based on the gameCount.
func generateGameName(gameCount int) string {
	return fmt.Sprintf("game_%03d", gameCount)
}
