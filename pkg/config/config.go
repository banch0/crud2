package config

// Config ...
type Config struct {
	Addr   string `json:"addr"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	DBUser string `json:"dbuser"`
	DBPass string `json:"dbpass"`
	DBPort string `json:"dbport"`
	DBHost string `json:"dbhost"`
	DBName string `json:"dbname"`
}