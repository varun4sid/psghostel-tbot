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

func GetLiveToken(rollno, password string) (string, string, error) {
	client, err := GetAuthenticatedClient(rollno, password)
	if err != nil {
		return "", "", err
	}

	response, err := client.Get(qrCode_url)
	if err != nil {
		return "", "", fmt.Errorf("UNABLE TO FETCH QR CODE PAGE : %w", err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("UNABLE TO PARSE QR CODE PAGE : %w", err)
	}

	qrCodeEncoded, _ := doc.Find("img[alt='QR Code']").Attr("src")

	var tokens []Token
	var tokenDate string
	var caption string
	doc.Find("table.card-table tbody tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() >= 6 {
			tokenName := cells.Eq(2).Text()
			tokenQty := cells.Eq(5).Text()
			tokenDate = cells.Eq(3).Text()

			tokens = append(tokens, Token{
				Name:     tokenName,
				Quantity: tokenQty,
			})
		}

		caption = fmt.Sprintf("DATE : %s\n", tokenDate)
		for _, token := range tokens {
			caption += fmt.Sprintf("%s x%s\n", token.Name, token.Quantity)
		}
	})

	return qrCodeEncoded, caption, nil
}
