/*
* This is a template for creating an app. It only has one command which
* prints out the name of the app.
 */
package main

import (
	"os"

	"bytes"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/network"
    "github.com/dedis/crypto/random"
    "github.com/dedis/cothority/crypto"

	"github.com/boombuler/barcode"
    "github.com/boombuler/barcode/qr"
	
	"github.com/jung-kurt/gofpdf"
	pdf_barcode "github.com/jung-kurt/gofpdf/contrib/barcode"


    "image/png"

   // log_2 "log"
	/*"io/ioutil"

	"path"

	"errors"

	"encoding/base64"

	"net"

	"fmt"

	"strings"

	"github.com/BurntSushi/toml"


	"github.com/dedis/crypto/abstract"
	"github.com/dedis/crypto/anon"
	"github.com/dedis/crypto/random"

	"github.com/dedis/cothority/app/lib/config"

	"github.com/dedis/cothority/sda"
	"github.com/dedis/cothority/services/pop"
	"github.com/dedis/cothority/network"

	s "github.com/SSSaaS/sssa-golang"*/

	"gopkg.in/urfave/cli.v1"
	//Falta agregar la libreria de cothority, pero espera, igual y no es necesario NO SE
)

func main() {
	app := cli.NewApp()
	app.Name = "Template"
	app.Usage = "Used for building other apps."
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:      "main",
			Usage:     "main command",
			Aliases:   []string{"m"},
			ArgsUsage: "additional parameters",
			Action:    cmdMain,
		},
		{
			Name:      "generate_keys",
			Usage:     "Generate public key, private key and ephemeral public key",
			Aliases:   []string{"m"},
			ArgsUsage: "salt",
			Action:    gen_keys,
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "debug, d",
			Value: 0,
			Usage: "debug-level: 1 for terse, 5 for maximal",
		},
	}
	app.Before = func(c *cli.Context) error {
		log.SetDebugVisible(c.Int("debug"))
		return nil
	}
	app.Run(os.Args)

}

// Main command.
func cmdMain(c *cli.Context) error {
	log.Info("Main command")
	return nil
}

func create_QRCode(file_name string, key_encode string) error{
	f, _ := os.Create(file_name)
	defer f.Close()

	qrcode, err := qr.Encode(key_encode, qr.L, qr.Auto)
	if err != nil {
	    return err
	} else {
	    qrcode, err = barcode.Scale(qrcode, 300, 300)
	    if err != nil {
	        return err
	    } else {
	        png.Encode(f, qrcode)
	        return nil
	    }
	}

}

func create_PDF(file_name string, qr_puk string, qr_ek string, qr_prk string) error{

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "PSEUDONYM PARTY SET OF KEYS")
	pdf.Ln(10)
	pdf.Cell(0, 10, "Ephemeral Public Key")


	/*pdf.Image(qr_ek, 10, 10, 30, 10, false, "", 0, "")

	//pdf.MultiCell(40, 10, "Pseudonnym Party Set of Keys", "", "", false)
	err := pdf.OutputFileAndClose(file_name)
	if (err!=nil){
		log.Info(err)
	}*/

	//pdf = gofpdf.New("L", "mm", "A4", "")
	//pdf.AddPage()

	key := pdf_barcode.RegisterQR(pdf, qr_ek, qr.L, qr.Auto)
	pdf_barcode.Barcode(pdf, key, 15, 30, 40, 40, false)
	pdf.Ln(55)
	pdf.Cell(0,10, qr_ek)
	pdf.Ln(15)

	pdf.Cell(0, 0, "Public Key")
	key = pdf_barcode.RegisterQR(pdf, qr_puk, qr.L, qr.Auto)
	pdf_barcode.Barcode(pdf, key, 15, 100, 40, 40, false)
	pdf.Ln(55)
	pdf.Cell(0,10, qr_puk)
	pdf.Ln(25)

	pdf.Line(0, 160, 260, 160)


	pdf.Cell(0, 0, "Private Key")
	key = pdf_barcode.RegisterQR(pdf, qr_prk, qr.L, qr.Auto)
	pdf_barcode.Barcode(pdf, key, 15, 180, 40, 40, false)
	pdf.Ln(55)
	pdf.Cell(0,10, qr_prk)
	pdf.Ln(15)

	err := pdf.OutputFileAndClose("pseudonym_party_keys.pdf")
	if (err!=nil){
		log.Info(err)
		return err
	}
	return nil
}
func gen_keys(c *cli.Context) error{
	log.Info("Create public key, private key and ephemeral key")
	private_k := network.Suite.Scalar().Pick(random.Stream) //Generate private key
	public_K := network.Suite.Point().Mul(nil, private_k) //Generate public key
	//El party_ID es el config file que recibe como parametro
	//Revisar como mandar a llamar parametros
	party_ID := []byte("Hello World!")
	party_hash_sum, err := crypto.HashBytes(network.Suite.Hash(),party_ID)
	party_cipher := network.Suite.Cipher(party_hash_sum)
	party_Base, _ := network.Suite.Point().Pick(nil,party_cipher)
	ephemeral_K := network.Suite.Point().Mul(party_Base,private_k)
	pub_string, err := crypto.Pub64(network.Suite,public_K)
	eph_string, err := crypto.Pub64(network.Suite,ephemeral_K)
	//priv_string, err := crypto.Pub64(network.Suite,private_k)
	_ = err
	var buff bytes.Buffer
	buff.Reset()
	crypto.WriteScalar64(network.Suite,&buff,private_k) //Write the scalar version on 64 bits
	priv_string := buff.String()
	log.Info(pub_string)
	log.Info(eph_string)
	log.Info(priv_string)	

	//create_QRCode("qrcode_pub.png",pub_string)
	//create_QRCode("qrcode_eph.png",eph_string)
	//create_QRCode("qrcode_priv.png",priv_string)
	//Create PDF
	create_PDF("party_keys.pdf",pub_string,eph_string,priv_string)


	return nil
}