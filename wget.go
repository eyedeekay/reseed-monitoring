package monitor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"time"
	//"path"

	"github.com/cretz/bine/tor"
)

func FetchReseed(url, cert string) ([]byte, error) {
	url = url + "/i2pseeds.su3"
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	t, err := tor.Start(nil, nil)
	if err != nil {
		return nil, err
	}
	defer t.Close()
	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, err
	}
	if cert != "le" {
		caCert, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		//	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
		//	if err != nil {
		//		log.Fatal(err)
		//	}

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
					//				Certificates: []tls.Certificate{cert},
				},
				DialContext: dialer.DialContext,
			},
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
	client := &http.Client{
		/*Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				//				Certificates: []tls.Certificate{cert},
			},
		},*/
		Transport: &http.Transport{DialContext: dialer.DialContext},
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
