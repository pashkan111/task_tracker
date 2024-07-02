package entities

type User struct {
	Id             int
	PassportSerie  int
	PassportNumber int
	Surname        string
	Name           string
}

type UserCreateRequest struct {
	PassportNumber string `json:"passportNumber" validate:"required"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
}

type UserCreateResponse struct {
	Id             int    `json:"id"`
	PassportSerie  int    `json:"passportSerie"`
	PassportNumber int    `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
}

type UserUpdateRequest struct {
	PassportSerie  int    `json:"passportSerie"`
	PassportNumber int    `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
}

type UserUpdateResponse struct {
	Id             int    `json:"id"`
	PassportSerie  int    `json:"passportSerie"`
	PassportNumber int    `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
}
