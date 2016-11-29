package pop

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	
	"github.com/dedis/crypto/abstract"
	"github.com/dedis/cothority/network"
)

// Register messages, for network to handle them.
func init() {
	for _, msg := range []interface{}{
		FinalStatement{}, HashConfigurationFile{},SendHashConfigFileResponse{},
		ConfigurationFile{}, CountRequest{}, CountResponse{},
	} {
		network.RegisterPacketType(msg)
	}
}


type FinalStatement struct{
	AttendeesPublic []abstract.Point //The set of public keys
	//Config *ConfigurationFile //Configuration file obtained at the party
	Party_ID HashConfigurationFile //
}


type HashConfigurationFile struct{
	Value []byte
	//OrganizersPublic []abstract.Point //List of organizers public keys, NOT SURE IF THIS WILL WORK FINE FROM THE FIRST TIME
	//OrganizersPublic []byte
	/*StartingTime float64 //Starting time of party
	EndingTime float64 //Ending time of party
	Duration float64 //Duration of the party
	Use string
	ExpirationTime float64 //Nof sure if other type os better
	//Also need to store the server data, maybe a network body variable is enough
	//Agregar el marshalled value
	*/
}

type ConfigurationFile struct{
	//OrganizersPublic []abstract.Point //List of organizers public keys, NOT SURE IF THIS WILL WORK FINE FROM THE FIRST TIME
	Data []byte
	/*StartingTime float64 //Starting time of party
	EndingTime float64 //Ending time of party
	Duration float64 //Duration of the party
	Use string
	ExpirationTime float64 //Nof sure if other type os better
	//Also need to store the server data, maybe a network body variable is enough
	//Agregar el marshalled value
	*/
}
type SendHashConfigFileResponse struct {
	answer []byte
}

// CountRequest will return how many times the protocol has been run.
type CountRequest struct {
}

// CountResponse returns the number of protocol-runs
type CountResponse struct {
	Count int
}
