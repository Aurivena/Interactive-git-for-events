package entity

type AppError struct {
	Code   int         `json:"code"`
	Detail ErrorDetail `json:"detail"`
}

type ErrorDetail struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Solution string `json:"solution"`
}
