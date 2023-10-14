package test

import (
	"database/sql"
	"fmt"
	"graha/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {

	accounts, err := account_list(c.Get("db").(*sql.DB))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("Database error scan %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, accounts)
}

func TestAuth(c echo.Context) error {
	return c.String(http.StatusOK, "Authorized!")
}

func account_list(db *sql.DB) ([]models.Account, error) {

	query := `WITH RECURSIVE tree(root, id) AS (
		SELECT root,
 					 id,
					 name,
					 name_en,
					 description,
					 is_active
 		  FROM accounts
	   WHERE root = 0
 UNION ALL
	  SELECT a.root,
           a.id,
					 a.name,
					 a.name_en,
					 a.description,
					 a.is_active
			FROM accounts a
INNER JOIN tree t
				ON a.root = t.id
)
    SELECT root, id, name, name_en, description, is_active
		  FROM tree
  ORDER BY root, id;`

	q, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer q.Close()
	accounts := make([]models.Account, 0)
	// defer rs.Close()

	for q.Next() {
		var p models.Account
		err = q.Scan(
			&p.Root,
			&p.ID,
			&p.Name,
			&p.Name_en,
			&p.Description,
			&p.Is_active,
		)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		accounts = append(accounts, p)
	}

	return accounts, nil
}
