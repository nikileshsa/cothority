package lib


import (
	"github.com/dedis/crypto/nist"
	"github.com/dedis/crypto/random"
	"testing"
)

func TestElgamal(t *testing.T) {

	suite := nist.NewAES128SHA256P256()

	a := suite.Scalar().Pick(random.Stream)
	A := suite.Point().Mul(nil, a)

	m := []byte("elgamal encryption")
	K, C, _ := ElGamalEncrypt(suite, A, m)

	mm, err := ElGamalDecrypt(suite, a, K, C)

	if err != nil {
		panic("decryption failed: " + err.Error())
	}
	if string(mm) != string(m) {
		panic("decryption produced wrong output: " + string(mm))
	}
	println("Decryption succeeded: " + string(mm))

}



