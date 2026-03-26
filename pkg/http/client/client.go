package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pangduckwai/sea9go/internal/errors"
)

func getTlsConfig(path ...string) (tlsCfg *tls.Config, err error) {
	var pths, pthc, pthk string
	switch len(path) {
	case 2:
		err = errors.Fatal("[CERT] both the client cert and client key are requried for mTLS.")
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
		err = errors.Fatal("[CERT] invalid parameters found.")
		return
	}

	var bufs, bufc, bufk []byte
	var certPool *x509.CertPool
	var keyPair tls.Certificate
	var cp, kp int
	errs := make([]error, 0)

	if pths != "" {
		bufs, err = os.ReadFile(pths)
		if err != nil {
			if !os.IsNotExist(err) {
				err = errors.Fatalf("[CERT] error reading server cert: %v", err)
				return
			} else {
				errs = append(errs, errors.NonFatalf("[CERT] server cert '%v' missing", pths))
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
				err = errors.Fatalf("[CERT] error reading client cert: %v", err)
				return
			} else {
				errs = append(errs, errors.NonFatalf("[CERT] client cert '%v' missing", pthc))
			}
		} else {
			kp++
		}
		bufk, err = os.ReadFile(pthk)
		if err != nil {
			if !os.IsNotExist(err) {
				err = errors.Fatalf("[CERT] error reading client key: %v", err)
				return
			} else {
				errs = append(errs, errors.NonFatalf("[CERT] client key '%v' missing", pthk))
			}
		} else {
			kp++
		}
		if kp == 2 {
			keyPair, err = tls.X509KeyPair(bufc, bufk)
			if err != nil {
				err = errors.Fatalf("[CERT] error preparing key pair: %v", err)
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

	if len(errs) > 0 {
		err = errors.NonFatal(errors.Errors(errs...).Error())
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
	tlsCfg, err = getTlsConfig(path...)
	if err != nil {
		if errors.IsFatal(err) {
			return
		} else {
			log.Println(err)
			err = nil
		}
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
func ClientInsecure(timeout int, path ...string) (client *http.Client, err error) {
	var tlsCfg *tls.Config
	tlsCfg, err = getTlsConfig(path...)
	if err != nil {
		if errors.IsFatal(err) {
			return
		} else {
			log.Println(err)
			err = nil
		}
	}
	tlsCfg.InsecureSkipVerify = true

	client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}
	return
}
