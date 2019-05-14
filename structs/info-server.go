package infoserver

// Server info about server associated with domain
type Server struct {
	IPAddress string
	Address   string
	SslGrade  string
	Country   string
	Owner     string
}

// InfoServer info about server configuration with domain and know if was change
type InfoServer struct {
	Servers          []Server
	ServersChanged   bool
	SslGrade         string
	PreviousSslGrade string
	Logo             string
	IsDown           bool
}
