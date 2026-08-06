package scraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	index_page_url   = "https://edviewx.psgtech.ac.in/Hostel/Home/Index"
	auth_url         = "https://edviewx.psgtech.ac.in/Hostel/Login/Authenticate"
	student_info_url = "https://edviewx.psgtech.ac.in/Hostel/Student/studDetails"
	get_token_url    = "https://edviewx.psgtech.ac.in/Hostel/Student/StudentGetToken"
	qrCode_url       = "https://edviewx.psgtech.ac.in/Hostel/QRCode/QRcodeGenerate"
)

func GetLiveToken(rollno, password string) (string, error) {
	client, err := GetAuthenticatedClient(rollno, password)
	if err != nil {
		return "", err
	}

	response, err := client.Get(qrCode_url)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO FETCH QR CODE PAGE : %w", err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO PARSE QR CODE PAGE : %w", err)
	}

	qrCodeEncoded, _ := doc.Find("img[alt='QR Code']").Attr("src")

	return qrCodeEncoded, nil
}
