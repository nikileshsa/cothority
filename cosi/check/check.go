package check

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dedis/onet"
	"github.com/dedis/onet/app/config"
	"github.com/dedis/onet/crypto"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"

	// CoSi-protocol is not part of the cothority.
	"github.com/dedis/cothority/cosi/service"
	"github.com/dedis/crypto/cosi"
)

// RequestTimeOut is how long we're willing to wait for a signature.
var RequestTimeOut = time.Second * 10

// Config contacts all servers and verifies if it receives a valid
// signature from each.
// If the roster is empty it will return an error.
// If a server doesn't reply in time, it will return an error.
func Config(tomlFileName string) error {
	f, err := os.Open(tomlFileName)
	log.ErrFatal(err, "Couldn't open group definition file")
	group, err := config.ReadGroupDescToml(f)
	log.ErrFatal(err, "Error while reading group definition file", err)
	if len(group.Roster.List) == 0 {
		log.ErrFatalf(err, "Empty entity or invalid group defintion in: %s",
			tomlFileName)
	}
	log.Info("Checking the availability and responsiveness of the servers in the group...")
	return Servers(group)
}

// Servers contacts all servers in the entity-list and then makes checks
// on each pair. If server-descriptions are available, it will print them
// along with the IP-address of the server.
// In case a server doesn't reply in time or there is an error in the
// signature, an error is returned.
func Servers(g *config.Group) error {
	success := true
	// First check all servers individually
	for _, e := range g.Roster.List {
		desc := []string{"none", "none"}
		if d := g.GetDescription(e); d != "" {
			desc = []string{d, d}
		}
		el := onet.NewRoster([]*network.ServerIdentity{e})
		success = checkList(el, desc) == nil && success
	}
	if len(g.Roster.List) > 1 {
		// Then check pairs of servers
		for i, first := range g.Roster.List {
			for _, second := range g.Roster.List[i+1:] {
				log.Lvl3("Testing connection between", first, second)
				desc := []string{"none", "none"}
				if d1 := g.GetDescription(first); d1 != "" {
					desc = []string{d1, g.GetDescription(second)}
				}
				es := []*network.ServerIdentity{first, second}
				success = checkList(onet.NewRoster(es), desc) == nil && success
				es[0], es[1] = es[1], es[0]
				desc[0], desc[1] = desc[1], desc[0]
				success = checkList(onet.NewRoster(es), desc) == nil && success
			}
		}
	}

	if !success {
		return errors.New("At least one of the tests failed")
	}
	return nil
}

// checkList sends a message to the cothority defined by list and
// waits for the reply.
// If the reply doesn't arrive in time, it will return an
// error.
func checkList(list *onet.Roster, descs []string) error {
	serverStr := ""
	for i, s := range list.List {
		name := strings.Split(descs[i], " ")[0]
		serverStr += fmt.Sprintf("%s_%s ", s.Address, name)
	}
	log.Lvl3("Sending message to: " + serverStr)
	msg := "verification"
	fmt.Printf("Checking server(s) %s: ", serverStr)
	sig, err := signStatement(strings.NewReader(msg), list)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = verifySignatureHash([]byte(msg), sig, list)
	if err != nil {
		fmt.Printf("Invalid signature: %s\n", err.Error())
		return err
	}
	fmt.Println("Success")
	return nil
}

// signStatement signs the contents passed in the io.Reader
// (pass an io.File or use an strings.NewReader for strings). It uses
// the roster el to create the collective signature.
// In case the signature fails, an error is returned.
func signStatement(read io.Reader, el *onet.Roster) (*service.SignatureResponse,
	error) {
	//publics := entityListToPublics(el)
	client := service.NewClient()
	msg, _ := crypto.HashStream(network.Suite.Hash(), read)

	pchan := make(chan *service.SignatureResponse)
	var err error
	go func() {
		log.Lvl3("Waiting for the response on SignRequest")
		response, e := client.SignatureRequest(el, msg)
		if e != nil {
			err = e
			close(pchan)
			return
		}
		pchan <- response
	}()

	select {
	case response, ok := <-pchan:
		log.Lvl5("Response:", response)
		if !ok || err != nil {
			return nil, errors.New("received an invalid response")
		}
		err = cosi.VerifySignature(network.Suite, el.Publics(), msg, response.Signature)
		if err != nil {
			return nil, err
		}
		return response, nil
	case <-time.After(RequestTimeOut):
		return nil, errors.New("timeout on signing request")
	}
}

// verifySignatureHash verifies if the message b is correctly signed by signature
// sig from roster el.
// If the signature-check fails for any reason, an error is returned.
func verifySignatureHash(b []byte, sig *service.SignatureResponse, el *onet.Roster) error {
	// We have to hash twice, as the hash in the signature is the hash of the
	// message sent to be signed
	//publics := entityListToPublics(el)
	fHash, _ := crypto.HashBytes(network.Suite.Hash(), b)
	hashHash, _ := crypto.HashBytes(network.Suite.Hash(), fHash)
	if !bytes.Equal(hashHash, sig.Hash) {
		return errors.New("You are trying to verify a signature " +
			"belonging to another file. (The hash provided by the signature " +
			"doesn't match with the hash of the file.)")
	}
	err := cosi.VerifySignature(network.Suite, el.Publics(), fHash, sig.Signature)
	if err != nil {
		return errors.New("Invalid sig:" + err.Error())
	}
	return nil
}