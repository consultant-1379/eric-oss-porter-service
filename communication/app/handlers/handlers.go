package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// NewsFeedItem represents a single newsfeed item
type NewsFeedItem struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Timestamp string `json:"timestamp"`
	Is_enable string `json:"is_enable"`
	}

	// GetNewsFeedHandler retrieves newsfeed data from the database and returns it as a JSON response
	func GetNewsFeedHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

	// newsfeedLimit := 20

	// Retrieve newsfeed data from the database
	newsfeedItems, err := GetNewsFeed(db)
	if err != nil {
	return c.JSONPretty(http.StatusInternalServerError, map[string]string{
	"error": "Internal Server Error",
	}, " ")
	}

	// Return the response
	//return c.JSONPretty(http.StatusOK, newsfeedItems, " ")
	response := map[string]interface{}{
		"format": "newsfeed",
		"fields": []map[string]string{
			{
				"type": "Id",
				"field_name": "id",
			  },
			  {
				"type": "Header",
				"field_name": "title",
			  },
			  {
				"type": "Details",
				"field_name": "content",
			  },
			  {
				"type": "Date",
				"field_name":"timestamp",
			  },
			  {
				"type": "Is_enable",
				"field_name":"is_enable",
			  },
		},
		"data":newsfeedItems,
	}
	return c.JSONPretty(http.StatusOK, response, " ")
	}
	}

	// GetNewsFeed retrieves newsfeed data from the database
	func GetNewsFeed(db *sql.DB) ([]NewsFeedItem, error) {
	query := "SELECT id,title,content, created_at,is_enable FROM communication WHERE type = 'newsfeed' AND is_enable = 'true' ORDER BY created_at DESC "
	rows, err := db.Query(query)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	var newsfeedItems []NewsFeedItem
	for rows.Next() {
	var id string
	var title string
	var content string
	var createdAt time.Time
	var is_enable string
	err := rows.Scan(&id,&title,&content, &createdAt,&is_enable)
	if err != nil {
	return nil, err
	}

	item := NewsFeedItem{
	Id: id,
	Title: title,
	Content: content,
	Timestamp: createdAt.Format("Monday 02 January, 2006 03:04 PM"),
	Is_enable: is_enable,
	}
	newsfeedItems = append(newsfeedItems, item)
	}

	if err := rows.Err(); err != nil {
	return nil, err
	}

	return newsfeedItems, nil
	}

	// AnnouncementItem represents a single announcement item
	type AnnouncementItem struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Is_enable string `json:"is_enable"`
	//Content string `json:"content"`
	//Timestamp string `json:"timestamp"`
	}

	// GetAnnouncementsHandler retrieves announcements data from the database and returns it as a JSON response
	func GetAnnouncementsHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
	//announcementLimit := 2
	// Retrieve announcements data from the database
	announcementItems, err := GetAnnouncements(db)
	if err != nil {
	return c.JSONPretty(http.StatusInternalServerError, map[string]string{
	"error": "Internal Server Error",
	}, " ")
	}

	// Return the response
	//return c.JSONPretty(http.StatusOK, announcementItems, " ")
	response := map[string]interface{}{
		"format": "announcements",
		"fields": []map[string]string{
			{
				"type": "Id",
				"field_name": "id",
			  },
			  {
				"type": "Header",
				"field_name": "title",
			  },
			  {
				"type": "Is_enable",
				"field_name": "is_enable",
			  },
		},
		"data":announcementItems,
	}
	return c.JSONPretty(http.StatusOK, response, " ")
	}

	}

	// GetAnnouncements retrieves announcements data from the database
	func GetAnnouncements(db *sql.DB) ([]AnnouncementItem, error) {
	query := "SELECT id,title,is_enable FROM communication WHERE type = 'announcements' AND is_enable='true' ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	var announcementItems []AnnouncementItem
	for rows.Next() {
	var id string
	var title string
	var is_enable string
	//var content string
	//var createdAt time.Time
	err := rows.Scan(&id,&title,&is_enable)
	if err != nil {
	return nil, err
	}

	item := AnnouncementItem{
	Id: id,
	Title: title,
	Is_enable: is_enable,
	//Content: content,
	//Timestamp: createdAt.Format("Monday 02 January, 2006 03:04 PM"),
	}
	announcementItems = append(announcementItems, item)
	}

	if err := rows.Err(); err != nil {
	return nil, err
	}

	return announcementItems, nil
	}

	// CommunicationColumn represents a single column in the communication response
	type CommunicationColumn struct {
	Type string `json:"type"`
	FieldName string `json:"field_name"`
	}

	// CommunicationData represents the data for newsfeed and announcement
	type CommunicationData struct {
	Format string `json:"format"`
	Fields []CommunicationColumn `json:"fields"`
	Data []CommunicationItem `json:"data"`
	}

	type CommunicationAnnouncementData struct{
		Format string `json:"format"`
		Fields []CommunicationColumn `json:"fields"`
		Data []AnnouncementItem `json:"data"`
	}

	// CommunicationResponse represents the response structure for the /communication endpoint
	type CommunicationResponse struct {
	Newsfeed CommunicationData `json:"newsfeed"`
	Announcement CommunicationAnnouncementData `json:"announcement"`
	}

	// CommunicationItem represents a single communication item
	type CommunicationItem struct {
		Id string `json:"id"`
		Title string `json:"title"`
		Content string `json:"content"`
		Timestamp string `json:"timestamp"`
		Is_enable string `json:"is_enable"`
	}


// GetCommunicationHandler retrieves communication data (newsfeed and announcements) from the database and returns it as a JSON response
func GetCommunicationHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
	// Retrieve newsfeed and announcements data from the database
	newsfeedItems, err := GetNewsFeed(db)
	if err != nil {
	return c.JSONPretty(http.StatusInternalServerError, map[string]string{
	"error": "Internal Server Error",
	}, " ")
	}

	announcementItems, err := GetAnnouncements(db)
	if err != nil {
	return c.JSONPretty(http.StatusInternalServerError, map[string]string{
	"error": "Internal Server Error",
	}, " ")
	}

	// Create the communication response
	response := CommunicationResponse{
	Newsfeed: CommunicationData{
	Format: "newsfeed",
	Fields: []CommunicationColumn{
	{Type: "Id",FieldName: "id"},
	{Type: "Title", FieldName: "title"},
	{Type: "Details", FieldName: "content"},
	{Type: "Date", FieldName: "timestamp"},
	{Type: "Is_enable", FieldName: "is_enable"},
	},
	Data: make([]CommunicationItem, len(newsfeedItems)),
	},
	Announcement: CommunicationAnnouncementData{
	Format: "announcement",
	Fields: []CommunicationColumn{
	{Type: "Id",FieldName: "id"},
	{Type: "Title", FieldName: "title"},
	{Type: "Is_enable", FieldName: "is_enable"},
	//{Type: "Details", FieldName: "content"},
	//{Type: "Date", FieldName: "timestamp"},
	},
	Data: make([]AnnouncementItem, len(announcementItems)),
	},
	}

	// Populate the newsfeed items in the response
	for i, item := range newsfeedItems {
	response.Newsfeed.Data[i] = CommunicationItem{
	Id: item.Id,
	Title: item.Title,
	Content: item.Content,
	Timestamp: item.Timestamp,
	Is_enable: item.Is_enable,
	}
	}

	// Populate the announcement items in the response
	for i, item := range announcementItems {
	response.Announcement.Data[i] = AnnouncementItem{
	Id: item.Id,
	Title: item.Title,
	Is_enable: item.Is_enable,
	}
	}

	// Return the response
	return c.JSONPretty(http.StatusOK, response, " ")
	}
	}

type CommunicationRequest struct {
	Title string `json:"title"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func CreateCommunication(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body
		req := new(CommunicationRequest)
		if err := c.Bind(req); err != nil {
			log.Println(err)
			return c.JSONPretty(http.StatusBadRequest, "Invalid request payload"," ")
		}

		// Get the current timestamp
		now := time.Now()

		// Insert the communication entry into the database
		insertStmt := "INSERT INTO communication (title, type, content, created_at,is_enable) VALUES ($1, $2, $3, $4,true)"
		_, err := db.Exec(insertStmt,req.Title, req.Type, req.Content, now)
		if err != nil {
			log.Println(err)
			return c.JSONPretty(http.StatusInternalServerError, "Failed to create communication entry"," ")
		}

		// Return the created communication entry
		response := struct {
			Title string `json:"title"`
			Type      string    `json:"type"`
			Content   string    `json:"content"`
			CreatedAt time.Time `json:"created_at"`
		}{
			Title: req.Title,
			Type:      req.Type,
			Content:   req.Content,
			CreatedAt: now,
		}

		return c.JSONPretty(http.StatusCreated, response," ")
	}
}

type Mail struct {
	MailCommunication string    `json:"mailcommunication"`
	CreatedAt         time.Time `json:"created_at"`
}

func GetMailCommunication(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := "SELECT content, created_at FROM mail ORDER BY id DESC LIMIT 1"
		row := db.QueryRow(query)

		var content string
		var createdAt time.Time
		err := row.Scan(&content, &createdAt)
		if err != nil {
			log.Println(err)
			return c.JSONPretty(http.StatusInternalServerError, "Failed to retrieve mail communication"," ")
		}

		response := struct {
			MailCommunication string `json:"mailcommunication"`
		}{
			MailCommunication: content,
		}

		return c.JSONPretty(http.StatusOK, response, " ")
	}
}

func CreateMailCommunication(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(CommunicationRequest)
		if err := c.Bind(req); err != nil {
			log.Println(err)
			return c.JSONPretty(http.StatusBadRequest, "Invalid request payload"," ")
		}

		now := time.Now()

		insertStmt := "INSERT INTO mail (content, created_at) VALUES ($1, $2)"
		_, err := db.Exec(insertStmt, req.Content, now)
		if err != nil {
			log.Println(err)
			return c.JSONPretty(http.StatusInternalServerError, "Failed to create mail communication", " ")
		}

		response := struct {
			MailCommunication string    `json:"mailcommunication"`
			CreatedAt         time.Time `json:"created_at"`
		}{
			MailCommunication: req.Content,
			CreatedAt:         now,
		}

		return c.JSONPretty(http.StatusCreated, response, " ")
	}
}

// UpdateNewsFeedItemHandler updates a newsfeed item by its ID sent in the request body
func UpdateNewsFeedItemHandler(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Parse the request body
        req := new(struct {
            ID      string `json:"id"`
            Title   string `json:"title"`
            Content string `json:"content"`
			Is_enable string `json:"is_enable"`
        })

        if err := c.Bind(req); err != nil {
            log.Println(err)
            return c.JSONPretty(http.StatusBadRequest, "Invalid request payload", " ")
        }

        // Update the newsfeed item in the database using the ID from the request body
        updateStmt := "UPDATE communication SET title = $1, content = $2, is_enable = $3 WHERE id = $4"
        _, err := db.Exec(updateStmt, req.Title, req.Content, req.Is_enable, req.ID)
        if err != nil {
            log.Println(err)
            return c.JSONPretty(http.StatusInternalServerError, "Failed to update newsfeed item", " ")
        }

        // Return the updated newsfeed item
        response := struct {
            Message string `json:"message"`
        }{
            Message: "Newsfeed item updated successfully",
        }

        return c.JSONPretty(http.StatusOK, response, " ")
    }
}

// UpdateNewsFeedItemHandler updates a newsfeed item by its ID sent in the request body
func UpdateAnnouncementItemHandler(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Parse the request body
        req := new(struct {
            ID      string `json:"id"`
            Title   string `json:"title"`
			Is_enable string `json:"is_enable"`
        })

        if err := c.Bind(req); err != nil {
            log.Println(err)
            return c.JSONPretty(http.StatusBadRequest, "Invalid request payload", " ")
        }

        // Update the newsfeed item in the database using the ID from the request body
        updateStmt := "UPDATE communication SET title = $1, is_enable = $2 WHERE id = $3"
        _, err := db.Exec(updateStmt, req.Title, req.Is_enable, req.ID)
        if err != nil {
            log.Println(err)
            return c.JSONPretty(http.StatusInternalServerError, "Failed to update newsfeed item", " ")
        }

        // Return the updated newsfeed item
        response := struct {
            Message string `json:"message"`
        }{
            Message: "Announcement item updated successfully",
        }

        return c.JSONPretty(http.StatusOK, response, " ")
    }
}