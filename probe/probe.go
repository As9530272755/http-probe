package probe

import (
	"fmt"
	"http-probe/config"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

func DoHttpProbe(url string) string {
	twSec := config.GlobalTwSec
	dnsStr := ""
	targetAddr := ""
	var t0, t1, t2, t3, t4 time.Time
	start := time.Now()

	req, _ := http.NewRequest("GET", url, nil)

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t0 = time.Now()
		},

		DNSDone: func(di httptrace.DNSDoneInfo) {
			t1 = time.Now()
			ips := make([]string, 0)
			for _, d := range di.Addrs {
				ips = append(ips, d.IP.String())
			}
			dnsStr = strings.Join(ips, ",")
		},

		ConnectStart: func(network, addr string) {
			if t1.IsZero() {
				t1 = time.Now()
			}
		},

		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("[无法建立连接][addr:%v][err:%v]", addr, err)
				return
			}

			targetAddr = addr
			t2 = time.Now()
		},

		GotConn: func(_ httptrace.GotConnInfo) {
			t3 = time.Now()
		},

		GotFirstResponseByte: func() {
			t4 = time.Now()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := http.Client{
		Timeout: time.Duration(twSec) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("[http 探测出错]\n"+
			"[http 探测的目标：%s]\n"+
			"[错误详情：%v]\n"+
			"[总耗时：%s]\n",
			url,
			err,
			msDurationStr(time.Now().Sub(start)),
		)
		log.Printf(msg)
		return msg
	}
	defer resp.Body.Close()

	end := time.Now()

	if t0.IsZero() {
		t0 = t1
	}

	dnsLookup := msDurationStr(t1.Sub(t0))
	tcpConnection := msDurationStr(t3.Sub(t1))
	serverProcessing := msDurationStr(t4.Sub(t3))
	totoal := msDurationStr(end.Sub(t0))
	probeResStr := fmt.Sprintf(
		"[http 探测目标:%s]\n"+
			"[dns 解析结果:%s]\n"+
			"[链接 ip and port:%s]\n"+
			"[状态码:%d]\n"+
			"[dns 解析耗时:%s]\n"+
			"[tcp 解析耗时:%s]\n"+
			"[服务器处理耗时:%s]\n"+
			"[总耗时:%s]",
		url,
		dnsStr,
		targetAddr,
		resp.StatusCode,
		dnsLookup,
		tcpConnection,
		serverProcessing,
		totoal,
	)
	return probeResStr
}

func msDurationStr(d time.Duration) string {
	return fmt.Sprintf("%dms", int(d/time.Millisecond))
}
