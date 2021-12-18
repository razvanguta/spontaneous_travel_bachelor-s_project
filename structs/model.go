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
