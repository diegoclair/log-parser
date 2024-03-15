package service_test

import (
	"context"
	"reflect"
	"slices"
	"sort"
	"sync"
	"testing"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/application/contract"
	"github.com/diegoclair/log-parser/application/dto"
	"github.com/diegoclair/log-parser/application/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	initGameEvent          = `0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\`
	userTest1Event         = `20:34 ClientUserinfoChanged: 2 n\Test1\t\0\model\xian/default\hmodel\`
	userTest2Event         = `20:34 ClientUserinfoChanged: 3 n\Test2\t\0\model\sarge/default\hmodel\`
	killEvent              = `22:06 Kill: 3 2 7: Test2 killed Test1 by MOD_ROCKET_SPLASH`
	killEventDifferentName = `22:06 Kill: 3 2 7: Xxxxx1 killed Xxxxx2 by MOD_ROCKET_SPLASH` // name should be found by id 3 and 2
	worldKillEvent         = `22:06 Kill: 1022 2 19: <world> killed Test1 by MOD_TRIGGER_HURT`
	samePlayerKillEvent    = `22:06 Kill: 3 3 7: Test2 killed Test2 by MOD_ROCKET_SPLASH`
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
				initGameEvent, // we need a initGame line to be possible to send a report
				userTest1Event,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Test1"},
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
				initGameEvent,
			},
		},
		want: []dto.Report{},
	},
	{
		name: "should process a new game event with players and send a report",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Test1"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
	{
		name: "should process two new game events and send two reports",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				initGameEvent,
				userTest2Event,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Test1"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
			{
				GameName:     "game_002",
				Players:      []string{"Test2"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
}

var processUserChangedEventTests = []test{
	{
		name: "should process a user changed event",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Test1"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
	{
		name: "if we have two userChangedEvent and one return error, should send only one user to report",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				`ClientUserinfoChanged: 3 n\Test2\t\0\model\sarge/default\hmodel\`,
			},
		},
		want: []dto.Report{
			{
				GameName:     "game_001",
				Players:      []string{"Test1"},
				Kills:        make(map[string]int),
				KillsByMeans: make(map[string]int),
			},
		},
	},
}

var processKillEventTests = []test{
	{
		name: "should process a kill event",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				userTest2Event,
				killEvent,
			},
		},
		want: []dto.Report{
			{
				GameName:   "game_001",
				TotalKills: 1,
				Players:    []string{"Test1", "Test2"},
				Kills: map[string]int{
					"Test2": 1,
				},
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
				},
			},
		},
	},
	{
		name: "should skip kill event if line is a invalid kill event",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				userTest2Event,
				` Kill: 3 2 7: Test2 killed Test1 by MOD_ROCKET_SPLASH`,
				killEvent,
			},
		},
		want: []dto.Report{
			{
				GameName:   "game_001",
				TotalKills: 1,
				Players:    []string{"Test1", "Test2"},
				Kills: map[string]int{
					"Test2": 1,
				},
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
				},
			},
		},
	},
	{
		name: "should process a kill event and get player name by id",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				userTest2Event,
				killEventDifferentName,
			},
		},
		want: []dto.Report{
			{
				GameName:   "game_001",
				TotalKills: 1,
				Players:    []string{"Test1", "Test2"},
				Kills: map[string]int{
					"Test2": 1,
				},
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
				},
			},
		},
	},
	{
		name: "should count worldPlayer kill as total kills and decrease killed player kills",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				userTest2Event,
				worldKillEvent,
			},
		},
		want: []dto.Report{
			{
				GameName:   "game_001",
				TotalKills: 1,
				Players:    []string{"Test1", "Test2"},
				Kills: map[string]int{
					"Test1": -1,
				},
				KillsByMeans: map[string]int{
					"MOD_TRIGGER_HURT": 1,
				},
			},
		},
	},
	{
		name: "should not count as kill for an user if the killer is the same as the killed",
		args: args{
			lines: []string{
				initGameEvent,
				userTest1Event,
				samePlayerKillEvent,
			},
		},
		want: []dto.Report{
			{
				GameName:   "game_001",
				TotalKills: 1,
				Players:    []string{"Test1"},
				Kills:      make(map[string]int),
				KillsByMeans: map[string]int{
					"MOD_ROCKET_SPLASH": 1,
				},
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
		processUserChangedEventTests,
		processKillEventTests,
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
		test{
			name: "item line should be skipped and not processed",
			args: args{
				lines: []string{
					initGameEvent,
					"20:34 Item: 2 weapon_rocketlauncher",
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

			for i := range tt.args.lines {
				lineChan <- tt.args.lines[i]
			}
			close(lineChan)

			if len(tt.args.lines) == 0 {
				assert.Equal(t, 0, len(writerChan))
				return
			}

			wg.Wait()

			require.Equal(t, len(tt.want), len(reports))
			for i := range reports {
				sort.Strings(reports[i].Players)
				sort.Strings(tt.want[i].Players)
				require.Equal(t, tt.want[i], reports[i])
			}

			assert.True(t, reflect.DeepEqual(tt.want, reports))
		})
	}
}
