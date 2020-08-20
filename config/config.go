package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

// the default file name of the configuration file
const DefaultConfigName = "sparta.yml"

// defaultConfig creates a default configuration where any values that
//               need to default to "non-empty" values (like booleans)
//               or strings can be set and so that the maps and arrays
//               are already constructed to prevent issues with unsafe
//               reads of the configuration
func defaultConfig() SpartaConfig {
	return SpartaConfig{
		OpenShift: OpenShift{},
		Cluster: Cluster{},
		Cloud: Cloud{},
		Subnets: Subnets{},
		ProviderAuth: ProviderAuth{
			Keys: true,
		},
		RedSord: RedSord{},
		Koffer: Koffer{
			Plugins: Plugins{},
		},
	}
}

// ViperSpartaConfig takes a viper instance, configures it for Sparta configuration using the given paths and search
//                   location and loads the values into the viper instance as well as returning the configuration
func ViperSpartaConfig(viperInstance *viper.Viper, fileName string, searchPaths ...string) (*SpartaConfig, error) {
	// break config name
	name := filepath.Base(fileName)
	ext := filepath.Ext(fileName)
	if strings.Index(ext, ".") == 0 {
		ext = ext[1:]
	}
	if len(ext) < 1 {
		ext = "yml"
	}

	// set config file
	viperInstance.SetConfigName(name)
	viperInstance.SetConfigType(ext)

	// add locations to search for the file with that name
	for _, location := range searchPaths {
		viperInstance.AddConfigPath(location)
	}

	// create new default config instance so that
	// parts can be defaulted before unmarshalling
	config := defaultConfig()

	// unmarshal config using viper
	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// unmarshal configuration into configuration option
	err = viperInstance.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// NewSpartaConfig creates a configuration from the given file name and searches
//                 the list of locations for that file.
func NewSpartaConfig(fileName string, searchPaths ...string) (*SpartaConfig, error) {
	// create a new viper instance
	viperInstance := viper.New()

	// return configuration
	return ViperSpartaConfig(viperInstance, fileName, searchPaths...)
}

func WriteConfig(config SpartaConfig, fullPath string) error {
	// create a new viper instance
	viperInstance := viper.New()

	// create map from mapstructure and struct
	encodedMap := make(map[string]interface{})
	err := mapstructure.Decode(&config, &encodedMap)
	if err != nil {
		return err
	}

	// load values into viper instance from map
	err = viperInstance.MergeConfigMap(encodedMap)
	if err != nil {
		return err
	}

	// write configuration to file
	err = viperInstance.WriteConfigAs(fullPath)
	if err != nil {
		return err
	}

	return nil
}