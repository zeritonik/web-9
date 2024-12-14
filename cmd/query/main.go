package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const connectionString = "host=localhost port=5432 user=postgres dbname=sandbox password=postgres"

type Handlers struct {
	db *sql.DB
}

func (h *Handlers) ServePost(c echo.Context) error {
	var data struct {
		Name string `json:name`
	}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err = h.db.Exec("UPDATE query_table SET NAME = $1", data.Name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK")
}

func (h *Handlers) ServeGet(c echo.Context) error {
	var name string
	row := h.db.QueryRow("SELECT name FROM query_table LIMIT 1")
	err := row.Scan(&name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	name = "Hello, " + name
	return c.String(http.StatusOK, name)
}

func main() {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	h := Handlers{db: db}

	e := echo.New()
	e.GET("/get", h.ServeGet)
	e.POST("/post", h.ServePost)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
