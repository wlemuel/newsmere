package domain

type Config struct {
	Debug    bool      `json:"debug"`
	Backends []Backend `json:"backends"`
	Services []Service `json:"services"`
}
