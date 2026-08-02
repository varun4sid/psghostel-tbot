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
)

func newAuthClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{Jar: jar}, nil
}

func CheckValidCredentials(rollno, password string) (bool, string, error) {
	var username string

	client, err := newAuthClient()
	if err != nil {
		return false, username, err
	}

	resp, err := client.Get(index_page_url)
	if err != nil {
		return false, username, fmt.Errorf("error fetching index page : %v", err)
	}
	defer resp.Body.Close()

	response, err := client.PostForm(auth_url, url.Values{
		"name":     {rollno},
		"password": {password},
	})
	if err != nil {
		return false, username, fmt.Errorf("%s : error during authentication : %v", rollno, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return false, username, fmt.Errorf("%s : authentication failed", rollno)
	}

	response, err = client.Get(student_info_url + "?rollno=" + rollno)
	if err != nil {
		return false, username, fmt.Errorf("%s : failed to fetch student info : %v", rollno, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return false, username, fmt.Errorf("%s : failed to fetch student info : %v", rollno, err)
	}

	var studentDetails []StudentDetails
	err = json.NewDecoder(response.Body).Decode(&studentDetails)
	if err != nil {
		return false, username, fmt.Errorf("%s : failed to decode student info : %v", rollno, err)
	}

	username = studentDetails[0].Name
	return true, username, nil
}
