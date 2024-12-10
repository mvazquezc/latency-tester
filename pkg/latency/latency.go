package latency

import (
	"github.com/tcnksm/go-httpstat"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Implements the Run method for the Interface LatencyTest
func RunLatencyTest(lt LatencyTest) (LatencyTestOutput, error) {
	return lt.Run()
}

// Implements the Run method for HTTP/s latency tests
func (lt LatencyHTTPTest) Run() (LatencyTestOutput, error) {
	lto := LatencyTestOutput{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", lt.Url, nil)
	if err != nil {
		return lto, err
	}

	// Create a httpstat powered context
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Send request by default HTTP client
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return lto, err
	}
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return lto, err
	}

	res.Body.Close()
	end := time.Now()
	result.End(end)
	// Close connections
	defer client.CloseIdleConnections()
	// Prepare result
	lto.DnsLookup = int(result.DNSLookup / time.Millisecond)
	lto.TcpConn = int(result.TCPConnection / time.Millisecond)
	lto.TlsHandshake = int(result.TLSHandshake / time.Millisecond)
	lto.ServerProcessing = int(result.ServerProcessing / time.Millisecond)
	lto.ContentTransfer = int(result.ContentTransfer(end) / time.Millisecond)
	lto.Total = int(result.Total(end) / time.Millisecond)
	return lto, nil
}

// Implements the Run method for TCP latency tests
func (lt LatencyTCPTest) Run() (LatencyTestOutput, error) {
	lto := LatencyTestOutput{}
	u, err := url.Parse(lt.Socket)
	if err != nil {
		return lto, err
	}

	timeout := 5 * time.Second

	startDNS := time.Now()

	_, err = net.LookupHost(u.Hostname())
	if err != nil {
		return lto, err
	}
	elapsedDNS := int(time.Since(startDNS) / time.Millisecond)

	startTCP := time.Now()
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return lto, err
	}
	defer conn.Close()

	if lt.SendTCPPing {

		strEcho := "PING\n"
		_, err = conn.Write([]byte(strEcho))
		if err != nil {
			return lto, err
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		reply := make([]byte, 1024)
		_, err = conn.Read(reply)
		if err != nil {
			return lto, err
		}
	}

	elapsedTCP := int(time.Since(startTCP) / time.Millisecond)

	lto.DnsLookup = elapsedDNS
	lto.TcpConn = elapsedTCP

	return lto, nil
}
