package pop

import (

	"github.com/dedis/crypto/abstract"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"
)

/*
*Register messages, for the network to handle them
*/

// Register messages, for network to handle them.
func init() {
	for _, msg := range []interface{}{
		SendConfig{}, SendConfigReply{},
		EndPartyResponse{},EndPartyResponse{},
		CheckConfig{}, CheckConfigReply{},
	} {
		network.RegisterPacketType(msg)
	}
}



// CheckConfig asks whether the pop-config and the attendees are available.
type CheckConfig struct {
	ConfigHash_ID   []byte
	Attendees []abstract.Point
}

// CheckConfigReply sends back an integer for the Pop:
// - 0 - no popconfig yet
// - 1 - popconfig, but other hash
// - 2 - popconfig with the same hash but no attendees in common
// - 3 - popconfig with same hash and at least one attendee in common
// if PopStatus == 3, then the Attendees will be the common attendees between
// the two nodes.
type CheckConfigReply struct {
	Status 		int
	ConfigHash_ID   []byte
	Attendees []abstract.Point
}

// SendConfig presents a Configuration File to be stored
type SendConfig struct {
	ConfigFile *ConfigurationFile
}

//COMO LO TIENEN AHORITA
//Identifies a configuration file with its hash
//TODO: StoreConfigReply will give in a later version a handler that
//can be used to indetifiy that config, no estoy segura de a que se refiere
type SendConfigReply struct{
	Hash_ID []byte
}

// Ask to end a particular pop party, based on the Hash_ID of the config file.
// TODO: support more than one popconfig
type EndPartyRequest struct {
	DescID    []byte
	Attendees []abstract.Point
}

// FinalizeResponse returns the FinalStatement if all conodes already received
// a PopDesc and signed off. The FinalStatement holds the updated PopDesc, the
// pruned attendees-public-key-list and the collective signature.
type EndPartyResponse struct {
	Final *FinalTranscript
}