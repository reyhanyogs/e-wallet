package config

type Config struct {
	Server   Server
	Database Database
	Mail     Email
}

type Server struct {
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Email struct {
	User string
	API  string
}
