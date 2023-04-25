package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mvazquezc/latency-tester/pkg/latency"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

func IntervalTimeToSeconds(interval string) int {
	r := regexp.MustCompile(`^(\d*)(s|m|h)`)
	captureGroups := r.FindStringSubmatch(interval)
	if len(captureGroups) < 1 {
		return -1
	}
	timeValue, err := strconv.Atoi(captureGroups[1])
	if err != nil {
		return -1
	}
	timeUnit := captureGroups[2]

	switch timeUnit {
	case "s":
		return timeValue
	case "m":
		return timeValue * 60
	case "h":
		return timeValue * 3600
	default:
		return 1
	}
}

func ValidateIntervalTime(interval string) bool {
	r := regexp.MustCompile(`^(\d*)(s|m|h)`)
	matched := r.MatchString(interval)
	return matched
}

func GetScheme(inputTarget string) string {
	u, _ := url.Parse(inputTarget)
	if u.Scheme == "tcp" {
		return "tcp"
	}
	return "http"
}

func ValidateTarget(inputTarget string) (bool, string) {
	u, err := url.Parse(inputTarget)
	if err != nil {
		return false, ""
	}
	if (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "tcp") && u.Host != "" {
		return true, u.Scheme
	}
	return false, ""
	//return err == nil && u.Scheme != "" && u.Host != ""
}

func WriteOutputTable(lto latency.AggregatedLatencyTestOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{color.InWhite("Test"), color.InWhite("Average Value")})
	t.AppendRow([]interface{}{color.InWhite("Average DNS Lookup"), color.InWhite(fmt.Sprintf("%d ms", lto.AvgDnsLookup))})
	t.AppendRow([]interface{}{color.InWhite("Average TCP Connection"), color.InWhite(fmt.Sprintf("%d ms", lto.AvgTcpConn))})
	t.AppendRow([]interface{}{color.InWhite("Average TLS Handshake"), color.InWhite(fmt.Sprintf("%d ms", lto.AvgTlsHandshake))})
	t.AppendRow([]interface{}{color.InWhite("Average Server Processing"), color.InWhite(fmt.Sprintf("%d ms", lto.AvgServerProcessing))})
	t.AppendRow([]interface{}{color.InWhite("Average Content Transfer"), color.InWhite(fmt.Sprintf("%d ms", lto.AvgContentTransfer))})
	t.AppendRow([]interface{}{color.InWhite("Average Total"), color.InWhite(fmt.Sprintf("%d ms", lto.AvgTotal))})
	t.SetStyle(table.StyleLight)
	t.Render()
}

func WriteOutputJson(lto latency.AggregatedLatencyTestOutput) {
	o, _ := json.MarshalIndent(lto, "", "    ")
	fmt.Println(string(o))
}

func WriteOutputYaml(lto latency.AggregatedLatencyTestOutput) {
	o, _ := yaml.Marshal(lto)
	fmt.Println(string(o))
}
