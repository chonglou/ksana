package ksana

type Migration interface {
	Add(Bean) error
	Migrate() error
	Rollback() error
}

type migration struct {
}

func (m *migration) Add( b Bean) error{
	return nil
}

func (m *migration) Migrate() error {
	return nil
}

func (m *migration) Rollback() error {
	return nil
}
