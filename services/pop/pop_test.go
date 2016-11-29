package pop

import (
	"testing"
	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/sda"
	"github.com/dedis/cothority/network"
	"github.com/dedis/crypto/abstract"

)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func NewTestClient(lt *sda.LocalTest) *Client {
	return &Client{Client: lt.NewClient(ServiceName)}
}

func setupConfigFile() *ConfigurationFile{
	rand := network.Suite.Cipher([]byte("example"))

	X := make([]abstract.Point, 3)
	for i := range X { // pick random points
		x := network.Suite.Scalar().Pick(rand) // create a private key x
    	X[i] = network.Suite.Point().Mul(nil, x)
	}
	return &ConfigurationFile{
		OrganizersPublic : X,
		//Data: []byte{1,2,3,4,5},
	}
}

func TestServiceTemplate(t *testing.T) {
	local := sda.NewLocalTest()
	// generate 5 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, el, _ := local.GenTree(5, true)
	defer local.CloseAll()
	//dst := el.RandomServerIdentity()
	// Send a request to the service
	client := NewTestClient(local)
	log.Lvl1("Sending request to service...")
	config := setupConfigFile()
	log.Lvl1("Config file  ", config)
	hash_value, err := client.SendConfigFileHash(el,config)
	log.ErrFatal(err, "Couldn't send")
	log.Lvl1("Config File was hashed with ",hash_value)
}
