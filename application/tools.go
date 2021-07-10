package application

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	hclis            map[string]*http.Transport
	defaultTransport *http.Transport
)

func Check(err *error) {
	if *err != nil {
		Logger.Error(*err)
	}
}

func getHcli(proxyname string) (hcli *http.Client, err error) {
	//defer func() {
	//	logs.Debug(hcli.Transport)
	//}()
	if proxyname == `` {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hcli = &http.Client{Transport: tr}
	} else {
		if transport, ok := hclis[proxyname]; !ok {
			err = fmt.Errorf(`proxy name not found`)
			return
		} else {
			client := &http.Client{
				Transport: transport,
			}
			hcli = client
		}
	}
	return
}

func SendRequest(url string) ([]byte, error){
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//func SendRequest(method, url string, headers map[string]string, data []byte, proxyname string) (respbody []byte, err error) {
//	defer Check(&err)
//	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
//	if err != nil {
//		return
//	}
//	for key, value := range headers {
//		req.Header.Add(key, value)
//	}
//	hc, err := getHcli(proxyname)
//	if err != nil {
//		return
//	}
//	if v, ok := headers[`host`]; ok {
//		req.Host = v
//	}
//	fmt.Println(fmt.Sprintf("before Do, request: %v", req))
//	resp, err := hc.Do(req)
//	fmt.Println("after Do")
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	defer resp.Body.Close()
//	var body []byte
//	if resp != nil {
//		body, err = ioutil.ReadAll(resp.Body)
//		if err != nil {
//			return
//		}
//	}
//	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
//		respbody = body
//		return
//	}
//	if method == http.MethodDelete && resp.StatusCode == http.StatusNoContent {
//		respbody = body
//		return
//	}
//	err = fmt.Errorf("[%d]%s", resp.StatusCode, body)
//	return
//}
