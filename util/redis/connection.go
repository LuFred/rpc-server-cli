package redis

// connectionURL implements a redis connection struct.
type connectionURL struct {
	Password string
	Database int
	Host     string
	Options  map[string]string
}
