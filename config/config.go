package config

import (
	"fmt"
	"github.com/alexflint/go-cloudfile"
	"github.com/alexflint/go-cloudfile/s3file"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"net/url"
	"path/filepath"
	"strings"
)

// the default file name of the configuration file
const DefaultConfigName = "sparta.yml"

// viper keys for configuration to be passed in
const ViperS3Key = "s3key"
const ViperS3Secret = "s3secret"
const ViperS3Region = "s3region"
const ViperS3Url = "s3url"

// set the default region this way so the test case can add
// a custom region so that the mock/fake s3 can operate locally
var defaultRegion = aws.USGovWest

// DefaultConfig creates a default configuration where any values that
//               need to default to "non-empty" values (like booleans)
//               or strings can be set and so that the maps and arrays
//               are already constructed to prevent issues with unsafe
//               reads of the configuration
func DefaultConfig() SpartaConfig {
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
func ViperSpartaConfig(viperInstance *viper.Viper, configUrl string, searchPaths ...string) (*SpartaConfig, error) {
	parsedUrl, err := url.Parse(configUrl)
	if err != nil {
		return nil, err
	}

	// create new default config instance so that
	// parts can be defaulted before unmarshalling
	config := DefaultConfig()

	// break config name into parts and determine the
	name := filepath.Base(parsedUrl.Path)
	ext := filepath.Ext(parsedUrl.Path)
	if strings.Index(ext, ".") == 0 {
		ext = ext[1:]
	}
	if len(ext) < 1 {
		ext = "yml"
	}

	// set config file
	viperInstance.SetConfigName(name)
	viperInstance.SetConfigType(ext)

	// if given a file scheme or assume no scheme means "file"
	if parsedUrl.Scheme == "file" || len(parsedUrl.Scheme) < 1 {
		// set the location based on if the path is absolute or not
		if filepath.IsAbs(configUrl) {
			viperInstance.SetConfigFile(configUrl)
		} else {
			// add locations to search for the file with that name
			for _, location := range searchPaths {
				viperInstance.AddConfigPath(location)
			}
		}

		// unmarshal config using viper
		err = viperInstance.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// perform special handling when s3 is used to get the correct configuration for aws
		if parsedUrl.Scheme == "s3" {
			// get authentication from file
			auth, err := aws.GetAuth(viper.GetString(ViperS3Key), ViperS3Secret)

			// use public / no credential authentication
			if err != nil {
				auth = aws.Auth{
					AccessKey: "",
					SecretKey: "",
					Token:     "",
				}
			}

			// figure out region
			regionString := strings.TrimSpace(viper.GetString(ViperS3Region))
			region := defaultRegion
			if len(regionString) > 0 {
				if foundValue, found := aws.Regions[regionString]; found {
					region = foundValue
				} else {
					return nil, fmt.Errorf("The region '%s' is not a valid AWS region", regionString)
				}
			}

			// overwrite s3url for region if provided
			s3UrlString := strings.TrimSpace(viperInstance.GetString(ViperS3Url))
			if len(s3UrlString) > 0 {
				region.S3Endpoint = s3UrlString
			}

			// create driver using region and authentication
			cloudfile.Drivers["s3:"] = &s3file.Driver{
				Region: region,
				Auth:   auth,
			}
		}

		r, err := cloudfile.Open(configUrl)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		err = viperInstance.ReadConfig(r)
		if err != nil {
			return nil, err
		}
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
func NewSpartaConfig(configUrl string, searchPaths ...string) (*SpartaConfig, error) {
	// create a new viper instance
	viperInstance := viper.New()

	// return configuration
	return ViperSpartaConfig(viperInstance, configUrl, searchPaths...)
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