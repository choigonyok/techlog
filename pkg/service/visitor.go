package service

// GetDate returns stored visitor/date
func (svc *Service) GetDate() (string, error) {
	m, err := svc.provider.GetVisitor()
	return m.Date, err
}

// ResetToday resets stored visitor/today to 1 and visitor/date to today
func (svc *Service) ResetToday(today string) error {
	return svc.provider.ResetVisitorTodayAndDate(today)
}

// GetCounts returns stored visitor/today and visitor/total
func (svc *Service) GetCounts() (int, int, error) {
	m, err := svc.provider.GetVisitor()
	return m.Today, m.Total, err
}

// AddToday updates stored visitor/today, visitor/total to visitor/today + 1, visitor/total + 1
func (svc *Service) AddTodayAndTotal() error {
	m, err := svc.provider.GetVisitor()
	newToday := m.Today + 1
	newTotal := m.Total + 1
	if err != nil {
		return err
	}
	return svc.provider.UpdateVisitorToday(newToday, newTotal)
}
