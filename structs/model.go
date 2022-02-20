package structs

type Comment struct {
	ID           string
	ID2          string
	Username     string
	Email        string
	ErrMessage   string
	IsAdmin      string
	IsClient     string
	IsAgency     string
	IsTheSame    string
	Path         string
	MoneyBalance string
}

type Agency struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Description   string `json:"description"`
	Email         string `json:"email"`
	Profile_image string `json:"profile_image"`
	Is_admin      string `json:"is_admin"`
}

type Client struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	City      string `json:"city"`
	Hotel     string `json:"hotel"`
	Date      string `json:"date"`
	PathToPDF string `json:"path_to_pdf"`
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
	IsClient     string `json:"is_client"`
	ClientID     string `json:"clientId"`
	TripID       string `json:"tripId"`
	AgencyID     string `json:"agencyId"`
	PathToPDF    string `json:"path_to_pdf"`
}

type Review struct {
	Client    string `json:"client"`
	Title     string `json:"title"`
	Comment   string `json:"comment"`
	Stars     string `json:"stars"`
	Date      string `json:"date"`
	IsTheSame string `json:"same"`
}
