package main

import (
    "net"
    "time"
    "github.com/miekg/dns"
    "fmt"
)

func sendDNSQuery(query []byte, serverAddress string) (*dns.Msg, error) {
    conn, err := net.Dial("udp", serverAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to server %s: %v", serverAddress, err)
    }
    defer conn.Close()

    _, err = conn.Write(query)
    if err != nil {
        return nil, fmt.Errorf("failed to send query to server %s: %v", serverAddress, err)
    }

    buffer := make([]byte, 1024)
    readTimeout := 2 * time.Second
    conn.SetReadDeadline(time.Now().Add(readTimeout))
    n, err := conn.Read(buffer)
    if err != nil {
        return nil, fmt.Errorf("failed to read from server %s: %v", serverAddress, err)
    }

    msg := new(dns.Msg)
    if err := msg.Unpack(buffer[:n]); err != nil {
        return nil, fmt.Errorf("failed to unpack DNS response from server %s: %v", serverAddress, err)
    }

    return msg, nil
}
