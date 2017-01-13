package pop

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"errors"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"
	"github.com/dedis/cothority/crypto"
	"github.com/dedis/cothority/protocols/cosi"
	"github.com/dedis/crypto/abstract"
	crypto_cosi "github.com/dedis/crypto/cosi"

)

// ServiceName is the name to refer to the Template service from another
// package.
const ServiceName = "PoPService"

func init() {
	sda.RegisterNewService(ServiceName, newPoPService)
}

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*sda.ServiceProcessor
	path string
	data *StoredData //The data that the service will locally store
	// channel to return the configreply
	//Not sure if this is needed
	ccChannel chan *CheckConfigReply
}

type StoredData struct{
	//No creo que vaya a dejar el Pin
	//Public Key of a linked PoP, no estoy segura de a que se refiere esto
	Public abstract.Point
	//Party Transcript
	PartyTranscript *FinalTranscript
}

// StoreConfig saves the pop-config locally
func (s *Service) SendConfig(req *SendConfig) (network.Body, error) {
	log.Lvlf3("%s %v %x", s.Context.ServerIdentity(), req.ConfigFile, req.ConfigFile.Hash())
	if req.ConfigFile.Roster == nil {
		//The set of conodes has not been defined
		return nil, errors.New("The conode roster has not been defined yet")
	}
	if s.data.Public == nil {
		return nil, errors.New("Not linked yet") //Not sure what this means
	}
	//It only reserves the space for the signature, but it does not store the signature
	s.data.ConfigFile = &FinalTranscript{ConfigFile: req.ConfigFile, Signature: []byte{}}
	//Stores ConfigFile Hash, though I am not sure why
	//Returns a StoreConfigReply
	return &SendConfigReply{req.ConfigFile.Hash()}, nil
}


func (s *Service) NewProtocol(tn *sda.TreeNodeInstance, conf *sda.GenericConfig) (sda.ProtocolInstance, error) {
	log.Lvl3("Cosi Service received New Protocol event")
	pi, err := cosi.NewProtocol(tn)
	go pi.Dispatch()
	return pi, err
}

// newTemplate receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newPoPService(c *sda.Context, path string) sda.Service {
	s := &Service{
		ServiceProcessor: sda.NewServiceProcessor(c),
		path:             path,
	}
	if err := s.RegisterMessages(s.HashConfigurationFile,s.CheckHashConfigurationFile, s.SignatureRequestConfig, 
		s.HashFinalStatement, s.ChekHashFinalStatement,s.VerifyFinalStatement); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}