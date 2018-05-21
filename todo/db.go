package todo

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	DROP_TABLE      = "DROP TABLE TODO"
	RECREATE_SCHEMA = `CREATE TABLE public.todo (
	id serial NOT NULL,
	title varchar NULL,
	status varchar NULL,
	CONSTRAINT todo_pk PRIMARY KEY (id)
)
WITH (
	OIDS=FALSE
);

CREATE SEQUENCE IF NOT EXISTS public.todo_id_seq
NO MINVALUE
NO MAXVALUE;`
)

//OpenTestDB opens the psql database specified by %TEST_DB_USER and $TEST_DB_NAME
func OpenTestDB() *sql.DB {
	name, err := ENV_TEST_NAME.val()
	if err != nil {
		log.Fatal(err)
	}

	user, err := ENV_TEST_USER.val()
	if err != nil {
		log.Fatal(err)
	}

	db, _ := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, name))
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

//OpenDB opens the psql database specified by $DB_USER, $DB_HOST, $DB_PASSWORD, and $DB_NAME
func OpenDB() *sql.DB {
	var user, name, host, pass string
	var err error
	if user, err = ENV_USER.val(); err != nil {
		log.Fatal(err)
	} else if name, err = ENV_NAME.val(); err != nil {
		log.Fatal(err)
	} else if host, err = ENV_HOST.val(); err != nil {
		log.Fatal(err)
	} else if pass, err = ENV_PASS.val(); err != nil {
		log.Fatal(err)
	} else if user, err = ENV_USER.val(); err != nil {
		log.Fatal(err)
	}

	db, _ := sql.Open("postgres",
		fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
			user, name, pass, host,
		),
	)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
