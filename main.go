package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dirtyhairy/socks5-server/socks5"
)

type SpecList []string

func (s SpecList) String() string {
	return strings.Join(s, ",")
}

func (s *SpecList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func Usage() {
	fmt.Println(`usage: socks5-server [ -m mapping ] [ -m mapping ] ...

Valid destination mappings:

 * <source address>:<source port>:<target address>:<target port>
 * <source address>:<target address>

Valid parameters:
`)

	flag.PrintDefaults()
	fmt.Println("  -h    show this help")
}

func main() {
	var specs SpecList
	var err error

	flag.Usage = Usage
	flag.Var(&specs, "m", "specify a mapping")
	flag.Parse()

	mappings, err := socks5.MappingsFromSpecs(specs)
	if err != nil {
		fmt.Printf("failed to parse mappings: %v\n", err)
		os.Exit(1)
	}

	conf := &socks5.Config{Rewriter: mappings}
	server, err := socks5.New(conf)
	if err != nil {
		fmt.Printf("ERROR: unable to create server: %v\n", err)
		os.Exit(1)
	}

	resultChannel := make(chan int, 1)

	go func() {
		if err := server.ListenAndServe("tcp", "127.0.0.1:9998"); err != nil {
			fmt.Printf("ERROR: unable to start server: %v\n", err)
			os.Exit(1)

			resultChannel <- 1
		}

		resultChannel <- 0
	}()

	fmt.Println("server listening...")

	<-resultChannel
}
