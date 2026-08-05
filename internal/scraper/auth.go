package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type AuthPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type StudentDetails struct {
	Name string `json:"name"`
}

const (
	index_page_url   = "https://edviewx.psgtech.ac.in/Hostel/Home/Index"
	auth_url         = "https://edviewx.psgtech.ac.in/Hostel/Login/Authenticate"
	student_info_url = "https://edviewx.psgtech.ac.in/Hostel/Student/studDetails"
	get_token_url    = "https://edviewx.psgtech.ac.in/Hostel/Student/StudentGetToken"
	qrCode_url       = "https://edviewx.psgtech.ac.in/Hostel/QRCode/QRcodeGenerate"
)

func newAuthClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{Jar: jar}, nil
}

func GetAuthenticatedClient(rollno, password string) (*http.Client, error) {
	client, err := newAuthClient()
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO CREATE HTTP CLIENT: %w", err)
	}

	resp, err := client.Get(index_page_url)
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO GET INDEX PAGE: %w", err)
	}
	defer resp.Body.Close()

	response, err := client.PostForm(auth_url, url.Values{
		"name":     {rollno},
		"password": {password},
	})
	if err != nil || response.StatusCode >= 400 {
		return nil, fmt.Errorf("UNABLE TO AUTHENTICATE USER: %w", err)
	}
	defer response.Body.Close()

	return client, nil
}

func GetUserIfExists(rollno, password string) (string, error) {
	var err error

	client, err := GetAuthenticatedClient(rollno, password)
	if err != nil {
		return "", err
	}

	response, err := client.Get(student_info_url + "?rollno=" + rollno)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO FETCH STUDENT INFO: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("UNABLE TO FETCH STUDENT INFO: %w", err)
	}

	var studentDetails []StudentDetails
	err = json.NewDecoder(response.Body).Decode(&studentDetails)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO DECODE STUDENT INFO: %w", err)
	}

	return studentDetails[0].Name, nil
}
