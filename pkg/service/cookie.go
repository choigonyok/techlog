package service

func (svc *Service) SetCookieValueByUniqueID(uniqueID string) error {
	return svc.provider.SetCookieValueByUniqueID(uniqueID)
}

func (svc *Service) UpdateCookieValueByUniqueID(uniqueID string) error {
	return svc.provider.UpdateCookieValueByUniqueID(uniqueID)
}

func (svc *Service) GetCookieValue() (string, error) {
	return svc.provider.GetCookieValue()
}

func (svc *Service) VerifyAdminByCookieValue(value string) (bool, error) {
	storedValue, err := svc.provider.GetCookieValue()
	if storedValue == value {
		return true, nil
	} else {
		return false, err
	}
}
