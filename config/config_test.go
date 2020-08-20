package config

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func loadTestData(name string, t *testing.T) SpartaConfig {
	// get current gofile/runtime location
	_, filename, _, _ := runtime.Caller(0)
	parent := filepath.Dir(filename)
	parent, _ = filepath.Abs(parent)

	// join with testdata directory
	testdata := filepath.Join(parent, "testdata")
	t.Logf("Adding testdata: %s", testdata)

	spartaConfig, err := NewSpartaConfigFromNameAndLocations(name, testdata)
	if err != nil {
		t.Errorf("Error loading configuration: %s", err)
		t.FailNow()
	}
	if spartaConfig == nil {
		t.Error("Loaded configuration was nil")
		t.FailNow()
	}

	return *spartaConfig
}

func TestSampleYaml(t *testing.T) {
	config := loadTestData("sparta.yml", t)

	// create assertion
	a := assert.New(t)

	// check OpenShift object
	ocp := config.OpenShift
	a.Equal(ocp.Version, "4.5.4")

	// check cluster object
	cluster := config.Cluster
	a.Equal(cluster.Target, "govcloud")
	a.Equal(cluster.VpcName, "iamgroot")
	a.Equal(cluster.ClusterName, "i")
	a.Equal(cluster.BaseDomain, "am.groot")
	a.Equal(cluster.ClusterDomain, "i.am.groot")
	a.Equal(cluster.AmiId, "ami-e06e5081")

	// check cloud object
	cloud := config.Cloud
	a.Equal(cloud.Provider, "aws")
	a.Equal(cloud.Region, "us-gov-west-1")
	a.Equal(cloud.VpcId, "vpc-0aef6256b40f30778")
	a.Equal(cloud.CidrPrivate, "10.0.0.0/24")

	// check subnets
	subnets := config.Subnets
	a.Equal(1, len(subnets))
	a.Contains(subnets, "private")
	private := subnets["private"]
	a.Equal(3, len(private))
	a.Contains(private,"subnet-02bf7c8c69067b993")
	a.Contains(private,"subnet-0d75d5033bfc98414")
	a.Contains(private,"subnet-058e00cfdb41ca5ce")

	// check provider-auth
	providerAuth := config.ProviderAuth
	a.True(providerAuth.Keys)
	a.Equal("XXXXXXXXXXXXXXXXXXXX", providerAuth.Secret)
	a.Equal("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", providerAuth.Key)

	// check redsord
	redSord := config.RedSord
	a.False(redSord.Enabled)

	// check koffer
	koffer := config.Koffer
	a.True(koffer.Silent)

	// check koffer plugins
	plugins := koffer.Plugins
	a.Equal(3, len(plugins))
	a.Contains(plugins, "collector-infra")
	a.Contains(plugins, "collector-operators")
	a.Contains(plugins, "collector-apps")

	// check a plugin
	collectorInfra := plugins["collector-infra"]
	a.Equal("1.0.0", collectorInfra.Version)
	a.Equal("github.com", collectorInfra.Service)
	a.Equal("codesparta", collectorInfra.Organization)
	a.Equal("master", collectorInfra.Branch)

}
