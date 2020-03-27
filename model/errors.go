package model

type ApiError struct {
	StatusCode int
	Text       string
}

func (a *ApiError) Error() string {
	return a.Text
}
