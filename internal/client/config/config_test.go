package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type netAddressStringTestCase struct {
	name string
	na   *netAddress
	want string
}

var netAddressStringTests = []netAddressStringTestCase{
	{"Localhost", &netAddress{"localhost", 8080}, "localhost:8080"},
	{"Example", &netAddress{"example.com", 1234}, "example.com:1234"},
	{"Empty", &netAddress{"", 0}, ":0"},
}

func TestNetAddressString(t *testing.T) {
	for _, tc := range netAddressStringTests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.na.String())
		})
	}
}

type netAddressSetTestCase struct {
	name      string
	flagValue string
	wantErr   bool
	wantHost  string
	wantPort  int
}

var netAddressSetTests = []netAddressSetTestCase{
	{"отрицательный тест: неправильный формат адреса", "invalid", true, "", 0},
	{"положительный тест: домен с портом", "example.com:1234", false, "example.com", 1234},
	{"положительный тест: localhost", "localhost:8080", false, "localhost", 8080},
	{"отрицательный тест: пустое значение", "", true, "", 0},
}

func TestNetAddressSet(t *testing.T) {
	for _, tc := range netAddressSetTests {
		t.Run(tc.name, func(t *testing.T) {
			na := &netAddress{}

			err := na.Set(tc.flagValue)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.wantHost, na.Host)
			assert.Equal(t, tc.wantPort, na.Port)
		})
	}
}
