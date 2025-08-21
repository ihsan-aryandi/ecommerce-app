package repository

func handleNullAndError[T any](t *T, found bool, err error) (*T, error) {
	if err != nil {
		return nil, err
	}

	if !found {
		t = nil
	}

	return t, nil
}
