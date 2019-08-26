package apps

type KeyNotFound struct {
	key string
}

func (e KeyNotFound) Error() string {
	return e.key + " not found"

}

type KeyExpired struct {
	key string
}

func (e KeyExpired) Error() string {
	return e.key + " expired"

}

type recordConflict struct {
	key string
}

func (e recordConflict) Error() string {
	return e.key + "record conflict"

}
