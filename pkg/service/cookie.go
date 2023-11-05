package service

// SetCookieValueByUniqueID stores new cookie value into database
func (svc *Service) SetCookieValueByUniqueID(uniqueID string) error {
	return svc.provider.SetCookieValueByUniqueID(uniqueID)
}

// UpdateCookieValueByUniqueID updates stored cookie value
func (svc *Service) UpdateCookieValueByUniqueID(uniqueID string) error {
	return svc.provider.UpdateCookieValueByUniqueID(uniqueID)
}

// GetCookieValue returns stored cookie value from database
func (svc *Service) GetCookieValue() (string, error) {
	return svc.provider.GetCookieValue()
}

// VerifyAdminByCookieValue compares client cookie value with stored cookie value from database
func (svc *Service) VerifyAdminByCookieValue(value string) (bool, error) {
	storedValue, err := svc.provider.GetCookieValue()
	if storedValue == value {
		return true, nil
	} else {
		return false, err
	}
}
