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
		HashConfigurationFile{},HashConfigFileResponse{},
		ConfigurationFile{},CheckHashConfigurationFile{},CheckHashConfigFileResponse{},
		SignatureResponseConfig{},CheckHashFinalStatement{},HashFinalStatement{},CheckHashFinalStatementResponse{},
		HashFinalStatementResponse{},VerificationStatement{},VerificationStatementResponse{}
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

type VerificationStatement struct{
	final_msg []byte
	Roster  *sda.Roster	//The set of servers that will be used to start the collective signature protocol
	Signature []byte
	ConodesPublic []abstract.Point
}

type VerificationStatementResponse struct{
	Success bool
}
//FinalStatement File
type CheckHashFinalStatement struct{
	Sum []byte
}

type HashFinalStatement struct{
	Sum []byte
}

type CheckHashFinalStatementResponse struct {
	Success bool
}

type HashFinalStatementResponse struct {
	Answer []byte
}

//Configuration File
type CheckHashConfigurationFile struct{
	Sum []byte
}

type HashConfigurationFile struct{
	Sum []byte
}

type CheckHashConfigFileResponse struct {
	Success bool
}

type HashConfigFileResponse struct {
	Answer []byte
}
