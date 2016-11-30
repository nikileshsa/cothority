package pop

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"time"
	"github.com/dedis/crypto/abstract"
	"github.com/dedis/cothority/network"
)

// Register messages, for network to handle them.
func init() {
	for _, msg := range []interface{}{
		FinalStatement{}, HashConfigurationFile{},SendHashConfigFileResponse{},
		ConfigurationFile{},CheckHashConfigurationFile{},SendCheckHashConfigFileResponse{},
	} {
		network.RegisterPacketType(msg)
	}
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
	Check_Value []byte
}

type HashConfigurationFile struct{
	Value []byte
}

type SendCheckHashConfigFileResponse struct {
	Success bool
}

type SendHashConfigFileResponse struct {
	Answer []byte
}
