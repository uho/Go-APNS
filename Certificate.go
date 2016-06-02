package goapns

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

var (
	ErrorCertificateExpired          = errors.New("Your certificate has expired. Please renew in Apples Developer Center")
	ErrorCertificatePrivateKeyNotRSA = errors.New("Apparently the private key is not in RSA format, aborting.")
)

func CertificateFromP12(filePath string, key string) (tls.Certificate, error) {
	p12Data, err := ioutil.ReadFile(filePath)
	fmt.Printf("Read Data %v error: %v\n", p12Data, err)
	if err != nil {
		return tls.Certificate{}, err
	}

	privateKey, crt, err := pkcs12.Decode(p12Data, key)
	if err != nil {
		fmt.Printf("Could not load cert with error %v", err)
		return tls.Certificate{}, err
	}
	fmt.Printf("Decoded. Private key %v crt %v, error %v \n", privateKey, crt, err)

	//ensure that private key is RSA
	privateRSAKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return tls.Certificate{}, ErrorCertificatePrivateKeyNotRSA
	}

	certificate := tls.Certificate{
		Certificate: [][]byte{crt.Raw},
		PrivateKey:  privateRSAKey,
		Leaf:        crt,
	}

	return certificate, nil
}

func verify(cert *x509.Certificate) error {
	_, err := cert.Verify(x509.VerifyOptions{})
	if err == nil {
		return nil
	}
	switch e := err.(type) {

	case x509.CertificateInvalidError:
		if e.Reason == x509.Expired {
			return ErrorCertificateExpired
		}
		return err

	default:
		return err
	}
}
