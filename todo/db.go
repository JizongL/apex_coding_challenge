package todo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	DROP_TODOS      = "DROP TABLE TODO"
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

func TestDB() *sql.DB {
	db, _ := sql.Open("postgres",
		fmt.Sprintf("user=%s dbname=test sslmode=disable", os.Getenv("DB_USER")))
	return db
}

func OpenDB() *sql.DB {
	var user, name string
	var err error
	if user, err = USER.env(); err != nil {
		log.Fatal(err)
	} else if name, err = NAME.env(); err != nil {
		log.Fatal(err)
	}
	db, _ := sql.Open("postgres",
		fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, name))
	return db
}
