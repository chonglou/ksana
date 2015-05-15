package auth

import (
	"fmt"
	"github.com/chonglou/ksana"
	"time"
)

type Contact struct {
	Qq       string
	Wechat   string
	Skype    string
	Linkedin string
	Factbook string
	Logo     string
}

type User struct {
	FirstName string
	LastName  string
	Email     string
	Contact   Contact
}

type Setting struct {
	Key string
	Val string
}

type Log struct {
	Id      int
	Message string
	Created time.Time
}

type Role struct {
	Id      int
	Name    string
	Rid     int
	Rtype   string
	Created time.Time
	Updated time.Time
}

type AuthEngine struct {
}

func (ae *AuthEngine) Router(path string, r ksana.Router) {
	r.Resources(fmt.Sprintf("%s/users", path), ksana.Controller{
		Index:   []ksana.Handler{},
		Show:    []ksana.Handler{},
		New:     []ksana.Handler{},
		Create:  []ksana.Handler{},
		Edit:    []ksana.Handler{},
		Update:  []ksana.Handler{},
		Destroy: []ksana.Handler{},
	})
}

func (ae *AuthEngine) Migration(m ksana.Migration) {
	m.Add(
		"201505151405",
		`
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    CREATE TABLE users(
      id SERIAL,
      email VARCHAR,
      password BIT(64),      
      first_name VARCHAR(32) NOT NULL,
      middle_name VARCHAR(32),
      last_name VARCHAR(32) NOT NULL,
      token VARCHAR(64) NOT NULL DEFAULT UUID_GENERATE_V4(),
      provider VARCHAR(16) NOT NULL DEFAULT 'local',
      locked TIMESTAMP,
      confirmed TIMESTAMP,
      updated TIMESTAMP,
      created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
    CREATE UNIQUE INDEX users_oauth_idx ON users (provider, token);
    CREATE INDEX users_email_idx ON users (email);
    CREATE INDEX users_first_name_idx ON users (first_name);
    CREATE INDEX users_last_name_idx ON users (last_name);
    CREATE INDEX users_middle_name_idx ON users (middle_name);

    CREATE TABLE contacts(
      user_id INTEGER,
      type VARCHAR(16),
      value VARCHAR(512)
      );
    CREATE UNIQUE INDEX contacts_user_idx ON contacts (user_id,type);

    CREATE TABLE settings(
      id VARCHAR(128) NOT NULL PRIMARY KEY,
      val BIT VARYING NOT NULL,
      iv BIT(32)
      );
    CREATE TABLE logs(
      id SERIAL,
      message VARCHAR(255),
      created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
    CREATE TABLE roles(
      id SERIAL,
      name VARCHAR(32),
      r_id INTEGER,
      r_type VARCHAR,
      created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
    CREATE TABLE users_roles(
      user_id INTEGER NOT NULL,
      role_id INTEGER NOT NULL,
      PRIMARY KEY (user_id, role_id)
      );
    `,
		`
    DROP TABLE users;
    DROP TABLE contacts;
    DROP TABLE settings;
    DROP TABLE logs;
    DROP TABLE roles;
    DROP TABLE users_roles;
    `)
}
