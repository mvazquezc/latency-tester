package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetScheme(t *testing.T) {
	type args struct {
		inputTarget string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test HTTP URL",
			args: args{
				inputTarget: "http://google.com",
			},
			want: "http",
		},
		{
			name: "Test HTTPs URL",
			args: args{
				inputTarget: "https://google.com",
			},
			want: "http",
		},
		{
			name: "Test TCP Socket",
			args: args{
				inputTarget: "tcp://google.com",
			},
			want: "tcp",
		},
		{
			name: "Test default",
			args: args{
				inputTarget: "some://malformed.url",
			},
			want: "http",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetScheme(tt.args.inputTarget)
			assert.Equal(t, got, tt.want, "Expected scheme %s, got: %s", tt.want, got)
		})
	}
}

func TestIntervalTimeToSeconds(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test seconds to seconds",
			args: args{
				interval: "30s",
			},
			want: 30,
		},
		{
			name: "Test minutes to seconds",
			args: args{
				interval: "5m",
			},
			want: 300,
		},
		{
			name: "Test hours to seconds",
			args: args{
				interval: "1h",
			},
			want: 3600,
		},
		{
			name: "Test invalid",
			args: args{
				interval: "1123j",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntervalTimeToSeconds(tt.args.interval)
			assert.Equal(t, got, tt.want, "Expected seconds %d, got: %d", tt.want, got)
		})
	}
}

func TestValidateIntervalTime(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test invalid",
			args: args{
				interval: "1123j",
			},
			want: false,
		},
		{
			name: "Test valid seconds",
			args: args{
				interval: "3600s",
			},
			want: true,
		},
		{
			name: "Test valid minutes",
			args: args{
				interval: "60m",
			},
			want: true,
		},
		{
			name: "Test valid hours",
			args: args{
				interval: "1h",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateIntervalTime(tt.args.interval)
			assert.Equal(t, got, tt.want, "Expected %v, got: %v", tt.want, got)
		})
	}
}

func TestValidateTarget(t *testing.T) {
	type args struct {
		inputTarget string
	}
	tests := []struct {
		name              string
		args              args
		wantedValidTarget bool
		wantedScheme      string
	}{
		{
			name: "Test valid HTTP target",
			args: args{
				inputTarget: "http://google.com",
			},
			wantedValidTarget: true,
			wantedScheme:      "http",
		},
		{
			name: "Test valid HTTPs target",
			args: args{
				inputTarget: "https://google.com",
			},
			wantedValidTarget: true,
			wantedScheme:      "https",
		},
		{
			name: "Test valid TCP target",
			args: args{
				inputTarget: "tcp://google.com",
			},
			wantedValidTarget: true,
			wantedScheme:      "tcp",
		},
		{
			name: "Test non-valid target",
			args: args{
				inputTarget: "notValid://google.com",
			},
			wantedValidTarget: false,
			wantedScheme:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValidTarget, gotScheme := ValidateTarget(tt.args.inputTarget)
			assert.Equal(t, gotValidTarget, tt.wantedValidTarget, "Expected %v, got: %v", tt.wantedValidTarget, gotValidTarget)
			assert.Equal(t, gotScheme, tt.wantedScheme, "Expected %s, got: %s", tt.wantedScheme, gotScheme)
		})
	}
}
