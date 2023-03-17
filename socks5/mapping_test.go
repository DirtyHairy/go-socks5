package socks5

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMapping2_Valid(t *testing.T) {
	mapping, err := MappingFromSpec("1.2.3.4:10.11.12.13")

	assert.Nil(t, err)
	assert.Equal(t, mapping.addressFrom, net.IPv4(1, 2, 3, 4))
	assert.Equal(t, mapping.portFrom, 0)
	assert.Equal(t, mapping.addressTo, net.IPv4(10, 11, 12, 13))
	assert.Equal(t, mapping.portTo, 0)
}

func TestParseMapping4_Valid(t *testing.T) {
	mapping, err := MappingFromSpec("1.2.3.4:23:10.11.12.13:48")

	assert.Nil(t, err)
	assert.Equal(t, mapping.addressFrom, net.IPv4(1, 2, 3, 4))
	assert.Equal(t, mapping.portFrom, 23)
	assert.Equal(t, mapping.addressTo, net.IPv4(10, 11, 12, 13))
	assert.Equal(t, mapping.portTo, 48)
}

func TestParseMapping3_Invalid(t *testing.T) {
	_, err := MappingFromSpec("1.2.3.4:23:10.11.12.13")
	assert.NotNil(t, err)
}

func TestParseMappingBadIP_Invalid(t *testing.T) {
	_, err := MappingFromSpec("1.2.3.4aa:23:10.11.12.13:48")
	assert.NotNil(t, err)
}

func TestParseMappingBadPort_Invalid(t *testing.T) {
	_, err := MappingFromSpec("1.2.3.4:23:10.11.12.13:aa")
	assert.NotNil(t, err)
}

func TestParseMappingsGood(t *testing.T) {
	mappings, err := MappingsFromSpecs([]string{"1.2.3.4:10.11.12.13", "1.2.3.4:23:10.11.12.13:48"})

	assert.Nil(t, err)
	assert.EqualValues(t, mappings, []Mapping{
		Mapping{
			addressFrom: net.IPv4(1, 2, 3, 4),
			portFrom:    0,
			addressTo:   net.IPv4(10, 11, 12, 13),
			portTo:      0,
		},
		Mapping{
			addressFrom: net.IPv4(1, 2, 3, 4),
			portFrom:    23,
			addressTo:   net.IPv4(10, 11, 12, 13),
			portTo:      48,
		},
	})
}

func TestParseMappingsBad(t *testing.T) {
	_, err := MappingsFromSpecs([]string{"1.2.3.4:10.11.12.13", "1.2.3.4:23:10.11.12.13:aa"})
	assert.NotNil(t, err)
}

func TestMapsUnspecificPort(t *testing.T) {
	mappings, _ := MappingsFromSpecs([]string{"1.2.3.4:10.11.12.13", "4.5.6.7:23:10.11.12.14:48"})

	mapped, mappingApplied := mappings.Map(AddrSpec{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 66,
	})

	assert.True(t, mappingApplied)
	assert.Equal(t, mapped, AddrSpec{
		IP:   net.IPv4(10, 11, 12, 13),
		Port: 66,
	})
}

func TestMapsSpecificPort(t *testing.T) {
	mappings, _ := MappingsFromSpecs([]string{"1.2.3.4:10.11.12.13", "4.5.6.7:23:10.11.12.14:48"})

	mapped, mappingApplied := mappings.Map(AddrSpec{
		IP:   net.IPv4(4, 5, 6, 7),
		Port: 23,
	})

	assert.True(t, mappingApplied)
	assert.Equal(t, mapped, AddrSpec{
		IP:   net.IPv4(10, 11, 12, 14),
		Port: 48,
	})
}
