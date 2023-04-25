package commands

import (
	"fmt"
	"github.com/mvazquezc/latency-tester/pkg/latency"
	"github.com/mvazquezc/latency-tester/pkg/utils"
	"log"
	"time"
)

func RunLatencyTestCmd(target string, numberOfRuns int, waitIntervalSeconds int, tcpPing bool) (latency.AggregatedLatencyTestOutput, error) {

	var (
		dnsLookup        int
		tcpConn          int
		tlsHandshake     int
		serverProcessing int
		contentTransfer  int
		total            int
		lt               latency.LatencyTest
	)
	fmt.Printf("Run command executed with target: %s, numberOfRuns: %d, waitInternvalSeconds: %d, tcpPing: %v.\n", target, numberOfRuns, waitIntervalSeconds, tcpPing)
	scheme := utils.GetScheme(target)

	if scheme == "tcp" {
		lt = latency.LatencyTCPTest{
			Socket:      target,
			SendTCPPing: tcpPing,
		}
	} else {
		lt = latency.LatencyHTTPTest{
			Url: target,
		}
	}

	for i := 1; i <= numberOfRuns; i++ {
		log.Printf("Request number [%d/%d]", i, numberOfRuns)

		lto, err := latency.RunLatencyTest(lt)

		if err != nil {
			return latency.AggregatedLatencyTestOutput{}, err
		}

		// Aggregate results
		dnsLookup += lto.DnsLookup
		tcpConn += lto.TcpConn
		tlsHandshake += lto.TlsHandshake
		serverProcessing += lto.ServerProcessing
		contentTransfer += lto.ContentTransfer
		total += lto.Total

		if numberOfRuns > 1 && i != numberOfRuns {
			time.Sleep(time.Duration(waitIntervalSeconds) * time.Second)
		}

	}
	lto := latency.AggregatedLatencyTestOutput{
		AvgDnsLookup:        dnsLookup / numberOfRuns,
		AvgTcpConn:          tcpConn / numberOfRuns,
		AvgTlsHandshake:     tlsHandshake / numberOfRuns,
		AvgServerProcessing: serverProcessing / numberOfRuns,
		AvgContentTransfer:  contentTransfer / numberOfRuns,
		AvgTotal:            total / numberOfRuns,
	}

	return lto, nil
}
