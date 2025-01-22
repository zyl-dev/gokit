package requests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type HTTPClient struct {
	client http.Client
	Debug  bool
}

var Client *HTTPClient

// NewHTTPTLSClient 初始化 http tls 客户端
func NewHTTPTLSClient() *HTTPClient {
	pool := x509.NewCertPool()
	caCertPath := "xxx.cer"
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
	}
	pool.AppendCertsFromPEM(caCrt)

	// 本地证书
	cliCrt, err := tls.LoadX509KeyPair("cert.cer", "cert.key")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
			// 忽略服务器端证书校验，仅测试
			// 这里为了测试方便忽略了验证
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: tr}
	return &HTTPClient{client: client}
}

// NewHTTPClient 初始化 http 客户端
func NewHTTPClient() *HTTPClient {
	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 100,
		MaxIdleConns:        0,
		MaxConnsPerHost:     0,
	}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}
	return &HTTPClient{client: client}
}

// SetTimeout 重置超时时间
func (c *HTTPClient) SetTimeout(t time.Duration) {
	c.client.Timeout = t
}

// MakeGetURL : 拼接 url
func (c *HTTPClient) MakeGetURL(endpoint string, args map[string]string) string {
	url := endpoint
	first := true
	if strings.Contains(url, "?") {
		first = false
	}
	for k, v := range args {
		if first {
			url = fmt.Sprintf("%s?%s=%s", url, k, v)
			first = false
		} else {
			url = fmt.Sprintf("%s&%s=%s", url, k, v)
		}
	}
	return url
}

// MakeRequest : 构造请求
func (c *HTTPClient) MakeRequest(ctx context.Context, method, endpoint string, data io.Reader, vs []interface{}) (*http.Request, error) {
	// https://www.v2ex.com/t/622953
	// ctx, cancel := context.WithTimeout(req.Context(), 1*time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, method, endpoint, data)
	if err != nil {
		return nil, err
	}
	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case map[string]string:
			for key, value := range vv {
				req.Header[key] = []string{value}
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, err
}

// Get : get 方法, 之后可以增加重试
func (c *HTTPClient) Get(endpoint string, vs ...interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := c.MakeRequest(ctx, "GET", endpoint, nil, vs)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		byts, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(byts))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(dump))
	}
	return body, nil
}

// Post : post 方法
func (c *HTTPClient) Post(endpoint string, data io.Reader, vs ...interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := c.MakeRequest(ctx, "POST", endpoint, data, vs)
	if err != nil {
		return nil, err
	}
	// 输出到 log
	if c.Debug {
		byts, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(byts))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(dump))
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
