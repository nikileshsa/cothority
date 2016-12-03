package pop

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	_ "time"
	"errors"
	"fmt"
	"bytes"
	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"
	"github.com/dedis/cothority/crypto"
	"github.com/dedis/cothority/protocols/cosi"
)

// ServiceName is the name to refer to the Template service from another
// package.
const ServiceName = "PoP"

func init() {
	sda.RegisterNewService(ServiceName, newService)
}

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*sda.ServiceProcessor
	path string
	// Count holds the number of calls to 'ClockRequest'
	PoPConfig ConfigurationFile
	Count int
	ConfigHash HashConfigurationFile
}

//HashConfigurationFile: hashes the configuration file and stores it in the service
//The hash is stored in s.hash_digest.Value of type HashConfigurationFile
func (s *Service) HashConfigurationFile(e *network.ServerIdentity, req *HashConfigurationFile) (network.Body, error) {
	log.Lvl1("Hash sum value received",req.Sum)
	s.ConfigHash.Sum = req.Sum
	return &SendHashConfigFileResponse{Answer: s.ConfigHash.Sum,}, nil
}

//CheckHashConfigurationFile: Verifies that the Hash of the configuration file to be signed is correct
//If it is correct, returns true, else false
func (s *Service) CheckHashConfigurationFile(e *network.ServerIdentity, req *CheckHashConfigurationFile) (network.Body, error) {
	log.Lvl1("Hash sum value received",req.Sum)
	reply := false
	if 	bytes.Equal(req.Sum,s.ConfigHash.Sum){
		reply = true
		log.Lvl1("The hash sum's are the same")
	}
	fmt.Println("CheckHashConfigurationFile salio bien")
	return &SendCheckHashConfigFileResponse{Success: reply,}, nil
}


/*SignatureRequestConfig: starts the collective signature of a file
Useful for ConfigurationFile
*/
func (s *Service) SignatureRequestConfig(e *network.ServerIdentity, req *SignatureRequestConfig) (network.Body, error) {
	log.Lvl1("Request collective signature for configuration file",req.Message)
	tree := req.Roster.GenerateBinaryTree() //Generates the tree formed by conode servers
	tni := s.NewTreeNodeInstance(tree, tree.Root, cosi.Name)
	pi, err := cosi.NewProtocol(tni) //called from 	"github.com/dedis/cothority/protocols/cosi"
	fmt.Println("Ya instancio el protocolo")
	if err != nil{
		return nil, errors.New("Error in creating CoSi protocol ")
	}
	s.RegisterProtocolInstance(pi)
	fmt.Println("Registra el protocolo")
	pcosi := pi.(*cosi.CoSi)
	pcosi.SigningMessage(req.Message)
	fmt.Println("Signing Message")
	hash_sum, err := crypto.HashBytes(network.Suite.Hash(), req.Message) //Calculate message hash
	fmt.Println("Crypto hash")
	if err != nil {
		return nil, errors.New("Error hashing the message ")
	}
	response := make (chan []byte)
	fmt.Println("Creating channel")
	pcosi.RegisterSignatureHook(func(sig []byte) {
		response <- sig
	})
	fmt.Println("Register Signature Hook")
	log.Lvl3("CoSi Service starting up root protocol")
	go pi.Dispatch()
	go pi.Start()
	sig := <-response
	if log.DebugVisible() > 1 {
		fmt.Printf("%s: Signed a message.\n")
	}
	fmt.Println("Sig is asigned repsonse")
	fmt.Println("The signature is:")
	fmt.Println(sig)
	return &SignatureResponseConfig{Sum: hash_sum, Signature: sig,}, nil
}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolSDA, this will not be called, as the SDA will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *sda.TreeNodeInstance, conf *sda.GenericConfig) (sda.ProtocolInstance, error) {
	fmt.Println("Entre a New Protocol")
	log.Lvl3("Cosi Service received New Protocol event")
	pi, err := cosi.NewProtocol(tn)
	go pi.Dispatch()
	return pi, err
}

// newTemplate receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *sda.Context, path string) sda.Service {
	s := &Service{
		ServiceProcessor: sda.NewServiceProcessor(c),
		path:             path,
	}
	if err := s.RegisterMessages(s.HashConfigurationFile,s.CheckHashConfigurationFile, s.SignatureRequestConfig); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}