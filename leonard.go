package leonard

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"sir"
	"time"
)

type Config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)

		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(rwTimeout))

		return conn, nil
	}
}

func FetchUrl(theurl string) string {
	var client *http.Client

	if proxy := os.Getenv("http_proxy"); proxy != `` {
		proxyUrl, err := url.Parse(proxy)
		sir.CheckError(err)

		transport := http.Transport{
			Dial:  TimeoutDialer(5*time.Second, 5*time.Second), // connect, read/write
			Proxy: http.ProxyURL(proxyUrl),
		}

		client = &http.Client{Transport: &transport}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(`GET`, theurl, nil)
	sir.CheckError(err)

	resp, err := client.Do(req)
	sir.CheckError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	sir.CheckError(err)

	return string(body)
}
