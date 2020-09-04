package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// load data in the same way for test cases
func loadTestData(name string, t *testing.T) SpartaConfig {
	// get current gofile/runtime location
	_, filename, _, _ := runtime.Caller(0)
	parent := filepath.Dir(filename)
	parent, _ = filepath.Abs(parent)

	// join with testdata directory
	testdata := filepath.Join(parent, "testdata")

	spartaConfig, err := NewSpartaConfig(name, testdata)
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

// single function for checking test values
func assertSampleData(config SpartaConfig, t *testing.T) {
	// create assertion
	a := assert.New(t)

	// check OpenShift object
	ocp := config.OpenShift
	a.Equal(ocp.Version, "4.5.4")

	// check cluster object
	cluster := config.Cluster
	a.NotNil(cluster)
	a.Equal(cluster.Target, "govcloud")
	a.Equal(cluster.VpcName, "iamgroot")
	a.Equal(cluster.ClusterName, "i")
	a.Equal(cluster.BaseDomain, "am.groot")
	a.Equal(cluster.ClusterDomain, "i.am.groot")
	a.Equal(cluster.AmiId, "ami-e06e5081")

	// check cloud object
	cloud := config.Cloud
	a.NotNil(cloud)
	a.Equal(cloud.Provider, "aws")
	a.Equal(cloud.Region, "us-gov-west-1")
	a.Equal(cloud.VpcId, "vpc-0aef6256b40f30778")
	a.Equal(cloud.CidrPrivate, "10.0.0.0/24")

	// check subnets
	subnets := config.Subnets
	a.NotNil(subnets)
	a.Equal(1, len(subnets))
	a.Contains(subnets, "private")
	private := subnets["private"]
	a.Equal(3, len(private))
	a.Contains(private,"subnet-02bf7c8c69067b993")
	a.Contains(private,"subnet-0d75d5033bfc98414")
	a.Contains(private,"subnet-058e00cfdb41ca5ce")

	// check provider-auth
	providerAuth := config.ProviderAuth
	a.NotNil(providerAuth)
	a.True(providerAuth.Keys)
	a.Equal("XXXXXXXXXXXXXXXXXXXX", providerAuth.Secret)
	a.Equal("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", providerAuth.Key)

	// check redsord
	redSord := config.RedSord
	a.NotNil(redSord)
	a.False(redSord.Enabled)

	// check koffer
	koffer := config.Koffer
	a.NotNil(koffer)
	a.True(koffer.Silent)

	// check koffer plugins
	plugins := koffer.Plugins
	a.NotNil(plugins)
	a.Equal(3, len(plugins))
	a.Contains(plugins, "collector-infra")
	a.Contains(plugins, "collector-operators")
	a.Contains(plugins, "collector-apps")

	// check a plugin
	collectorInfra := plugins["collector-infra"]
	a.NotNil(collectorInfra)
	a.Equal("4.5.6", collectorInfra.Version)
	a.Equal("github.com", collectorInfra.Service)
	a.Equal("codesparta", collectorInfra.Organization)
	a.Equal("master", collectorInfra.Branch)
}

// read a sample yaml file from the repo and assert it is correct
func TestReadSampleYaml(t *testing.T) {
	config := loadTestData("sparta.yml", t)
	assertSampleData(config, t)
}

// read a sample yaml file from the given endpoint and assert it is correct
func TestReadSampleHttpData(t *testing.T) {
	config, err := NewSpartaConfig("https://codesparta-testdata.s3-us-gov-west-1.amazonaws.com/sparta.yml")

	// create assertion that err is nil and config is not nil
	a := assert.New(t)
	a.Nil(err)
	a.NotNil(config)

	// check values
	assertSampleData(*config, t)
}

// read a sample file from an s3 bucket and assert it is correct
func TestReadSampleS3Data(t *testing.T) {
	config, err := NewSpartaConfig("s3://codesparta-testdata/sparta.yml")

	// create assertion that err is nil and config is not nil
	a := assert.New(t)
	a.Nil(err)
	a.NotNil(config)

	// check values
	assertSampleData(*config, t)
}

// test that a configuration can be loaded, written to disk, and that the round-trip comes out with the same values
func TestWriteSampleYaml(t *testing.T) {
	// create assertion
	a := assert.New(t)

	// load config
	config := loadTestData("sparta.yml", t)

	// create and defer removal of tmp dir
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	fileName := "test-output.yml"
	tmpFile := filepath.Join(tmpDir, fileName)

	err := WriteConfig(config, tmpFile)
	a.Nil(err)

	// read config
	writtenConfig, err := NewSpartaConfig(fileName, tmpDir)
	a.Nil(err)
	a.NotNil(writtenConfig)

	// assert written test data is the same
	assertSampleData(*writtenConfig, t)
}
