package monitor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	//"path"

	"github.com/cretz/bine/tor"
)

func FetchReseed(url, cert string) ([]byte, error) {
	url = url + "/i2pseeds.su3"
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	onion := false
	var t *tor.Tor
	if strings.Contains(url, ".onion") {
		onion = true
		var err error
		t, err = tor.Start(nil, nil)
		if err != nil {
			return nil, err
		}
		defer t.Close()
	}
	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	// Make connection
	var dialer *tor.Dialer
	if onion {
		var err error
		dialer, err = t.Dialer(dialCtx, nil)
		if err != nil {
			return nil, err
		}
	}
	if cert != "le" {
		caCert, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		var client *http.Client

		if onion {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: caCertPool,
					},
					DialContext: dialer.DialContext,
				},
			}
		} else {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: caCertPool,
					},
					//					DialContext: dialer.DialContext,
				},
			}
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "Wget/1.11.4")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}
	var client *http.Client
	if onion {
		client = &http.Client{
			Transport: &http.Transport{DialContext: dialer.DialContext},
		}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Wget/1.11.4")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
