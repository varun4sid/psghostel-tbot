package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

func decodeStudentDetails(r io.Reader) (string, error) {
	var studentDetails []StudentDetails
	if err := json.NewDecoder(r).Decode(&studentDetails); err != nil {
		return "", fmt.Errorf("UNABLE TO DECODE STUDENT INFO: %w", err)
	}

	if len(studentDetails) == 0 {
		return "", nil
	}

	return studentDetails[0].Name, nil
}

func GetUserIfExists(rollno, password string) (string, error) {
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

	return decodeStudentDetails(response.Body)
}
