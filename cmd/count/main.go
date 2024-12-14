package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const connectionString = "host=localhost port=5432 user=postgres dbname=sandbox password=postgres"

type Handlers struct {
	db *sql.DB
}

func (h *Handlers) ServeGet(c echo.Context) error {
	var count int
	err := h.db.QueryRow("SELECT count FROM count_table LIMIT 1").Scan(&count)
	var res struct {
		Count int    `json:"count"`
		Err   string `json:"error"`
	}
	if err != nil {
		res.Err = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}
	res.Count = count
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) ServePost(c echo.Context) error {
	var dcount struct {
		Count int `json:"count"`
	}
	var res struct {
		Err string `json:"error"`
	}
	err := json.NewDecoder(c.Request().Body).Decode(&dcount)
	if err != nil {
		res.Err = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	_, err = h.db.Exec("UPDATE count_table SET count = count + $1", dcount.Count)
	if err != nil {
		res.Err = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}
	return c.JSON(http.StatusOK, res)
}

func main() {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return
	}
	defer db.Close()

	h := Handlers{db: db}

	e := echo.New()
	e.GET("/get", h.ServeGet)
	e.POST("/post", h.ServePost)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
