package common

type JobNotification struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobNotification(m, e string) *JobNotification {
	return &JobNotification{
		Message: m,
		Error:   e,
	}
}
