package sockets

import (
	"net"
	"net/url"
	"os"
	"strings"
	"fmt"

	"golang.org/x/net/proxy"
)

// GetProxyEnv allows access to the uppercase and the lowercase forms of
// proxy-related variables.  See the Go specification for details on these
// variables. https://golang.org/pkg/net/http/
func GetProxyEnv(key string) string {
	proxyValue := os.Getenv(strings.ToUpper(key))
	if proxyValue == "" {
		return os.Getenv(strings.ToLower(key))
	}
	return proxyValue
}

// DialerFromEnvironment takes in a "direct" *net.Dialer and returns a
// proxy.Dialer which will route the connections through the proxy using the
// given dialer.
func DialerFromEnvironment(direct *net.Dialer) (proxy.Dialer, error) {
	fmt.Println("get proxy dialer")
	allProxy := GetProxyEnv("all_proxy")
	fmt.Printf("allProxy = %v",allProxy)
	if len(allProxy) == 0 {
		return direct, nil
	}

	proxyURL, err := url.Parse(allProxy)
	fmt.Printf("proxyURL = %v",proxyURL)
	if err != nil {
		return direct, err
	}

	proxyFromURL, err := proxy.FromURL(proxyURL, direct)
	fmt.Printf("proxyFromURL = %v",proxyFromURL)
	if err != nil {
		return direct, err
	}

	noProxy := GetProxyEnv("no_proxy")
	fmt.Printf("noProxy = %v",noProxy)
	if len(noProxy) == 0 {
		return proxyFromURL, nil
	}

	perHost := proxy.NewPerHost(proxyFromURL, direct)
	perHost.AddFromString(noProxy)

	return perHost, nil
}
