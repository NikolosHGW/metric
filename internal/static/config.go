package static

import (
	"fmt"

	"github.com/spf13/viper"
)

type StaticRulesConfig struct {
	RuleSet map[string]struct{}
}

func NewStaticRulesConfig() *StaticRulesConfig {
	viper.SetConfigName("staticchecks")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("ошибка при чтении конфигурационного файла: %w", err))
	}
	cnfg := &StaticRulesConfig{RuleSet: make(map[string]struct{})}
	untypedChecks, ok := viper.Get("checks").([]any)
	if !ok {
		panic(fmt.Errorf("неверное содержимое конфигурации в поле checks"))
	}

	checks := make([]string, 0, len(untypedChecks))
	for _, untypedCheck := range untypedChecks {
		value, ok := untypedCheck.(string)
		if !ok {
			panic("неверное значение конфигурации в поле checks")
		}
		checks = append(checks, value)
	}

	for _, check := range checks {
		cnfg.RuleSet[check] = struct{}{}
	}

	return cnfg
}
