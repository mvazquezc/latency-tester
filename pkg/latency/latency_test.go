package latency

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLatencyHTTPTest_Run(t *testing.T) {

	type fields struct {
		Url string
	}
	tests := []struct {
		name         string
		fields       fields
		want         LatencyTestOutput
		wantErr      bool
		wantErrError string
	}{
		{
			name: "Test HTTP URL",
			fields: fields{
				Url: "http://google.com",
			},
			want:    LatencyTestOutput{},
			wantErr: false,
		},
		{
			name: "Test HTTPs URL",
			fields: fields{
				Url: "https://redhat.com",
			},
			want:    LatencyTestOutput{},
			wantErr: false,
		},
		{
			name: "Test non existing DNS",
			fields: fields{
				Url: "http://wefnwejklflwewef.com",
			},
			want:         LatencyTestOutput{},
			wantErr:      true,
			wantErrError: "Get \"http://wefnwejklflwewef.com\": dial tcp: lookup wefnwejklflwewef.com: no such host",
		},
		{
			name: "Test HTTP Wrong URL",
			fields: fields{
				Url: "htt://google.com",
			},
			want:         LatencyTestOutput{},
			wantErr:      true,
			wantErrError: "Get \"htt://google.com\": unsupported protocol scheme \"htt\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := LatencyHTTPTest{
				Url: tt.fields.Url,
			}
			got, err := lt.Run()

			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.wantErrError, "Error not expected")
				return
			}

			if got.DnsLookup <= 0 || got.TcpConn <= 0 || got.TlsHandshake <= 0 || got.ServerProcessing <= 0 || got.ContentTransfer <= 0 || got.Total <= 0 {
				assert.Greater(t, got.DnsLookup, 0, "DNS Lookup not greater than 0")
				assert.Greater(t, got.TcpConn, 0, "TCP Connection not greater than 0")
				if strings.Contains(lt.Url, "https://") {
					assert.Greater(t, got.TlsHandshake, 0, "TLS Handshake not greater than 0")
				}
				assert.Greater(t, got.ServerProcessing, 0, "Server Processing not greater than 0")
				assert.Greater(t, got.ContentTransfer, 0, "Content Transfer not greater than 0")
				assert.Greater(t, got.Total, 0, "Total not greater than 0")
				return
			}

		})
	}
}

func TestLatencyTCPTest_Run(t *testing.T) {
	type fields struct {
		Socket      string
		SendTCPPing bool
	}
	tests := []struct {
		name         string
		fields       fields
		want         LatencyTestOutput
		wantErr      bool
		wantErrError string
	}{
		{
			name: "Test TCP Socket",
			fields: fields{
				Socket: "tcp://google.es:80",
			},
			want:    LatencyTestOutput{},
			wantErr: false,
		},
		{
			name: "Test TCP Socket with TCP Ping",
			fields: fields{
				Socket:      "tcp://redhat.es:80",
				SendTCPPing: true,
			},
			want:    LatencyTestOutput{},
			wantErr: false,
		},
		{
			name: "Test non existing DNS",
			fields: fields{
				Socket: "tcp://wefnwejklflwewef.com:80",
			},
			want:         LatencyTestOutput{},
			wantErr:      true,
			wantErrError: "lookup wefnwejklflwewef.com: no such host",
		},
		{
			name: "Test TCP Socket without port",
			fields: fields{
				Socket: "tcp://google.com",
			},
			want:         LatencyTestOutput{},
			wantErr:      true,
			wantErrError: "dial tcp: address google.com: missing port in address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := LatencyTCPTest{
				Socket:      tt.fields.Socket,
				SendTCPPing: tt.fields.SendTCPPing,
			}
			got, err := lt.Run()

			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.wantErrError, "Error not expected")
				return
			}

			if got.DnsLookup <= 0 || got.TcpConn <= 0 || got.TlsHandshake != 0 || got.ServerProcessing != 0 || got.ContentTransfer != 0 || got.Total != 0 {
				assert.Greater(t, got.DnsLookup, 0, "DNS Lookup not greater than 0")
				assert.Greater(t, got.TcpConn, 0, "TCP Connection not greater than 0")
				assert.Equal(t, got.TlsHandshake, 0, "TLS Handshake is not 0")
				assert.Equal(t, got.ServerProcessing, 0, "Server Processing is not 0")
				assert.Equal(t, got.ContentTransfer, 0, "Content Transfer is not 0")
				assert.Equal(t, got.Total, 0, "Total is not 0")
				return
			}
		})
	}
}
