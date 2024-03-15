package service_test

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/application/contract"
	"github.com/diegoclair/log-parser/application/dto"
	"github.com/diegoclair/log-parser/application/service"
	"github.com/stretchr/testify/assert"
)

func getQuakeService(t *testing.T) contract.QuakeService {
	services, err := service.New(logger.NewNoop())
	assert.NoError(t, err)

	return services.QuakeService
}

type args struct {
	lines []string
}

type test struct {
	name      string
	args      args
	setupTest func(a args)
	want      []dto.Report
}

var (
	initGameLine  = `0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\`
	userInfoLine1 = `20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\`
	userInfoLine2 = `20:34 ClientUserinfoChanged: 3 n\Mocinha\t\0\model\sarge/default\hmodel\`
	killLine      = `22:06 Kill: 3 2 7: Isgalamido killed Mocinha by MOD_ROCKET_SPLASH`
)

var sendLastGameReportTests = []test{
	{
		name: "If we don't have any line, we should not send any report to the writerChan",
		args: args{},
		want: []dto.Report{},
	},
	{
		name: "should be sent a report by last report function",
		args: args{
			lines: []string{
				initGameLine, // we need a initGame line to be possible to send a report
				userInfoLine1,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Isgalamido"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
}

var processNewGameEventTests = []test{
	{
		name: "should process a new game event",
		args: args{
			lines: []string{
				initGameLine,
			},
		},
		want: []dto.Report{},
	},
	{
		name: "should process a new game event with players and send a report",
		args: args{
			lines: []string{
				initGameLine,
				userInfoLine1,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Isgalamido"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
	{
		name: "should process two new game events and send two reports",
		args: args{
			lines: []string{
				initGameLine,
				userInfoLine1,
				initGameLine,
				userInfoLine2,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Isgalamido"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
			{
				GameName:     "game_002",
				Players:      []string{"Mocinha"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
}

func TestQuakeService_StartExtractingData(t *testing.T) {
	svc := getQuakeService(t)
	ctx := context.Background()

	tests := []test{}
	tests = slices.Concat(
		tests,
		sendLastGameReportTests,
		processNewGameEventTests,
	)

	tests = append(tests,
		test{
			name: "should skip first line that is not a new game event (continue for gameCount == 0)",
			args: args{
				lines: []string{
					"20:34 ----------------------------",
				},
			},
			want: []dto.Report{},
		},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTest != nil {
				tt.setupTest(tt.args)
			}

			writerChan := make(chan dto.Report)
			lineChan := make(chan string)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				svc.StartExtractingData(ctx, lineChan, writerChan)
			}()

			reports := []dto.Report{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				for report := range writerChan {
					reports = append(reports, report)
				}
			}()

			for _, line := range tt.args.lines {
				fmt.Println("Sending line: ", line)
				lineChan <- line
			}
			close(lineChan)

			fmt.Println("Waiting for writerChan")

			if len(tt.args.lines) == 0 {
				assert.Equal(t, 0, len(writerChan))
				return
			}

			wg.Wait()

			assert.Equal(t, tt.want, reports)

		})
	}
}
