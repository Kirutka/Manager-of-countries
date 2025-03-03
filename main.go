package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Country struct {
	ID             int
	Name           string
	Flag           string
	Territory      string
	FoundationYear int
	Currency       string
}

func main() {
	db, err := sql.Open("postgres", "user=postgres password=ваш-пароль host=localhost dbname=users sslmode=disable") // подключение к базе данных
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*") // подгружаем html

	r.GET("/", func(c *gin.Context) { // главная страница
		handleGetCountries(c, db)
	})

	r.POST("/delete", func(c *gin.Context) { // удаление
		handleDeleteCountry(c, db)
	})

	r.POST("/add", func(c *gin.Context) { // добавление
		handleAddCountry(c, db)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func handleGetCountries(c *gin.Context, db *sql.DB) { // функция получения страны
	countries, err := getCountries(db)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Ошибка получения данных", err)
		return
	}
	c.HTML(http.StatusOK, "index.html", countries)
}

func handleDeleteCountry(c *gin.Context, db *sql.DB) { // функция удаления страны
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Неверный ID", err)
		return
	}

	if err := deleteCountry(db, id); err != nil {
		handleError(c, http.StatusInternalServerError, "Ошибка удаления страны", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func handleAddCountry(c *gin.Context, db *sql.DB) { // добавление данных из формы на html
	name := c.PostForm("name")
	flag := c.PostForm("flag")
	territory := c.PostForm("territory")
	foundationYearStr := c.PostForm("foundation_year")
	currency := c.PostForm("currency")

	foundationYear, err := strconv.Atoi(foundationYearStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Неверный год основания", err)
		return
	}

	newCountry := Country{
		Name:           name,
		Flag:           flag,
		Territory:      territory,
		FoundationYear: foundationYear,
		Currency:       currency,
	}

	if err := addCountry(db, newCountry); err != nil {
		handleError(c, http.StatusInternalServerError, "Ошибка добавления страны", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func handleError(c *gin.Context, status int, message string, err error) { // функция для ошибок
	log.Printf("%s: %v", message, err)
	c.String(status, message)
}

func getCountries(db *sql.DB) ([]Country, error) { // функция с бд на получение стран
	rows, err := db.Query("SELECT id, name, flag, territory, foundation_year, currency FROM country")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []Country

	for rows.Next() {
		var country Country
		if err := rows.Scan(&country.ID, &country.Name, &country.Flag, &country.Territory, &country.FoundationYear, &country.Currency); err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, rows.Err()
}

func deleteCountry(db *sql.DB, id int) error { // функция с бд на удаление страны
	_, err := db.Exec("DELETE FROM country WHERE id = $1", id)
	return err
}

func addCountry(db *sql.DB, country Country) error { // функция с бд на добавление страны
	_, err := db.Exec("INSERT INTO country (name, flag, territory, foundation_year, currency) VALUES ($1, $2, $3, $4, $5)",
		country.Name, country.Flag, country.Territory, country.FoundationYear, country.Currency)
	return err
}
