package config

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

// Plugin is a struct that defines details about a Koffer plugin, it's
//        source, version, and other details for Koffer to collect
type Plugin struct {
	Version			string
	Service			string
	Organization 	string
	Branch			string
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
	ProviderAuth	ProviderAuth `mapstructure:"provider-pullsecret"`
	RedSord			RedSord
	Koffer			Koffer
}
