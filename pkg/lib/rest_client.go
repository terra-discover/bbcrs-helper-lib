package lib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RestClient struct {
	URL     string
	Method  string
	Timeout int
	Headers map[string]string
	Request interface{}
}

func (r *RestClient) SetURL(url string) *RestClient {
	r.URL = url
	return r
}

func (r *RestClient) SetMethod(method string) *RestClient {
	allowedMethods := []string{"get", "post", "delete", "patch", "put", "options", "head"}
	if _, hasSlice := FindSlice(allowedMethods, strings.ToLower(method)); !hasSlice {
		method = "get"
	}
	r.Method = method
	return r
}

func (r *RestClient) SetTimeout(timeout int) *RestClient {
	if timeout <= 1 {
		timeout = 15 //static default timeout.
	}
	r.Timeout = timeout
	return r
}

func (r *RestClient) SetHeaders(headers map[string]string) *RestClient {
	for hname, hval := range headers {
		r = r.AddHeader(hname, hval)
	}
	return r
}

func (r *RestClient) AddHeader(headerName, headerValue string) *RestClient {
	if len(r.Headers) == 0 {
		r.Headers = make(map[string]string)
	}
	r.Headers[headerName] = headerValue
	return r
}

func (r *RestClient) SetRequest(param interface{}) *RestClient {
	r.Request = param
	return r
}

func (r *RestClient) Execute() (httpBody string, httpStatus int) {
	var restrequest []byte
	restRequestID, _ := uuid.NewRandom()

	dataType := fmt.Sprintf("%T", r.Request)

	if strings.Contains(dataType, "string") {
		restrequest = []byte(r.Request.(string))
	} else {
		restrequest, _ = JSONMarshal(r.Request)
	}

	if len(r.Method) == 0 {
		r.Method = "GET"
	}
	if r.Timeout == 0 {
		r.Timeout = 15
	}

	// create request structure
	req, _ := http.NewRequest(r.Method, r.URL, strings.NewReader(string(restrequest)))

	for hname, hval := range r.Headers {
		req.Header.Set(hname, hval)
	}
	req.Close = true // this is required to prevent too many files open

	// Create HTTP Connection
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(r.Timeout) * time.Second,
	}

	// Now hit to destionation endpoint
	LogStruct(map[string]interface{}{
		"REQUESTID": restRequestID.String(),
		"URL":       r.URL,
		"METHOD":    r.Method,
		"REQUEST":   string(restrequest),
	}, "RESTCLIENT REQUEST LOG")
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Call URL Failed : %s", err.Error())
		httpStatus = 0
		httpBody = "Call URL Failed : " + err.Error()
		// if res != nil {
		// 	buff := new(bytes.Buffer)
		// 	buff.ReadFrom(res.Body)
		// 	httpBody = buff.String()
		// 	httpStatus = res.StatusCode
		// 	log.Printf("Body : %s", httpBody)
		// } else {
		// 	httpBody = "Call URL Failed : " + err.Error()
		// }
		return
	}
	defer res.Body.Close()

	buff := new(bytes.Buffer)
	buff.ReadFrom(res.Body)
	httpBody = buff.String()
	httpStatus = res.StatusCode

	LogStruct(map[string]interface{}{
		"REQUESTID": restRequestID.String(),
		"RESPONSE":  httpBody,
	}, "REST CLIENT RESPONSE LOG")
	return
}

// FindSlice Find string on slice
func FindSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
