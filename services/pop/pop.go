package pop

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	_ "time"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"
	_ "github.com/dedis/cothority/protocols"
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
	hash_digest HashConfigurationFile
}

//HashConfigurationFile: hashes the configuration file and stores it in the service
//The hash is stored in s.hash_digest.Value of type HashConfigurationFile
func (s *Service) HashConfigurationFile(e *network.ServerIdentity, req *HashConfigurationFile) (network.Body, error) {
	//log.Lvl1("HashConfigurationFile")
	s.hash_digest.Value = req.Value
	reply := &SendHashConfigFileResponse{s.hash_digest.Value}
	return reply, nil
}

//Previous testing of CountRequest
func (s *Service) CountRequest(e *network.ServerIdentity, req *CountRequest) (network.Body, error) {
	return &CountResponse{s.Count}, nil
}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolSDA, this will not be called, as the SDA will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *sda.TreeNodeInstance, conf *sda.GenericConfig) (sda.ProtocolInstance, error) {
	log.Lvl3("Not templated yet")
	return nil, nil
}

// newTemplate receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *sda.Context, path string) sda.Service {
	s := &Service{
		ServiceProcessor: sda.NewServiceProcessor(c),
		path:             path,
	}
	if err := s.RegisterMessages(s.HashConfigurationFile, s.CountRequest); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}