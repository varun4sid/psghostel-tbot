package scraper

type AuthPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type StudentDetails struct {
	Name string `json:"name"`
}
