// Package client is a convinent wrapper of `http.Client`, with support of specifying TLS server certificates and/or mTLS certificates when creating the clients.
package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getTlsConfig(
	path ...string, // [0] - path to server cert, [1] - path to mTLS client cert, [2] - path to mTLS client key
) (
	tlsCfg *tls.Config,
	es, ec, ek, err error, // server cert error, client cert error, client key error
) {
	var pths, pthc, pthk string
	switch len(path) {
	case 2:
		err = fmt.Errorf("[CERT] both the client cert and client key are requried for mTLS.")
		return
	case 3:
		pthc = path[1]
		pthk = path[2]
		fallthrough
	case 1:
		pths = path[0]
	case 0:
		//
	default:
		err = fmt.Errorf("[CERT] invalid parameters found.")
		return
	}

	var bufs, bufc, bufk []byte
	var certPool *x509.CertPool
	var keyPair tls.Certificate
	var cp, kp int

	if pths != "" {
		bufs, err = os.ReadFile(pths)
		if err != nil {
			if !os.IsNotExist(err) {
				err = fmt.Errorf("[CERT] error reading server cert: %v", err)
				return
			} else {
				es = fmt.Errorf("[CERT] server cert '%v' missing", pths)
				err = nil
			}
		} else {
			certPool = x509.NewCertPool()
			certPool.AppendCertsFromPEM(bufs)
			cp++
		}
	}

	if pthk != "" {
		bufc, err = os.ReadFile(pthc)
		if err != nil {
			if !os.IsNotExist(err) {
				err = fmt.Errorf("[CERT] error reading client cert: %v", err)
				return
			} else {
				ec = fmt.Errorf("[CERT] client cert '%v' missing", pthc)
				err = nil
			}
		} else {
			kp++
		}
		bufk, err = os.ReadFile(pthk)
		if err != nil {
			if !os.IsNotExist(err) {
				err = fmt.Errorf("[CERT] error reading client key: %v", err)
				return
			} else {
				ek = fmt.Errorf("[CERT] client key '%v' missing", pthk)
				err = nil
			}
		} else {
			kp++
		}
		if kp == 2 {
			keyPair, err = tls.X509KeyPair(bufc, bufk)
			if err != nil {
				err = fmt.Errorf("[CERT] error preparing key pair: %v", err)
				return
			}
			kp++
		}
	}

	if kp == 3 && cp == 1 { // mTLS enabled
		tlsCfg = &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{keyPair},
		}
	} else if cp == 1 { // TLS enabled
		tlsCfg = &tls.Config{
			RootCAs: certPool,
		}
	} else {
		tlsCfg = &tls.Config{} // No cert
	}

	return
}

// Client prepare a http client.
// - timeout: client timeout
// - path:
//   - path[0]: path to server cert
//   - path[1]: path to mTLS client cert
//   - path[2]: path to mTLS client key
//
// Please note if len(path) > 1, it must be >= 3, that is, if mTLS cert is provided, the key must also be provided.
func Client(
	timeout time.Duration,
	path ...string,
) (
	client *http.Client,
	err error,
) {
	var tlsCfg *tls.Config
	tlsCfg, es, ec, ek, err := getTlsConfig(path...)
	if es != nil {
		log.Println(es)
	}
	if ec != nil {
		log.Println(ec)
	}
	if ek != nil {
		log.Println(ek)
	}
	if err != nil {
		return
	}

	client = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}
	return
}

// ClientInsecure prepare a http client which skip TLS cert verification.
func ClientInsecure(
	timeout time.Duration,
	path ...string,
) (
	client *http.Client,
	err error,
) {
	var tlsCfg *tls.Config
	tlsCfg, es, ec, ek, err := getTlsConfig(path...)
	if es != nil {
		log.Println(es)
	}
	if ec != nil {
		log.Println(ec)
	}
	if ek != nil {
		log.Println(ek)
	}
	if err != nil {
		return
	}
	tlsCfg.InsecureSkipVerify = true

	client = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}
	return
}

// GetTlsCerts examine if any, the configured TLS and mTLS certificates in the given http client.
func GetTlsCerts(
	client *http.Client,
) (
	tls *x509.CertPool,
	mtls []tls.Certificate,
	err error,
) {
	trns, okay := client.Transport.(*http.Transport)
	if !okay {
		err = fmt.Errorf("HTTP transport type mismatched")
		return
	}
	tls = trns.TLSClientConfig.RootCAs
	mtls = trns.TLSClientConfig.Certificates
	return
}
