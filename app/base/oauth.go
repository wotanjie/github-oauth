package base

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	RequestUri = "https://github.com/login/oauth/authorize"
	BaseUri = "https://github.com/login/oauth/access_token"
	ClientId = "0caea0cb11135cc31419"
	ClientSecrect = "a87752dfaa87200feff508b54cac832e3501fdad"
)

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
}

func Auth(code string)(OAuthAccessResponse,error){
	url := fmt.Sprintf(BaseUri+"?client_id=%s&client_secret=%s&code=%s",ClientId,ClientSecrect,code)
	header := make(map[string]string)
	header["accept"] = "application/json"
	res,err := Reuqest("POST",url,nil,header)
	defer res.Body.Close()
	var t OAuthAccessResponse
	if err =json.NewDecoder(res.Body).Decode(&t) ;err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		return OAuthAccessResponse{AccessToken:""},err
	}else {
		return t,nil
	}
}

func CheckLogin(r *http.Request) (*OAuthAccessResponse,bool) {
	session,err := SessionStore.Get(r,"user")
	if err != nil {
		return nil,false
	}
	user := session.Values["auth"]
	auth ,ok := user.(*OAuthAccessResponse)
	return auth, user != nil && ok
}
