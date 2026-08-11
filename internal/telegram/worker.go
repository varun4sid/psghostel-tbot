package telegram

import (
	"database/sql"
	"encoding/base64"
	"log"
	"os"
	"psghostelbot/internal/crypt"
	"psghostelbot/internal/db"
	"psghostelbot/internal/scraper"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
)

func RunBatchScrape(database *sql.DB, psg *TelegramBot, mealName string) {
	students, err := db.GetAllUsers(database)
	if err != nil {
		log.Printf("UNABLE TO QUERY USERS : %v", err)
		return
	}

	if len(students) == 0 {
		log.Print("No users registered to scrape...")
		return
	}

	var QRcodesSent int64
	var wg sync.WaitGroup
	for _, student := range students {
		wg.Add(1)

		go func(s db.Student) {
			defer wg.Done()

			plainPassword, err := crypt.DecryptPassword(s.EncPassword, os.Getenv("AES_KEY"))
			if err != nil {
				log.Printf("UNABLE TO DECRYPT PASSWORD FOR ROLL NO %s : %v", s.Rollno, err)
				return
			}

			err = RunScraperForUser(s.Rollno, plainPassword, s.ChatID, psg)
			if err != nil {
				log.Printf("UNABLE TO SEND PHOTO TO CHAT ID %d FOR ROLL NO %s : %v", s.ChatID, s.Rollno, err)
				return
			} else {
				atomic.AddInt64(&QRcodesSent, 1)
			}

		}(student)
	}
	wg.Wait()

	log.Printf("%s session batch scrape count : %d", mealName, len(students))
	log.Printf("Successfully sent QR codes count : %d", atomic.LoadInt64(&QRcodesSent))
}

func CreateScheduler(database *sql.DB, psg *TelegramBot) *gocron.Scheduler {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*60*60+30*60) // Fallback to IST if timezone loading fail
	}

	scheduler := gocron.NewScheduler(loc)

	scheduler.Every(1).Day().At("07:00").Do(
		RunBatchScrape, database, psg, "Breakfast",
	)

	scheduler.Every(1).Day().At("11:30").Do(
		RunBatchScrape, database, psg, "Lunch",
	)

	scheduler.Every(1).Day().At("18:30").Do(
		RunBatchScrape, database, psg, "Dinner",
	)

	return scheduler
}

func GetQRImage(rollno string, password string, psg *TelegramBot) ([]byte, string) {
	qrCodeEncoded, caption, err := scraper.GetLiveToken(rollno, password)
	if err != nil {
		log.Printf("UNABLE TO GET LIVE TOKEN FOR ROLL NO %s : %v", rollno, err)
		return nil, ""
	}

	if qrCodeEncoded == "" {
		return nil, ""
	}

	var rawBase64 string = qrCodeEncoded
	if strings.Contains(rawBase64, "data:image/png;base64,") {
		rawBase64 = strings.TrimPrefix(rawBase64, "data:image/png;base64,")
	}

	imageBytes, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		log.Printf("UNABLE TO DECODE BASE64 IMAGE FOR ROLL NO %s : %v", rollno, err)
		return nil, ""
	}

	return imageBytes, caption
}

func RunScraperForUser(rollno string, password string, chatID int64, psg *TelegramBot) error {
	imageBytes, caption := GetQRImage(rollno, password, psg)
	if imageBytes == nil {
		log.Printf("UNABLE TO GET QR CODE FOR ROLL NO %s", rollno)
	}

	err := psg.sendPhoto(chatID, imageBytes, caption)
	return err
}
