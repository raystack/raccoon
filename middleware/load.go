package middleware

var loaded bool

func Load() {
	if loaded {
		return
	}
	loadCors()
}
