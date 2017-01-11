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
	"github.com/dedis/crypto/abstract"
	crypto_cosi "github.com/dedis/crypto/cosi"

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
	FinalHash  HashFinalStatement
	ConfigSig SignatureResponseConfig //Signature of the configuration file
	AttendeesPublic []abstract.Point //The set of public keys
	//Tags	[300][] bytes //Tags corresponding to each attendee when using a service

}

//HashConfigurationFile: hashes the configuration file and stores it in the service
//The hash is stored in s.hash_digest.Value of type HashConfigurationFile
func (s *Service) HashConfigurationFile(e *network.ServerIdentity, req *HashConfigurationFile) (network.Body, error) {
	log.Lvl1("Hash sum value received",req.Sum)
	s.ConfigHash.Sum = req.Sum
	return &HashConfigFileResponse{Answer: s.ConfigHash.Sum,}, nil
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
	return &CheckHashConfigFileResponse{Success: reply,}, nil
}

//HashConfigurationFile: hashes the configuration file and stores it in the service
//The hash is stored in s.hash_digest.Value of type HashConfigurationFile
func (s *Service) HashFinalStatement(e *network.ServerIdentity, req *HashFinalStatement) (network.Body, error) {
	log.Lvl1("Hash sum value received",req.Sum)
	fmt.Println("Hash Final Statement")
	s.FinalHash.Sum = req.Sum
	return &HashFinalStatementResponse{Answer: s.FinalHash.Sum,}, nil
}

//CheckHashConfigurationFile: Verifies that the Hash of the configuration file to be signed is correct
//If it is correct, returns true, else false
func (s *Service) ChekHashFinalStatement(e *network.ServerIdentity, req *CheckHashFinalStatement) (network.Body, error) {
	log.Lvl1("Hash sum value received",req.Sum)
	reply := false
	if 	bytes.Equal(req.Sum,s.FinalHash.Sum){
		reply = true
		log.Lvl1("The hash sum's are the same")
	}
	return &CheckHashFinalStatementResponse{Success: reply,}, nil
}

/*SignatureRequestConfig: starts the collective signature of a file
Useful for ConfigurationFile
*/
func (s *Service) SignatureRequestConfig(e *network.ServerIdentity, req *SignatureRequestConfig) (network.Body, error) {
	log.Lvl1("Request collective signature for configuration file",req.Message)
	tree := req.Roster.GenerateBinaryTree() //Generates the tree formed by conode servers
	tni := s.NewTreeNodeInstance(tree, tree.Root, cosi.Name)
	pi, err := cosi.NewProtocol(tni) //called from 	"github.com/dedis/cothority/protocols/cosi"
	if err != nil{
		return nil, errors.New("Error in creating CoSi protocol ")
	}
	s.RegisterProtocolInstance(pi)
	pcosi := pi.(*cosi.CoSi)
	pcosi.SigningMessage(req.Message)
	hash_sum, err := crypto.HashBytes(network.Suite.Hash(), req.Message) //Calculate message hash
	if err != nil {
		return nil, errors.New("Error hashing the message ")
	}
	response := make (chan []byte)
	pcosi.RegisterSignatureHook(func(sig []byte) {
		response <- sig
	})
	log.Lvl3("CoSi Service starting up root protocol")
	go pi.Dispatch()
	go pi.Start()
	sig := <-response
	if log.DebugVisible() > 1 {
		fmt.Printf("%s: Signed a message.\n")
	}
	log.Lvl1(sig)
	s.ConfigSig.Sum = hash_sum
	s.ConfigSig.Signature = sig
	return &SignatureResponseConfig{Sum: hash_sum, Signature: sig,}, nil
}

/*
Not sure if this goes here or not
An organizer sends a finalstatement and wants to verify its authenticity
Needs the signature
*/
func (s *Service) VerifyFinalStatement(e *network.ServerIdentity, req *VerificationStatement)(network.Body, error){
	err := crypto_cosi.VerifySignature(e.Suite(), req.ConodesPublic,req.final_msg, res.Signature)
	if err != nil{
		return &VerificationStatementResponse{Success: true,}, nil
	}else{
		return &VerificationStatementResponse{Success: false,}, err
	}

}

/*
Not sure if this goes here or not
An organizer sends a finalstatement and wants to verify its authenticity it
*/

/*To authenticate an attendee we need interaction between the parts:
The user wants to authenticate, I guess that this is done in an API, or should I do a new service
just for the user
*/
func (s *Service) AuthenticateAttendee(e *network.ServerIdentity, req *VerificationStatement)(network.Body, error){
	err := crypto_cosi.VerifySignature(e.Suite(), req.ConodesPublic,req.final_msg, res.Signature)
	if err != nil{
		return &VerificationStatementResponse{Success: true,}, nil
	}else{
		return &VerificationStatementResponse{Success: false,}, err
	}

}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolSDA, this will not be called, as the SDA will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *sda.TreeNodeInstance, conf *sda.GenericConfig) (sda.ProtocolInstance, error) {
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
	if err := s.RegisterMessages(s.HashConfigurationFile,s.CheckHashConfigurationFile, s.SignatureRequestConfig, 
		s.HashFinalStatement, s.ChekHashFinalStatement,s.VerifyFinalStatement); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}