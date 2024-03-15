package writer_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/application/dto"
	"github.com/diegoclair/log-parser/infra/writer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type writerMock struct {
	wrote []byte
}

func newWriterMock() *writerMock {
	return &writerMock{}
}

func (w *writerMock) Write(b []byte) (n int, err error) {
	report := dto.QuakeDataReport{}

	err = json.Unmarshal(b, &report)
	if err != nil {
		return 0, err
	}

	if _, ok := report["game_error"]; ok {
		return 0, errors.New("error")
	}

	w.wrote = b

	return 0, nil
}

func (w *writerMock) wroteBytes() []byte {
	return w.wrote
}

func TestWriter_StartWriting(t *testing.T) {
	ctx := context.Background()
	writerMock := newWriterMock()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	loggerMock := logger.NewMockLogger(ctrl)

	wr := writer.NewWriter(ctx, writerMock, loggerMock)

	type args struct {
		data dto.Report
	}

	tests := []struct {
		name      string
		args      args
		setupTest func(a args, log *logger.MockLogger)
		want      dto.QuakeDataReport
		wantErr   bool
	}{
		{
			name: "should write data to file",
			args: args{
				data: dto.Report{
					GameName:   "game1",
					TotalKills: 35,
					Players:    []string{"player1", "player2"},
					Kills:      map[string]int{"player1": 10, "player2": 23},
				},
			},
			want: dto.QuakeDataReport{
				"game1": {
					GameName:   "game1",
					TotalKills: 35,
					Players:    []string{"player1", "player2"},
					Kills:      map[string]int{"player1": 10, "player2": 23},
				},
			},
		},
		{
			name: "should log error when fail to write data",
			args: args{
				data: dto.Report{
					GameName: "game_error",
				},
			},
			setupTest: func(a args, log *logger.MockLogger) {
				log.EXPECT().Errorf(ctx, "Error to write data: %v", errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTest != nil {
				tt.setupTest(tt.args, loggerMock)
			}

			data := make(chan dto.Report)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				wr.StartWriting(ctx, data)
			}()

			data <- tt.args.data
			close(data)
			wg.Wait()

			b, err := json.Marshal(tt.want)
			require.NoError(t, err)

			if !tt.wantErr {
				require.Equal(t, b, writerMock.wroteBytes())
			}
		})
	}
}
