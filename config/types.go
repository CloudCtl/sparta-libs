package config

import "fmt"

// OpenShift defines the openshift configuration details related to artifacts and sources
type OpenShift struct {
	Version		string
}

// Cluster defines the cluster logical details like names and networking
type Cluster struct {
	Target			string
	VpcName			string `mapstructure:"vpc-name"`
	ClusterName		string `mapstructure:"cluster-name"`
	BaseDomain		string `mapstructure:"base-domain"`
	ClusterDomain	string `mapstructure:"cluster-domain"`
	AmiId			string `mapstructure:"ami-id"`
}

// Cloud is information on the target cloud environment and the VPCs that
//       are created there
type Cloud struct {
	Provider		string
	Region			string
	VpcId			string `mapstructure:"vpc-id"`
	CidrPrivate		string `mapstructure:"cidr-private"`
}

// SubnetGroup is a list of VPC subnet IDs
type SubnetGroup []string

// Subnets is a map of names (private/public) to a list of subnet ids (SubnetGroup) found
//         in that subnet name
type Subnets map[string]SubnetGroup

// ProviderAuth is the authentication structure for the provider (such as AWS)
type ProviderAuth struct {
	Keys			bool
	Secret			string
	Key				string
}

// RedSord
type RedSord struct {
	Enabled        	bool
}

// An EnvironmentVariable is a struct that holds a name/value pair because a straight-up map[string]string
//                        can't be used since viper keys are case insensitive and will be all lower case in the
//                        map resulting to unexpected behavior in the plugin
//                        see: https://github.com/spf13/viper/issues/411
//                        and: https://github.com/spf13/viper/issues/373
type EnvironmentVariable struct {
	Name	string
	Value	string
}

// Pair converts the name/value into the name=value format that is expected
//      when the command is being run
func (env EnvironmentVariable) Pair() string {
	return fmt.Sprintf("%s=%s", env.Name, env.Value)
}

// EnvironmentVariables is a slice of name->values that is used to send environment
//                      to the execution of a plugin. This is used in place of a map
//                      because viper keys are case-insensitive and the resulting map
//                      would all be lowercase
type EnvironmentVariables []EnvironmentVariable

// List creates a slice of "name=value" pairs that is suitable for use with the
//      plugin command.
func (vars EnvironmentVariables) List() []string {
	output := make([]string, 0, len(vars))
	for idx := 0; idx < len(vars); idx++ {
		output = append(output, vars[idx].Pair())
	}
	return output
}

// Map creates a map of values matching the map[name]value structure
func (vars EnvironmentVariables) Map() map[string]string {
	output := make(map[string]string)
	for idx := 0; idx < len(vars); idx++ {
		output[vars[idx].Name] = vars[idx].Value
	}
	return output
}

// Plugin is a struct that defines details about a Koffer plugin, it's
//        source, version, and other details for Koffer to collect
type Plugin struct {
	Version			string
	Service			string
	Organization 	string
	Branch			string
	Env				EnvironmentVariables
}

// Plugins is a map that names plugins to their source
type Plugins map[string]Plugin

// Koffer is a collector for CodeSparta artifacts
type Koffer struct {
	Silent			bool
	Plugins			Plugins
}

// SpartaConfig is the holder for all of the parts of the configuration
type SpartaConfig struct {
	OpenShift		OpenShift
	Cluster			Cluster
	Cloud			Cloud
	Subnets			Subnets
	ProviderAuth	ProviderAuth `mapstructure:"provider-auth"`
	RedSord			RedSord
	Koffer			Koffer
}
