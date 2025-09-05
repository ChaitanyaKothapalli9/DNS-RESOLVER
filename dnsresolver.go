package main

import (
    "fmt"
    "github.com/miekg/dns"  
    "time"
    "math/rand"
    "net"
    "strings"
)

var rootServers = []string{
    "198.41.0.4", "199.9.14.201", "192.33.4.12", 
}

var dnsCache = map[string]cachedRecord{}

type cachedRecord struct {
    records []string
    expiry  time.Time
}

func resolveCNAME(name string) []string {
    cnames, err := queryDNS(name, dns.TypeCNAME)
    if err != nil || len(cnames) == 0 {
        return []string{fmt.Sprintf("%s is not a CNAME", name)}
    }
    var canonicalResults []string
    for _, cname := range cnames {
        aRecords, err := queryDNS(cname, dns.TypeA)
        if err == nil && len(aRecords) > 0 {
            canonicalResults = append(canonicalResults, fmt.Sprintf("%s,%s", cname, aRecords[0]))
        } else {
            canonicalResults = append(canonicalResults, fmt.Sprintf("%s has no A records", cname))
        }
    }
    return canonicalResults
}

func extractRecords(msg *dns.Msg) []string {
    var records []string
    for _, ans := range msg.Answer {
        records = append(records, ans.String())
    }
    return records
}

func recursiveQuery(domain string, qtype uint16, servers []string) ([]string, error) {
    var lastResponse *dns.Msg
    for _, server := range servers {
        serverAddress := fmt.Sprintf("%s:53", server)
        m := new(dns.Msg)
        m.SetQuestion(dns.Fqdn(domain), qtype)
        m.RecursionDesired = false
        query, _ := m.Pack()
        response, err := sendDNSQuery(query, serverAddress)
        if err != nil {
            continue
        }
        lastResponse = response
        if len(response.Answer) > 0 && response.Authoritative {
            return extractRecords(response), nil
        }
        if len(response.Ns) > 0 {
            var newServers []string
            for _, rr := range response.Ns {
                if ns, ok := rr.(*dns.NS); ok {
                    newServers = append(newServers, ns.Ns)
                }
            }
            if len(newServers) > 0 {
                return recursiveQuery(domain, qtype, newServers)
            }
        }
    }
    if lastResponse != nil {
        return extractRecords(lastResponse), nil
    }
    return nil, fmt.Errorf("resolution failed")
}

func resolve(name string, t RecordType) []string {
    cached, found := dnsCache[name]
    if found && time.Now().Before(cached.expiry) {
        return cached.records
    }
    results, err := recursiveQuery(name, uint16(t), rootServers)
    if err != nil {
        fmt.Println("Error resolving name:", err)
        return []string{"No record found"}
    }
    dnsCache[name] = cachedRecord{results, time.Now().Add(300 * time.Second)}
    return results
}

func queryDNS(domain string, qtype uint16) ([]string, error) {
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(domain), qtype)
    m.RecursionDesired = true
    query, err := m.Pack()
    if err != nil {
        return nil, fmt.Errorf("Failed to pack query: %v", err)
    }
    response, err := sendDNSQuery(query, rootServers[0]+":53")
    if err != nil {
        return nil, err
    }
    var records []string
    for _, ans := range response.Answer {
        records = append(records, ans.String())
    }
    return records, nil
}

func resolveAAAA(name string) []string {
    responses, err := queryDNS(name, dns.TypeAAAA)
    if err != nil || len(responses) == 0 {
        return []string{}
    }
    for _, response := range responses {
        ip := net.ParseIP(strings.TrimSuffix(response, "."))
        if ip != nil && ip.To4() == nil {
            return []string{ip.String()}
        }
    }
    return []string{}
}

func resolveNS(name string) []string {
    nsRecords, err := queryDNS(name, dns.TypeNS)
    if err != nil || len(nsRecords) == 0 {
        return []string{}
    }
    var results []string
    for _, ns := range nsRecords {
        ipRecords, err := queryDNS(ns, dns.TypeA)
        if err == nil && len(ipRecords) > 0 {
            results = append(results, fmt.Sprintf("%s,%s", ns, ipRecords[0]))
        } else {
            results = append(results, ns)
        }
        break
    }
    return results
}

func queryFromServers(name string, t RecordType, servers []string) ([]string, []string, error) {
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(name), uint16(t))
    m.RecursionDesired = true
    c := new(dns.Client)
    for _, server := range servers {
        r, _, err := c.Exchange(m, server+":53")
        if err != nil {
            continue
        }
        if r.Rcode != dns.RcodeSuccess {
            continue
        }
        var records []string
        for _, ans := range r.Answer {
            records = append(records, ans.String())
        }
        return records, extractServers(r), nil
    }
    return nil, nil, fmt.Errorf("no successful response")
}

func resolveTXT(name string) []string {
    txtRecords, err := queryDNS(name, dns.TypeTXT)
    if err != nil || len(txtRecords) == 0 {
        return []string{"No TXT record found"}
    }
    rand.Seed(time.Now().UnixNano())
    randomNum := rand.Intn(100)
    for i, txt := range txtRecords {
        txtRecords[i] = fmt.Sprintf("%s Random number: %d", txt, randomNum)
    }
    return txtRecords
}

func extractServers(msg *dns.Msg) []string {
    var servers []string
    for _, rr := range msg.Ns {
        if a, ok := rr.(*dns.NS); ok {
            servers = append(servers, a.Ns)
        }
    }
    return servers
}
