package auth

type JWTConfig struct {
	ServerHost      string        `yaml:"server_host"`
	DurationInHours int           `yaml:"duration_in_hours"`
	Credentials     JWTCredential `yaml:"jwt_credentials"`
}

type JWTCredential struct {
	Secret string `yaml:"secret"`
}

type UserJWTPayload struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
