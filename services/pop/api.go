package pop

import (
	_ "errors"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"

	"github.com/satori/go.uuid"
	"github.com/BurntSushi/toml"

   _ "github.com/dedis/cothority/crypto"
	"github.com/dedis/crypto/base64"

	"github.com/dedis/crypto/abstract"
	_ "github.com/dedis/crypto/anon"
)

// Client is a structure to communicate with Guard service
type Client struct {
	*sda.Client
	//Saves the data, which I am not sure what it is yet
}


type FinalTranscript struct{
	ConfigFile *ConfigurationFile
	Attendees []abstract.Point
	Signature []byte
}

type FinalTranscriptToml struct {
	ConfigFile *ConfigurationFileToml
	Attendees	[]string
	Signature 	string
}

type ConfigurationFile struct {
	Name     string
	DateTime string
	Location string
	Cothority   *sda.Roster
}

type ConfigurationFileToml struct {
	Name     string
	DateTime string
	Location string
	Cothority   [][]string
}

//Register the Packet files, that will be sent over the network
func init() {
	network.RegisterPacketType(&FinalTranscript{})
	network.RegisterPacketType(&ConfigurationFile{})
}

// NewClient makes a new Client
func NewClient(cothority *sda.Roster) *Client {
	return &Client{Client: sda.NewClient(ServiceName)}

}

// NewFinalStatementFromString creates a final statement from a string
func NewFinalStatementFromString(s string) *FinalTranscript {
	fsToml := &FinalTranscriptToml{}
	_, err := toml.Decode(s, fsToml)
	if err != nil {
		log.Error(err)
		return nil
	}
	sis := []*network.ServerIdentity{}
	for _, s := range fsToml.ConfigFile.Cothority {
		uid, err := uuid.FromString(s[2])
		if err != nil {
			log.Error(err)
			return nil
		}
		sis = append(sis, &network.ServerIdentity{
			Address:     network.Address(s[0]),
		//	Description: s[1],
			ID:          network.ServerIdentityID(uid),
			Public:      B64ToPoint(s[3]),
		})
	}
	coth := sda.NewRoster(sis)
	config_file := &ConfigurationFile{
		Name:     fsToml.ConfigFile.Name,
		DateTime: fsToml.ConfigFile.DateTime,
		Location: fsToml.ConfigFile.Location,
		Cothority:   coth,
	}
	atts := []abstract.Point{}
	for _, p := range fsToml.Attendees {
		atts = append(atts, B64ToPoint(p))
	}
	sig := make([]byte, 64)
	sig, err = base64.StdEncoding.DecodeString(fsToml.Signature)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &FinalTranscript{
		ConfigFile:      config_file,
		Attendees: atts,
		Signature: sig,
	}
}

// SendConfig sends the configuration to the conode for later usage.
func (c *Client) SendConfig(dst network.Address, configFile *ConfigurationFile) error {
	si := &network.ServerIdentity{Address: dst}
	//Cambiar esto, porque no estoy utilizando protobuf
	//Este es con el send normal, que no recuerdo como utilizar
	err := c.SendProtobuf(si, &SendConfig{p}, nil)
	if err != nil {
		return err
	}
	return nil
}

/*
*Functions to manage parsing of points and scalars to base64
*/
// PointToB64 converts an abstract.Point to a base64-point.
func PointToB64(p abstract.Point) string {
	pub, err := p.MarshalBinary()
	if err != nil {
		log.Error(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(pub)
}

// B64ToPoint converts a base64-string to an abstract.Point.
func B64ToPoint(str string) abstract.Point {
	public := network.Suite.Point()
	buf, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Error(err)
		return nil
	}
	err = public.UnmarshalBinary(buf)
	if err != nil {
		log.Error(err)
		return nil
	}
	return public
}

// ScalarToB64 converts an abstract.Scalar to a base64-string.
func ScalarToB64(s abstract.Scalar) string {
	sec, err := s.MarshalBinary()
	if err != nil {
		log.Error(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(sec)
}

// B64ToScalar converts a base64-string to an abstract.Scalar.
func B64ToScalar(str string) abstract.Scalar {
	scalar := network.Suite.Scalar()
	buf, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Error(err)
		return nil
	}
	err = scalar.UnmarshalBinary(buf)
	if err != nil {
		log.Error(err)
		return nil
	}
	return scalar
}
