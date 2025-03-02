package user

type ListUserRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Ids    []string
}

type GetUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}
