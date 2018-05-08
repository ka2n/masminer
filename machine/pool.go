package machine

// Pool common formatted pool
type Pool struct {
	URL  string
	Pass string
	User string
	Algo string

	Options map[string]string
}
