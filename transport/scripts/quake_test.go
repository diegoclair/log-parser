package scripts_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/transport/scripts"
)

func TestReadLinesFromQuakeLog(t *testing.T) {
	ctx := context.Background()

	type args struct {
		file     io.Reader
		lineChan chan string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should read lines from file",
			args: args{
				file:     strings.NewReader("line1\nline2\nline3"),
				lineChan: make(chan string),
			},
			want: []string{"line1", "line2", "line3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := scripts.NewQuakeLogParser(logger.NewNoop())

			go func() {
				for line := range tt.args.lineChan {
					if line != tt.want[0] {
						t.Errorf("ReadLinesFromQuakeLog() = %v, want %v", line, tt.want[0])
					}
					tt.want = tt.want[1:]
				}
			}()

			q.ReadLinesFromQuakeLog(ctx, tt.args.file, tt.args.lineChan)
		})
	}
}
