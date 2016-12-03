package pop

import (
	"testing"
	"time"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/sda"
	"github.com/dedis/cothority/network"
	"github.com/dedis/crypto/abstract"
	"github.com/dedis/crypto/cosi"
	"github.com/stretchr/testify/assert"



)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func NewTestClient(lt *sda.LocalTest) *Client {
	return &Client{Client: lt.NewClient(ServiceName)}
}

//Sets up an example configuration file
func setupConfigFile() *ConfigurationFile{
	rand := network.Suite.Cipher([]byte("example"))

	X := make([]abstract.Point, 3)
	for i := range X { // pick random points
		x := network.Suite.Scalar().Pick(rand) // create a private key x
    	X[i] = network.Suite.Point().Mul(nil, x)
	}
	return &ConfigurationFile{
		OrganizersPublic : X,
		StartingTime: 5.5,
		EndingTime: 5.8,
		Duration: 66.6,
		Context: []byte("IFF Forum"),
		Date: time.Date(2016, 5, 1, 12, 0, 0, 0, time.UTC),
	}
}

/*func TestSendConfigFileHash(t *testing.T) {
	local := sda.NewLocalTest()
	// generate 5 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, el, _ := local.GenTree(5, true)
	defer local.CloseAll()
	//dst := el.RandomServerIdentity() //For now a random server
	client := NewTestClient(local)
	//config := setupConfigFile()
	config := []byte("This would be a set of bytes representing the configuration file")
	//Not sure if the configuration file will be a stream of bytes or a 
	hash_value, err := client.SendConfigFileHash(el,config)
	log.ErrFatal(err, "Problem inside SendConfigFileHash")
	log.Lvl1("Config File was hashed with ",hash_value)
}*/

/*
Tests that the configuration file is signed
*/
func TestSignConfigFile(t *testing.T){
	local := sda.NewLocalTest()
	hosts, el, _:= local.GenTree(5, true)
	defer local.CloseAll()
	client := NewTestClient(local)
	//First organizers store the config-file
	config := []byte("This would be a set of bytes representing the configuration file")
	hash_value, err_configFile := client.SendConfigFileHash(el,config)
	_ = hash_value
	log.ErrFatal(err_configFile, "Problem inside SendConfigFileHash")
	log.Lvl1("Configuration file stored correctly in the server")
	//Start config file sitnature
	res_Signature, err_StartSignature := client.Start_signature_ConFigFile(el, config)
	log.ErrFatal(err_StartSignature, "Problem in the signing process")
	log.Lvl1("Configuration file signed")

	//I need to get the signature to print it
	//Send the configuration file to every server
	log.ErrFatal(err_StartSignature, "Couldn't send")
	// verify the response still
	assert.Nil(t, cosi.VerifySignature(hosts[0].Suite(), el.Publics(),
		config, res_Signature.Signature))
}

