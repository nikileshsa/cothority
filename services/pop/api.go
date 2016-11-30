package pop

import (
	"errors"
	"crypto/sha512"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/network"
	"github.com/dedis/cothority/sda"
	_ "github.com/dedis/crypto/abstract"
	_ "github.com/dedis/crypto/anon"

)

// Client is a structure to communicate with the PoP service from outside
type Client struct {
	*sda.Client
}

// NewClient instantiates a new client with name 'n'
func NewClient() *Client {
	return &Client{Client: sda.NewClient(ServiceName)}
}

//Send the configuration file of the party, it gets hashed, the output is the hashed value
//this hashed value is also stored in the server
func (c *Client) SendConfigFileHash(dst *network.ServerIdentity, data network.Body) ([]byte, error){
	//Change so that the Number of Organizers might in data
	//dst := r.RandomServerIdentity()
	if data != nil {
		config, err := network.MarshalRegisteredType(data)
		if err != nil {
			return nil, err
		}
		hash_config := sha512.New()
		hash_config.Write(config)
		hash_config_buff := hash_config.Sum(nil)
		log.Lvl1("Hash sum value ",hash_config_buff)
		r, err := c.Send(dst, &HashConfigurationFile{
				Value: hash_config_buff,
			})
		if err != nil {
			return nil, err
		}
		replyVal := r.Msg.(SendHashConfigFileResponse)
		log.Lvl1("The hash value replied is ",replyVal.Answer)
		return replyVal.Answer, nil
	}else{
		return nil, errors.New("Empty Configuration File")
	}
}


//Start_Signature returns
/*
Starts a collective signature round. It is expected that each conode signs a public key of a party attendat
Recibe? Statement to be signed, that is the public key
Regresa? The aggregate commit of the signed key
*/
/*Start_signature_ConFigFile signs the configuration file
/Steps:
1. Receives configuraiton file
2. Checks that it is valid (with the hash stored previously)
3. If the file is valid, then the file is signed
4. Returns error if the file has not valid
*/

func (c *Client) Start_signature_ConFigFile(dst *network.ServerIdentity, data network.Body) (error){
	if data != nil {
			config, err := network.MarshalRegisteredType(data)
			if err != nil {
				return nil, err
			}
			hash_config := sha512.New()
			hash_config.Write(config)
			hash_config_buff := hash_config.Sum(nil)
			//The value to check
			r, err := c.Send(dst, &CheckHashConfigurationFile{
					Check_Value: hash_config_buff,
				})
			if err != nil {
				return nil, err
			}
			replyVal := r.Msg.(SendCheckHashConfigFileResponse)
			log.Lvl1("Success ",replyVal.Success)
			if replyVal.Sucess == false{
				return nil, errors.New("Configuration file is incorrect")
			}
			//If configuration file is correct, start the signing process
			//How do I start a cosi round?
		}else{
			return nil, errors.New("Empty Configuration File")
		}
}

/*
Organizers send statements containing public keys and party configuration information.
This information is stored so that it will be compared when the signature round starts.
NOTA, no estoy muy segura como esta la signature round, osea como empieza pues y como se comunican
Receive a toml file, with an array that contains the public keys to be signed
public_keys = ["ZxYyfezvhCIw5c7C7KIYIJ4xCgo9VNh/YbylBIotOHk=", "ZxYyfezvhCIw5c7C7KIYIJ4xCgo9VNh/YbylBIotOHk=", "ZxYyfezvhCIw5c7C7KIYIJ4xCgo9VNh/YbylBIotOHk="]
*/
//func (c *Client) Send_statements(in io.Reader) (error){
/*
Tiene que guardar eso que lee en algun lado
Va a firmar que? Lista de clase publicas, archvos de video y configuraciones de la fiesta
Que regresamos? pues estas cosas firmadas no?
*/
//}

//func (c *Client) Set_up() (error){
/*
Supongo que aqui solicitan que se inicie una ronda de firmas con los nodos
*/
//}

//func (c *Client) Interaction_UnlimitID() (error){
//Empty for now, but it will be the authentication service that will connect to the IdP and so on when using UnlimitID, i think. But not sure
//}

