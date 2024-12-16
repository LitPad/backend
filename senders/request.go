package senders

import (
	"time"

	"github.com/valyala/fasthttp"
)

var httpClient = fasthttp.Client{
	ReadTimeout: 30 * time.Second,
	WriteTimeout: 30 * time.Second,
}

func MakeRequest(method, url string, headers map[string]string, body []byte)(*fasthttp.Response, error){
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	for key, value := range headers{
		req.Header.Set(key, value)
	}


	if method == fasthttp.MethodPost && body != nil{
		req.Header.Set("Content-Type", "application/json") 
		req.SetBody(body)
	}

	if err := httpClient.Do(req, res); err !=nil{
		return nil, err
	}


	responseCopy := fasthttp.AcquireResponse()
	res.CopyTo(responseCopy)

	return responseCopy, nil
}