package ksana

type databaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Ssl      string `json:"ssl"`
}



func NewOrm(path string, cfg *Config) (Model, error) {
	db := Connection{}

	err := db.Open(path, cfg)
	if err != nil {
		return nil, err
	}

	return &model{db: &db}, nil
}
