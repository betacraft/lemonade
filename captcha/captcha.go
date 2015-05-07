package captcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Validate(response string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	params := req.URL.Query()
	params.Add("secret", "6LefTAYTAAAAAODl3wracNKxSIGndmfa1nl7kycB")
	params.Add("response", response)
	req.URL.RawQuery = params.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var data map[string]interface{}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return false
	}
	success, ok := data["success"].(bool)
	if !ok {
		return false
	}
	return success
}
