package connection

import (
	"fmt"

	"github.com/project-flogo/core/app/resolve"
)

type Config struct {
	Ref      string                 `json:"ref,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

func ToConfig(config map[string]interface{}) (*Config, error) {

	if v, ok := config["ref"]; ok {
		if ref, ok := v.(string); ok {
			cfg := &Config{}
			cfg.Ref = ref
			if v, ok := config["settings"]; ok {
				if settings, ok := v.(map[string]interface{}); ok {
					// Resolve property/env value
					for name, value := range settings {
						strVal, ok := value.(string)
						if ok && len(strVal) > 0 && strVal[0] == '=' {
							var err error
							value, err = resolve.Resolve(strVal[1:], nil)
							if err != nil {
								return nil, err
							}
						}
						cfg.Settings[name] = value
					}
					return cfg, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("invalid connection config: %+v", config)
}

func ResolveConfig(config *Config) error {

	for name, value := range config.Settings {

		if strVal, ok := value.(string); ok && len(strVal) > 0 && strVal[0] == '=' {
			var err error
			value, err = resolve.Resolve(strVal[1:], nil)
			if err != nil {
				return err
			}

			config.Settings[name] = value
		}
	}
	return nil
}

