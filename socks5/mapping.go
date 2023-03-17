package socks5

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

type Mapping struct {
	addressFrom net.IP
	portFrom    int

	addressTo net.IP
	portTo    int
}

type MappingSet []Mapping

func MappingFromSpec(spec string) (mapping Mapping, err error) {
	atoms := strings.Split(spec, ":")

	if len(atoms) == 2 {
		mapping.addressFrom = net.ParseIP(atoms[0])
		mapping.addressTo = net.ParseIP(atoms[1])
	}

	if len(atoms) == 4 {
		mapping.addressFrom = net.ParseIP(atoms[0])
		mapping.addressTo = net.ParseIP(atoms[2])

		var e error

		mapping.portFrom, e = strconv.Atoi(atoms[1])
		if e != nil {
			err = errors.New("unable to parse source port")
			return
		}

		mapping.portTo, e = strconv.Atoi(atoms[3])
		if e != nil {
			err = errors.New("unable to parse destination port")
			return
		}
	}

	if mapping.addressFrom == nil {
		err = errors.New("unable to parse source address")
		return
	}

	if mapping.addressTo == nil {
		err = errors.New("unable to parse destination address")
		return
	}

	return
}

func (s Mapping) Apply(incoming AddrSpec) (spec AddrSpec, mappingApplied bool) {
	mappingApplied = false
	spec = incoming

	if !s.addressFrom.Equal(incoming.IP) {
		return
	}

	if s.portFrom != 0 && s.portFrom != incoming.Port {
		return
	}

	spec.IP = s.addressTo
	if s.portTo != 0 {
		spec.Port = s.portTo
	}

	mappingApplied = true

	return
}

func (s MappingSet) Apply(incoming AddrSpec) (spec AddrSpec, mappingApplied bool) {
	mappingApplied = false
	spec = incoming

	if incoming.IP == nil {
		return
	}

	for _, mapping := range s {
		spec, mappingApplied = mapping.Apply(incoming)
		if mappingApplied {
			return
		}
	}

	return
}

func MappingsFromSpecs(specs []string) (mappings MappingSet, err error) {
	mappings = make([]Mapping, len(specs))

	for i, spec := range specs {
		var e error

		mappings[i], e = MappingFromSpec(spec)
		if e != nil {
			err = errors.New(fmt.Sprintf("failed to parse mapping %v: %v", spec, e))
			return
		}
	}

	return
}

func (s MappingSet) Rewrite(ctx context.Context, request *Request, logger *log.Logger) (context.Context, *AddrSpec) {
	mapped, mappingApplied := s.Apply(*request.DestAddr)

	if mappingApplied {
		logger.Printf("mapped %v to %v\n", *request.DestAddr, mapped)
	}

	return ctx, &mapped
}
