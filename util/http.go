/*
@Time : 21-1-18
@Author : jzd
@Project: networkmgr
*/
package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mepmgr/common"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	//max memory 512MB
	maxMemory = 1024 * 1024 * 512

	HttpsAuthNone   int = 0
	HttpsAuthSingle int = 1
	HttpsAuthBoth   int = 2
)

type HttpsConf struct {
	HttpsSslType   int
	HttpsRootCrt   string
	HttpsClientCrt string
	HttpsClientKey string
}

// CopyBody returns the raw request body data as bytes.
func SafeBodyCopy(body io.Reader) []byte {
	if body == nil {
		return []byte{}
	}
	var respBody []byte
	safe := &io.LimitedReader{R: body, N: maxMemory}
	respBody, _ = ioutil.ReadAll(safe)
	return respBody
}

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return common.Unknown, err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return common.Unknown, errors.New("no valid ip found")
}

func DoRequest(method string, url string, host string, headers map[string]string, data []byte, timeoutMs int64, httpsConf *HttpsConf) (int, []byte, map[string][]string, error) {
	var reader io.Reader
	if data != nil && len(data) > 0 {
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return 0, nil, nil, err
	}

	if host != "" {
		req.Host = host
	}

	// I strongly advise setting user agent as some servers ignore request without it
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if _, f := headers["Content-Type"]; !f {
		req.Header.Set("Content-Type", "application/json")
	}

	var (
		statusCode int
		body       []byte
		timeout    time.Duration
		ctx        context.Context
		cancel     context.CancelFunc
		header     map[string][]string
	)
	timeout = time.Duration(time.Duration(timeoutMs) * time.Millisecond)
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	err = httpDo(ctx, httpsConf, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		statusCode = resp.StatusCode
		header = resp.Header

		return nil
	})

	return statusCode, body, header, err
}

// httpDo issues the HTTP request and calls f with the response. If ctx.Done is
// closed while the request or f is running, httpDo cancels the request, waits
// for f to exit, and returns ctx.Err. Otherwise, httpDo returns f's error.
func httpDo(ctx context.Context, httpsConf *HttpsConf, req *http.Request, f func(*http.Response, error) error) error {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// 是否发起https请求，若是，携带ca
	if httpsConf.HttpsSslType != HttpsAuthNone {
		pool := x509.NewCertPool()
		// 携带ca根证书
		pool.AppendCertsFromPEM([]byte(httpsConf.HttpsRootCrt))
		tlsConfig.InsecureSkipVerify = false
		tlsConfig.RootCAs = pool

		// 是否开启双向认证，若是，携带客户端证书
		if httpsConf.HttpsSslType == HttpsAuthBoth {
			cliCrt, err := tls.X509KeyPair([]byte(httpsConf.HttpsClientCrt), []byte(httpsConf.HttpsClientKey))
			if err != nil {
				fmt.Println("Loadx509keypair err:", err)
				return err
			}
			// 携带客户端证书
			tlsConfig.Certificates = []tls.Certificate{cliCrt}
		}
	}

	// Run the HTTP request in a goroutine and pass thehe response to f.
	tr := &http.Transport{
		TLSClientConfig:   tlsConfig,
		DisableKeepAlives: true,
	}

	client := &http.Client{Transport: tr}

	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
