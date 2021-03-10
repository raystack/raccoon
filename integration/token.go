package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchOTP(phNumber string) (string, error) {
	url := "https://integration-goid.gojekapi.com/goid/login/request"
	b, _ := json.Marshal(map[string]string{
		"client_id":     "gojek:consumer:app",
		"client_secret": "xKMPsxLFkMVlpZPRojqPwLl54X1Qch",
		"phone_number":  phNumber,
		"country_code":  "+62",
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Add("accept-language", "en")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-appid", "com.gojek.app")
	req.Header.Add("x-appversion", "3.53.0")
	req.Header.Add("x-deviceos", "Android 27")
	req.Header.Add("x-devicetoken", "test4")
	req.Header.Add("x-imei", "1")
	req.Header.Add("x-phonemake", "Xiaomi")
	req.Header.Add("x-phonemodel", "Redmi 5A")
	req.Header.Add("x-pushtokentype", "FCM")
	req.Header.Add("x-request-id", "5e770d9b-2966-4ebf-9f56-14e56892c908")
	req.Header.Add("x-uniqueid", "8b0df26d987081f9")
	req.Header.Add("x-user-type", "customer")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("something went wrong when sending the request: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("something went wrong when reading response: %v", err)
	}
	resp.Body.Close()
	return getToken(body)
}

func getToken(resp []byte) (string, error) {
	type data struct {
		OtpToken string `json:"otp_token"`
	}
	type Body struct {
		Data data `json:"data"`
	}

	b := Body{}
	err := json.Unmarshal(resp, &b)
	if err != nil {
		return "", fmt.Errorf("something went wrong when unmarshalling response: , %v", err)
	}
	return b.Data.OtpToken, nil
}

func login(phNumber string) (auth Auth, err error) {
	url := fmt.Sprintf("http://10.120.4.21:9000/otp?phone=%s&env=integration", phNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("clientid", "red_robin")
	req.Header.Add("passkey", "robin1234")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return auth, fmt.Errorf("request err: %v", err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("not 200!")
		return auth, fmt.Errorf("status code not 200 when login")
	}

	type Body struct {
		Data []Auth `json:"data"`
	}
	b := &Body{}
	r, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(r, b)
	if err != nil {
		return auth, err
	}
	return b.Data[len(b.Data)-1], nil
}

type Auth struct {
	Otp      string `json:"otp"`
	OtpToken string `json:"otp_token"`
}

func fetchAccessToken(auth Auth) (string, error) {
	type ReqBody struct {
		ClientID     string            `json:"client_id"`
		ClientSecret string            `json:"client_secret"`
		GrantType    string            `json:"grant_type"`
		Data         map[string]string `json:"data"`
	}
	reqBody := ReqBody{
		ClientID:     "gojek:consumer:app",
		ClientSecret: "xKMPsxLFkMVlpZPRojqPwLl54X1Qch",
		GrantType:    "otp",
		Data: map[string]string{
			"otp":       auth.Otp,
			"otp_token": auth.OtpToken,
		},
	}

	b, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://integration-goid.gojekapi.com/goid/token", bytes.NewBuffer(b))
	req.Header.Add("accept-language", "en")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-appid", "com.gojek.app")
	req.Header.Add("x-appversion", "3.53.0")
	req.Header.Add("x-deviceos", "Android 27")
	req.Header.Add("x-devicetoken", "test4")
	req.Header.Add("x-imei", "1")
	req.Header.Add("x-phonemake", "Xiaomi")
	req.Header.Add("x-phonemodel", "Redmi 5A")
	req.Header.Add("x-pushtokentype", "FCM")
	req.Header.Add("x-request-id", "5e770d9b-2966-4ebf-9f56-14e56892c908")
	req.Header.Add("x-uniqueid", "8b0df26d987081f9")
	req.Header.Add("x-user-type", "customer")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 201 && res.StatusCode != 200 {
		return "", fmt.Errorf("status code: %v", res.StatusCode)
	}

	r, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("fail to read body. %v", err)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(r, &m)
	if err != nil {
		return "", fmt.Errorf("fail to read body. %v", err)
	}
	return m["access_token"].(string), nil
}

func FetchAccessToken() (string, error) {
	ph := "820040000101"
	otpToken, err := fetchOTP(ph)
	if err != nil {
		return "", fmt.Errorf("fail to fetch OTP. %v", err)
	}
	a, err := login(ph)
	if err != nil {
		return "", fmt.Errorf("fail to get credential. %v", err)
	}
	token, err := fetchAccessToken(Auth{
		Otp:      a.Otp,
		OtpToken: otpToken,
	})
	if err != nil {
		return "", fmt.Errorf("fail to fetch access token. %v", err)
	}
	return token, nil
}
