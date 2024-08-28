package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Signum   string `json:"signum"`
	Role     string `json:"role"`
	ProgramName string `json:"programname"`
}

func RegisterUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(User)

		if err := c.Bind(user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if user.Signum == "" || user.ProgramName == ""{
			return echo.NewHTTPError(http.StatusBadRequest, "missing or invalid parameters")
		}

		// Check if the user already exists in the database
		err := db.QueryRow("SELECT signum FROM user_management WHERE signum = $1", user.Signum).Scan(&user.Signum)
		if err == nil {
			//return echo.NewHTTPError(http.StatusConflict, "User already registered")
			return c.JSONPretty(http.StatusConflict, map[string]interface{}{
				"message": "User already registered",
			}, " ")
		} else if err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO user_management (signum, role, programname) VALUES ($1, $2, $3)", user.Signum, user.Role, user.ProgramName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"message": "User Registered Successfully",
		}, " ")
	}
}

func GetUserRole(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		signum := c.Param("signum")

		user := User{}
		err := db.QueryRow("SELECT role, programname FROM user_management WHERE signum = $1", signum).Scan(&user.Role, &user.ProgramName)
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusOK, map[string]interface{}{
				"signum":       signum,
				"role":         "normal",
				"access_level": "2",
			}," ")
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}

		if user.Role == "admin" {
			return c.JSONPretty(http.StatusOK, map[string]interface{}{
				"signum":       signum,
				"programname":user.ProgramName,
				"role":         "admin",
				"access_level": "0",
			}, " ")
		}

		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"signum":       signum,
			"programname":user.ProgramName,
			"role":         "registered",
			"access_level": "1",
		}, " ")
	}
}

type AccessLevel struct {
	Application string `json:"application"`
	Feature     string `json:"feature"`
	AccessLevel string `json:"accesslevel"`
}

func GetAllAccessLevels(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := db.Query("SELECT application, feature, accesslevel FROM access_level")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
		}
		defer rows.Close()

		var accessLevels []AccessLevel
		for rows.Next() {
			var accessLevel AccessLevel
			err := rows.Scan(&accessLevel.Application, &accessLevel.Feature, &accessLevel.AccessLevel)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}
			accessLevels = append(accessLevels, accessLevel)
		}

		return c.JSONPretty(http.StatusOK, accessLevels, " ")
	}
}

type HomeHandler struct {
	db *sql.DB
}

func NewHomeHandler(db *sql.DB) *HomeHandler {
	return &HomeHandler{
		db: db,
	}
}

// struct for homepage

type Homepage struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

// function for getting homepage link

func (h *HomeHandler) GetHomepageDocumentLink(c echo.Context) error {
	var documentation Homepage
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'Homepage Link'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Homepage link not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve Homepage Link")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"homepage_link": documentation.DocumentLink}, " ")
}

// function for updating the homepage link

func (h *HomeHandler) UpdateHomepageDocumentationLink(c echo.Context) error {
	homepage := new(Homepage)
	if err := c.Bind(homepage); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid homepage data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'Homepage Link'"
	_, err := h.db.Exec(updateStmt, homepage.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update Homepage Link")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "Homepage Link updated successfully"}, " ")
}