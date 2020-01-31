package base

import (
	"io"
	"net/http"
)

var httpClient = http.Client{}

func Reuqest(method string,url string,body io.Reader,header map[string]string)(*http.Response,error){
	req,err := http.NewRequest(method,url,body)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k,v := range(header)  {
			req.Header.Set(k,v)
		}
	}
	return httpClient.Do(req)
}