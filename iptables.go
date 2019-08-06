package ultproxy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/coreos/go-iptables/iptables"
)

const (
	NAT        = "nat"
	PREROUTING = "PREROUTING"
	OUTPUT     = "OUTPUT"
)

type IPTables struct {
	iptables      *iptables.IPTables
	httpRule      []string
	httpsRule     []string
	tcpRule       []string
	err           error
}

type IPTablesConfig struct {
	HTTPToPort  int
	HTTPSToPort int
	TCPToPort   int
	TCPDPorts   []int
}

func NewIPTables(c *IPTablesConfig) (*IPTables, error) {
	t, err := iptables.New()
	if err != nil {
		return nil, err
	}

	var tcpDPorts []string
	for _, v := range c.TCPDPorts {
		tcpDPorts = append(tcpDPorts, strconv.Itoa(v))
	}

	httpRule := []string{NAT, PREROUTING, "-p", "tcp", "--dport", "80", "-j", "REDIRECT", "--to-ports", strconv.Itoa(c.HTTPToPort)}
	httpsRule := []string{NAT, PREROUTING, "-p", "tcp", "--dport", "443", "-j", "REDIRECT", "--to-ports", strconv.Itoa(c.HTTPSToPort)}
	tcpRule := []string{NAT, PREROUTING, "-p", "tcp", "-m", "multiport", "--dport", strings.Join(tcpDPorts, ","), "-j", "REDIRECT", "--to-ports", strconv.Itoa(c.TCPToPort)}

	return &IPTables{
		iptables:      t,
		httpRule:      httpRule,
		httpsRule:     httpsRule,
		tcpRule:       tcpRule,
	}, nil
}

func (t *IPTables) Start() error {
	t.Check(t.httpRule)
	t.Check(t.httpsRule)
	t.Check(t.tcpRule)

	t.insertRule(t.httpRule)
	t.insertRule(t.httpsRule)
	t.insertRule(t.tcpRule)

	return t.err
}

func (t *IPTables) Stop() error {
	t.deleteRule(t.httpRule)
	t.deleteRule(t.httpsRule)
	t.deleteRule(t.tcpRule)

	return t.err
}

func (t *IPTables) Show() string {
	s := fmt.Sprintf(`iptables -t %s -I %s
iptables -t %s -I %s
iptables -t %s -I %s`,
		t.tcpRule[0], strings.Join(t.tcpRule[1:], " "),
		t.httpsRule[0], strings.Join(t.httpsRule[1:], " "),
		t.httpRule[0], strings.Join(t.httpRule[1:], " "),
	)

	return s
}

func (t *IPTables) Check(rule []string) {
	if t.err != nil || len(rule) < 3 {
		return
	}

	exists, err := t.iptables.Exists(rule[0], rule[1], rule[2:]...)
	if exists {
		t.err = fmt.Errorf("Same iptables rule already exists : iptables -t %s -I %s", rule[0], strings.Join(rule[1:], " "))
	}

	if err != nil {
		t.err = fmt.Errorf("Checking iptables rule failed : %s", err.Error())
	}
}

func (t *IPTables) insertRule(rule []string) {
	if t.err != nil || len(rule) < 3 {
		return
	}

	if err := t.iptables.Insert(rule[0], rule[1], 1, rule[2:]...); err != nil {
		t.err = fmt.Errorf("Insert iptables rule failed : %s", err.Error())
	}
}

func (t *IPTables) deleteRule(rule []string) {
	// Don't skip when it has error for deleting all rules
	if len(rule) < 3 {
		return
	}

	if err := t.iptables.Delete(rule[0], rule[1], rule[2:]...); err != nil {
		t.err = fmt.Errorf("Delete iptables rule failed : %s", err.Error())
	}
}
