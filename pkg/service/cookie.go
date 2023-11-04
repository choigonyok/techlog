package service

func (svc *Service) SetNewCookieValueByUniqueID(uniqueID string) error {
	return svc.provider.SetNewCookieValueByUniqueID(uniqueID)
}

func (svc *Service) VerifyAdminByCookieValue(value string) (bool, error) {
	storedValue, err := svc.provider.GetCookieValue()
	if storedValue == value {
		return true, nil
	} else {
		return false, err
	}
}
