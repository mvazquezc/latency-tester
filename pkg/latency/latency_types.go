package latency

type LatencyTest interface {
	Run() (LatencyTestOutput, error)
}

type LatencyTestOutput struct {
	DnsLookup        int
	TcpConn          int
	TlsHandshake     int
	ServerProcessing int
	ContentTransfer  int
	Total            int
}

type AggregatedLatencyTestOutput struct {
	AvgDnsLookup        int `json:"avg_dns_lookup", yaml:"avg_dns_lookup"`
	AvgTcpConn          int `json:"avg_tcp_conn", yaml:"avg_tcp_conn"`
	AvgTlsHandshake     int `json:"avg_tls_handshake", yaml:"avg_tls_handshake"`
	AvgServerProcessing int `json:"avg_server_processing", yaml:"avg_server_processing"`
	AvgContentTransfer  int `json:"avg_content_transfer", yaml:"avg_content_transfer"`
	AvgTotal            int `json:"avg_total", yaml:"avg_total"`
}

type LatencyHTTPTest struct {
	Url string
}

type LatencyTCPTest struct {
	Socket      string
	SendTCPPing bool
}
