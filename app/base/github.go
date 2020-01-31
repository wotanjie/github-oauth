package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"github-oauth/app/models"
	"net/http"
)

var (
	sessionStore = SessionStore
)

func FetchUserInfo(r *http.Request,w http.ResponseWriter)(*models.User,error) {
	userData, isLogin := CheckLogin(r)
	if !isLogin {
		return nil,errors.New("unlogin")
	}
	accessToken := userData.AccessToken
	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("token %s", accessToken)
	header["accept"] = "application/json"
	resp,err := Reuqest("GET", "https://api.github.com/user", nil,header)
	if err != nil{
		return nil, err
	}
	//未授权
	if resp.StatusCode == 401 {
		session,err := sessionStore.Get(r,"user")
		if err != nil {
			return nil, err
		}
		session.Options.MaxAge = -1
		_ = session.Save(r,w)
		return nil,errors.New("Bad credentials")
	}
	defer resp.Body.Close()
	var user models.User
	_ = json.NewDecoder(resp.Body).Decode(&user)
	return &user,nil
}
