package ksana_orm

func New(path string, cfg *Config) (Model, error) {
	db := Connection{}

	err := db.Open(path, cfg)
	if err != nil {
		return nil, err
	}

	return &model{db: &db}, nil
}
