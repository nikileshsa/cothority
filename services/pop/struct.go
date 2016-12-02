package pop

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"time"
	"github.com/dedis/crypto/abstract"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"
)

// Register messages, for network to handle them.
func init() {
	for _, msg := range []interface{}{
		FinalStatement{}, 
		HashConfigurationFile{},SendHashConfigFileResponse{},
		ConfigurationFile{},CheckHashConfigurationFile{},SendCheckHashConfigFileResponse{},
		SignatureResponseConfig{},
	} {
		network.RegisterPacketType(msg)
	}
}

// SignatureRequest is what the Cosi service is expected to receive from clients.
type SignatureRequestConfig struct {
	Message []byte
	Roster  *sda.Roster	//The set of servers that will be used to start the collective signature protocol
}

// SignatureResponse is what the Cosi service will reply to clients.
type SignatureResponseConfig struct {
	Sum       []byte
	Signature []byte
}

type ConfigurationFile struct{
	OrganizersPublic []abstract.Point //List of organizers public keys
	StartingTime float64 //Starting time of party not sure if better use duration
	EndingTime float64 //End time of the party
	Duration float64 //Measured in hours and minutes
	Context []byte //Scope, what the token will be used for
	Date time.Time
}

type FinalStatement struct{
	AttendeesPublic []abstract.Point //The set of public keys
	//Config *ConfigurationFile //Configuration file obtained at the party
	Party_ID HashConfigurationFile
	RealStartingTime float64
	RealEndingTime float64
	//Not sure how to put the observers video Files
}

type CheckHashConfigurationFile struct{
	Sum []byte
}

type HashConfigurationFile struct{
	Sum []byte
}

type SendCheckHashConfigFileResponse struct {
	Success bool
}

type SendHashConfigFileResponse struct {
	Answer []byte
}
