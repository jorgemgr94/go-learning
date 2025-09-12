package config

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}
