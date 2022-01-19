package structs

type Comment struct {
	ID         string
	Username   string
	Email      string
	ErrMessage string
	IsAdmin    string
	IsClient   string
	IsAgency   string
	IsTheSame  string
	Path       string
}

type Agency struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Description   string `json:"description"`
	Email         string `json:"email"`
	Profile_image string `json:"profile_image"`
}

type Trip struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Hotel        string `json:"hotel"`
	Stars        string `json:"stars"`
	Price        string `json:"price"`
	Path_img1    string `json:"img1"`
	Path_img2    string `json:"img2"`
	Path_img3    string `json:"img3"`
	Date         string `json:"date"`
	NumberOfDays string `json:"days"`
	AgencyName   string `json:"agencyName"`
	IsTheSame    string `json:"same"`
	Country      string `json:"country"`
	City         string `json:"city"`
}
