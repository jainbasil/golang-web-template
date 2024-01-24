package services

// Service provides convenient functions related to the business logic that
// can be invoked from other ends like controllers or message queue handlers etc.
type Service struct{}

func NewService() *Service {
	return &Service{}
}
