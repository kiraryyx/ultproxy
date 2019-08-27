package main

import (
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/comail/colog"
	transproxy "github.com/kiraryyx/ultproxy"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func initProxies() {
	// seed the global random number generator, used in secureoperator
	rand.Seed(time.Now().UTC().UnixNano())

	// setup logger
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LInfo)
	level, err := colog.ParseLevel(config.Logging.Level)
	if err != nil {
		log.Fatalf("alert: Invalid log level: %s", err.Error())
	}
	colog.SetMinLevel(level)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lmicroseconds,
	})
	colog.ParseFields(true)
	colog.Register()

	if config.LiveServer.StartExplicitOnly {
		startExplicitProxyOnly(level)
	} else {
		startAllProxy(level)
	}
}

func startExplicitProxyOnly(level colog.Level) {
	startExplicitProxy()

	// serve until exit
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Printf("info: Proxy servers stopping.")
	log.Printf("info: ultproxy exited.")
}

func startAllProxy(level colog.Level) {
	var err error

	// Decide upstream proxy for each proxy
	upstreamStrForTCP := config.LiveServer.TCP.UpstreamProxyURL
	if upstreamStrForTCP == "" {
		upstreamStrForTCP = config.LiveServer.UpstreamProxyURL
	}

	upstreamURLForTCP, err := url.Parse(upstreamStrForTCP)
	if err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	upstreamStrForHTTP := config.LiveServer.HTTP.UpstreamProxyURL
	if upstreamStrForHTTP == "" {
		upstreamStrForHTTP = config.LiveServer.UpstreamProxyURL
	}

	upstreamURLForHTTP, err := url.Parse(upstreamStrForHTTP)
	if err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	upstreamStrForHTTPS := config.LiveServer.HTTPS.UpstreamProxyURL
	if upstreamStrForHTTPS == "" {
		upstreamStrForHTTPS = config.LiveServer.UpstreamProxyURL
	}

	upstreamURLForHTTPS, err := url.Parse(upstreamStrForHTTPS)
	if err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	np := parseNoProxy(config.LiveServer.NoProxy)
	// start servers
	tcpProxy := transproxy.NewTCPProxy(
		transproxy.TCPProxyConfig{
			ListenAddress:    config.LiveServer.TCP.ListenAddress,
			NoProxy:          np,
			UpstreamProxyURL: upstreamURLForTCP,
		},
	)
	if err := tcpProxy.Start(); err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	httpProxy := transproxy.NewHTTPProxy(
		transproxy.HTTPProxyConfig{
			ListenAddress:    config.LiveServer.HTTP.ListenAddress,
			NoProxy:          np,
			UpstreamProxyURL: upstreamURLForHTTP,
			Verbose:          level == colog.LDebug,
		},
	)
	if err := httpProxy.Start(); err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	httpsProxy := transproxy.NewHTTPSProxy(
		transproxy.HTTPSProxyConfig{
			ListenAddress:    config.LiveServer.HTTPS.ListenAddress,
			NoProxy:          np,
			UpstreamProxyURL: upstreamURLForHTTPS,
		},
	)
	if err := httpsProxy.Start(); err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	startExplicitProxy()

	log.Printf("info: All proxy servers started.")

	httpToPort := toPort(config.LiveServer.HTTP.ListenAddress)
	httpsToPort := toPort(config.LiveServer.HTTPS.ListenAddress)
	tcpToPort := toPort(config.LiveServer.TCP.ListenAddress)
	tcpDPorts := toPorts(config.LiveServer.TCP.DestPorts)

	var t *transproxy.IPTables

	if config.LiveServer.Iptables.EnableAutoConfig {
		t, err = transproxy.NewIPTables(&transproxy.IPTablesConfig{
			HTTPToPort:  httpToPort,
			HTTPSToPort: httpsToPort,
			TCPToPort:   tcpToPort,
			TCPDPorts:   tcpDPorts,
		})
		if err != nil {
			log.Printf("alert: %s", err.Error())
		}

		t.Start()

		log.Printf(`info: iptables rules inserted as follows.
---
%s
---`, t.Show())
	}

	// serve until exit
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Printf("info: Proxy servers stopping.")

	// start shutdown process
	if config.LiveServer.Iptables.EnableAutoConfig {
		t.Stop()
		log.Printf("info: iptables rules deleted.")
	}

	log.Printf("info: ultproxy exited.")
}

func startExplicitProxy() {
	upstreamStrForExplicit := config.LiveServer.Explicit.UpstreamProxyURL
	if upstreamStrForExplicit == "" {
		upstreamStrForExplicit = config.LiveServer.UpstreamProxyURL
	}

	upstreamURLForExplicit, err := url.Parse(upstreamStrForExplicit)
	if err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	upstreamStrForExplicitWithAuth := config.LiveServer.Explicit.UpstreamProxyURL
	if upstreamStrForExplicitWithAuth == "" {
		upstreamStrForExplicitWithAuth = upstreamStrForExplicit
	}

	upstreamURLForExplicitWithAuth, err := url.Parse(upstreamStrForExplicitWithAuth)
	if err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	explicitProxyWithAuth := transproxy.NewExplicitProxy(
		transproxy.ExplicitProxyConfig{
			ListenAddress:         config.LiveServer.ExplicitWithAuth.ListenAddress,
			UpstreamProxyURL:      upstreamURLForExplicitWithAuth,
			UseProxyAuthorization: true,
		},
	)
	if err := explicitProxyWithAuth.Start(); err != nil {
		log.Fatalf("alert: %s", err.Error())
	}

	explicitProxy := transproxy.NewExplicitProxy(
		transproxy.ExplicitProxyConfig{
			ListenAddress:         config.LiveServer.Explicit.ListenAddress,
			UpstreamProxyURL:      upstreamURLForExplicit,
			UseProxyAuthorization: false,
		},
	)
	if err := explicitProxy.Start(); err != nil {
		log.Fatalf("alert: %s", err.Error())
	}
}

func toPort(addr string) int {
	array := strings.Split(addr, ":")
	if len(array) != 2 {
		log.Printf("alert: Invalid address, no port: %s", addr)
	}

	i, err := strconv.Atoi(array[1])
	if err != nil {
		log.Printf("alert: Invalid address, the port isn't number: %s", addr)
	}

	if i > 65535 || i < 0 {
		log.Printf("alert: Invalid address, the port must be an integer value in the range 0-65535: %s", addr)
	}

	return i
}

func toPorts(ports []int) []int {
	var p []int

	for _, i := range ports {
		if i > 65535 || i < 0 {
			log.Printf("alert: Invalid port, It must be an integer value in the range 0-65535: %d", ports)
		}

		p = append(p, i)
	}

	return p
}

func parseNoProxy(noProxy []string) transproxy.NoProxy {
	var ipArray []string
	var cidrArray []*net.IPNet
	var domainArray []string

	for _, v := range noProxy {
		ip := net.ParseIP(v)
		if ip != nil {
			ipArray = append(ipArray, v)
			continue
		}

		_, ipnet, err := net.ParseCIDR(v)
		if err == nil {
			cidrArray = append(cidrArray, ipnet)
			continue
		}

		domainArray = append(domainArray, v)
	}

	return transproxy.NoProxy{
		IPs:     ipArray,
		CIDRs:   cidrArray,
		Domains: domainArray,
	}
}
