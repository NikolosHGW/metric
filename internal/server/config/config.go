package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type netAddress struct {
	Host string
	Port int
}

func (na *netAddress) String() string {
	return fmt.Sprintf("%v:%v", na.Host, na.Port)
}

func (na *netAddress) Set(flagValue string) error {
	hp := strings.Split(flagValue, ":")
	if len(hp) != 2 {
		return errors.New("internal/server/config netAdress_Set: need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	na.Host = hp[0]
	na.Port = port

	return nil
}

type config struct {
	Endpoint netAddress
}

func NewConfig() *config {
	return &config{
		Endpoint: netAddress{Host: "localhost", Port: 8080},
	}
}
