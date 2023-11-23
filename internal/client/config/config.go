package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env"
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
		return errors.New("internal/client/config netAdress_Set: need address in a form host:port")
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
	Endpoint       netAddress
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Address        string        `env:"ADDRESS"`
}

func (c *config) GetEndpointObject() flag.Value {
	return &c.Endpoint
}

func (c *config) GetPollIntervalPointer() *time.Duration {
	return &c.PollInterval
}

func (c *config) GetReportIntervalPointer() *time.Duration {
	return &c.ReportInterval
}

func (c *config) InitAdress() {
	c.Address = c.Endpoint.String()
}

func (c *config) InitEnv() {
	env.Parse(c)
}

func NewConfig() *config {
	return &config{
		Endpoint: netAddress{Host: "localhost", Port: 8080},
	}
}
