package entity_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/diegoclair/log-parser/domain/entity"
)

func TestIsKillByLine(t *testing.T) {
	type args struct {
		line string
	}

	tests := []struct {
		name      string
		args      args
		setupTest func(args args)
		want      entity.KillEvent
		wantErr   bool
		wantBool  bool
	}{
		{
			name: "Should return a kill event",
			args: args{
				line: "20:54 Kill: 1022 2 19: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			},
			want: entity.KillEvent{
				KillerID:   1022,
				KilledID:   2,
				DeathCause: "MOD_TRIGGER_HURT",
			},
			wantErr:  false,
			wantBool: true,
		},
		{
			name: "Should return false if line does not contain Kill",
			args: args{
				line: "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\sarge\\hmodel\\sarge",
			},
			want:     entity.KillEvent{},
			wantErr:  false,
			wantBool: false,
		},
		{
			name: "Should return error if number of matches is different than 7",
			args: args{
				line: " Kill: 1022 2 19: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			},
			want:    entity.KillEvent{},
			wantErr: true,
		},
		{
			name: "Should return error if the second match is not a number",
			args: args{
				line: "20:54 Kill: a 2 19: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			},
			setupTest: func(args args) {
				entity.KillRegex = regexp.MustCompile(fmt.Sprintf(`%s Kill: (.+) (\d+) \d+: (.+) killed (.+) by (.+)`, entity.TimeRegex))
			},
			want:     entity.KillEvent{},
			wantErr:  true,
			wantBool: false,
		},
		{
			name: "Should return error if the third match is not a number",
			args: args{
				line: "20:54 Kill: 1022 a 19: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			},
			setupTest: func(args args) {
				entity.KillRegex = regexp.MustCompile(fmt.Sprintf(`%s Kill: (\d+) (.+) \d+: (.+) killed (.+) by (.+)`, entity.TimeRegex))
			},
			want:     entity.KillEvent{},
			wantErr:  true,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTest != nil {
				tt.setupTest(tt.args)
			}

			gotBool, got, err := entity.IsKillByLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsKillByLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsKillByLine() got = %v, want %v", got, tt.want)
			}

			if gotBool != tt.wantBool {
				t.Errorf("IsKillByLine() gotBool = %v, wantBool %v", gotBool, tt.wantBool)
			}
		})
	}
}

func TestIsNewGameByLine(t *testing.T) {
	type args struct {
		line string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true if line contains New game",
			args: args{
				line: "some InitGame text",
			},
			want: true,
		},
		{
			name: "Should return false if line does not contain New game",
			args: args{
				line: "20:54 Kill: 1022 2 19: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.IsNewGameByLine(tt.args.line)

			if got != tt.want {
				t.Errorf("IsNewGameByLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUserInfoChangedByLine(t *testing.T) {
	type args struct {
		line string
	}

	tests := []struct {
		name      string
		args      args
		setupTest func(args args)
		want      entity.Player
		wantErr   bool
		wantBool  bool
	}{
		{
			name: "Should return a player",
			args: args{
				line: "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\sarge\\hmodel\\sarge",
			},
			want: entity.Player{
				ID:   2,
				Name: "Isgalamido",
			},
			wantErr:  false,
			wantBool: true,
		},
		{
			name: "Should return error if number of matches is different than 4",
			args: args{
				line: " ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\sarge\\hmodel",
			},
			want:    entity.Player{},
			wantErr: true,
		},
		{
			name: "Should return false if line does not contain ClientUserinfoChanged",
			args: args{
				line: "20:34 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			},
			want:     entity.Player{},
			wantErr:  false,
			wantBool: false,
		},
		{
			name: "Should return error if the second match is not a number",
			args: args{
				line: "20:34 ClientUserinfoChanged: a n\\Isgalamido\\t\\0\\model\\sarge\\hmodel\\sarge",
			},
			setupTest: func(args args) {
				entity.UserChangedRegex = regexp.MustCompile(fmt.Sprintf(`%s ClientUserinfoChanged: (.+) n\\(.+)\\t\\`, entity.TimeRegex))
			},
			want:     entity.Player{},
			wantErr:  true,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTest != nil {
				tt.setupTest(tt.args)
			}

			gotBool, got, err := entity.IsUserInfoChangedByLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserInfoChangedByLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsUserInfoChangedByLine() got = %v, want %v", got, tt.want)
			}

			if gotBool != tt.wantBool {
				t.Errorf("IsUserInfoChangedByLine() gotBool = %v, wantBool %v", gotBool, tt.wantBool)
			}
		})
	}
}
