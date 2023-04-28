package commands

import (
	"github.com/mvazquezc/latency-tester/pkg/latency"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRunLatencyTestCmd(t *testing.T) {
	type args struct {
		target              string
		numberOfRuns        int
		waitIntervalSeconds int
		tcpPing             bool
	}
	tests := []struct {
		name         string
		args         args
		want         latency.AggregatedLatencyTestOutput
		wantErr      bool
		wantErrError string
	}{
		{
			name: "Test HTTP URL",
			args: args{
				target:              "http://google.com",
				numberOfRuns:        1,
				waitIntervalSeconds: 1,
			},
			want:    latency.AggregatedLatencyTestOutput{},
			wantErr: false,
		},

		{
			name: "Test HTTPs URL",
			args: args{
				target:              "https://redhat.com",
				numberOfRuns:        1,
				waitIntervalSeconds: 1,
			},
			want:    latency.AggregatedLatencyTestOutput{},
			wantErr: false,
		},
		{
			name: "Test non existing DNS",
			args: args{
				target:              "http://wefnwejklflwewef.com",
				numberOfRuns:        1,
				waitIntervalSeconds: 1,
			},
			want:         latency.AggregatedLatencyTestOutput{},
			wantErr:      true,
			wantErrError: "Get \"http://wefnwejklflwewef.com\": dial tcp: lookup wefnwejklflwewef.com: no such host",
		},
		{
			name: "Test HTTP Wrong URL",
			args: args{
				target:              "htt://google.com",
				numberOfRuns:        1,
				waitIntervalSeconds: 1,
			},
			want:         latency.AggregatedLatencyTestOutput{},
			wantErr:      true,
			wantErrError: "Get \"htt://google.com\": unsupported protocol scheme \"htt\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunLatencyTestCmd(tt.args.target, tt.args.numberOfRuns, tt.args.waitIntervalSeconds, tt.args.tcpPing)

			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.wantErrError, "Error not expected")
				return
			}

			if got.AvgDnsLookup <= 0 || got.AvgTcpConn <= 0 || got.AvgTlsHandshake <= 0 || got.AvgServerProcessing <= 0 || got.AvgContentTransfer <= 0 || got.AvgTotal <= 0 {
				assert.Greater(t, got.AvgDnsLookup, 0, "DNS Lookup not greater than 0")
				assert.Greater(t, got.AvgTcpConn, 0, "TCP Connection not greater than 0")
				if strings.Contains(tt.args.target, "https://") {
					assert.Greater(t, got.AvgTlsHandshake, 0, "TLS Handshake not greater than 0")
				}
				assert.Greater(t, got.AvgServerProcessing, 0, "Server Processing not greater than 0")
				assert.Greater(t, got.AvgContentTransfer, 0, "Content Transfer not greater than 0")
				assert.Greater(t, got.AvgTotal, 0, "Total not greater than 0")
				return
			}

		})
	}
}
