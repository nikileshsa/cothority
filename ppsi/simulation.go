package ppsi

/*
import (
	"github.com/BurntSushi/toml"
	"github.com/dedis/crypto/ppsi_crypto_utils"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"
	//"github.com/dedis/onet/simul/monitor"
)

func init() {
	onet.SimulationRegister("PPSI", NewSimulation)
}

// Simulation implements the onet.Simulation of the CoSi protocol.
type ppsiSimulation struct {
	Tree onet.SimulationBFTree
	
}

// NewSimulation returns an onet.Simulation or an error if sth. is wrong.
// Used to register the CoSi protocol.
func NewSimulation(config string) (onet.Simulation, error) {
	cs := &ppsiSimulation{}
	_, err := toml.Decode(config, cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// Setup implements onet.Simulation.
func (cs *ppsiSimulation) Setup(dir string, hosts []string) (*onet.SimulationConfig, error) {
	sim := new(onet.SimulationConfig)
	cs.CreateRoster(sim, hosts, 2000)
	err := cs.CreateTreeBFhost(sim,5,6)//need to add a method CreateTreeBFhost(sim,bf,hosts) to simul.go
	return sim, err
}



// Run implements onet.Simulation.
func (cs *ppsiSimulation) Run(config *onet.SimulationConfig) error {
	size := len(config.Roster.List)
	set1 := []string{"543323345", "543323045", "843323345", "213323045", "843323345"}
	set2 := []string{"543323345", "543323045", "843343345", "213323045", "843323345"}
	set3 := []string{"543323345", "543323045", "843323345", "213323045", "843323345"}
	set4 := []string{"543323345", "543323045", "843333345", "548323032", "213323045"}
	set5 := []string{"543323345", "543323045", "843323345", "543323245", "213323045"}
	set6 := []string{"543323345", "543323045", "843333345", "543323032", "213323045"}

	setsToEncrypt := [][]string{set1, set2, set3, set4, set5, set6}
	
	
	suite := network.Suite
	publics := config.Roster.Publics()
	ppsi := ppsi_crypto_utils.NewPPSI2(suite, publics, 6)
	EncPhones := ppsi.EncryptPhones(setsToEncrypt, 6)

	log.Lvl2("Simulation starting with: number of sets=", len(setsToEncrypt), ", Rounds=", cs.Rounds)
	for round := 0; round < cs.Rounds; round++ {
		log.Lvl1("Starting round", round)
	//	roundM := monitor.NewTimeMeasure("round")
		
		node, err := config.Overlay.CreateProtocolOnet(Name, config.Tree)
		if err != nil {
			return err
		}
		
		proto := node.(*PPSI)
		
		proto.EncryptedSets = EncPhones
		
		done := make(chan bool)
		fn := func() {
			//roundM.Record()
			done <- true
		}
		proto.RegisterSignatureHook(fn)
		if err := proto.Start(); err != nil {
			log.Error("Couldn't start protocol in round", round)
		}
		<-done
	}
	log.Lvl1("PPSI Simulation finished")
	return nil
}

*/
