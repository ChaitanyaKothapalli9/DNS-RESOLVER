package main

import (
    "flag"
    "fmt"
    "os"
)

type RecordType uint16

const (
    TYPE_A     RecordType = 1
    TYPE_NS    RecordType = 2
    TYPE_CNAME RecordType = 5
    TYPE_TXT   RecordType = 16
    TYPE_AAAA  RecordType = 28
)

var RecordTypes = map[string]RecordType{
    "A":     TYPE_A,
    "NS":    TYPE_NS,
    "CNAME": TYPE_CNAME,
    "TXT":   TYPE_TXT,
    "AAAA":  TYPE_AAAA,
}

func main() {
    t := flag.String("type", "A", "the record type to query for each name")
    flag.Parse()
    names := flag.Args()

    if len(names) == 0 {
        fmt.Println("Not enough arguments, must pass in at least one name")
        os.Exit(1)
    }

    recordType, exists := RecordTypes[*t]
    if !exists {
        fmt.Printf("Specified record type %s doesn't exist. Must be one of %v\n", *t, RecordTypes)
        os.Exit(1)
    }

    for _, name := range names {
        results := resolve(name, recordType)
        if len(results) == 0 {
            fmt.Printf("%s, No record found\n", name)
        } else {
            for _, result := range results {
                fmt.Println(result)
            }
        }
    }
}
