package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type linkService struct {
	db *sql.DB
}

func NewLinkService(db *sql.DB) *linkService {
	return &linkService{db}
}

type link struct {
	FULL_LINK  string `json:"fullLink"`
	SHORT_LINK string `json:"shortLink"`
}

func (s *linkService) CreateShortLink(c *fiber.Ctx) link {

	var payload link

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("Error: %v", err)
	}

	var l link

	sqlQuery := `
	SELECT short_link FROM links
	WHERE full_link = $1;
	`
	if err := s.db.QueryRow(sqlQuery, payload.FULL_LINK).Scan(&l.SHORT_LINK); err != nil {
		fmt.Println(err)
	}

	if l.SHORT_LINK != "" {
		l.FULL_LINK = payload.FULL_LINK
		return link{
			FULL_LINK:  payload.FULL_LINK,
			SHORT_LINK: os.Getenv("DOMAIN") + "/" + l.SHORT_LINK,
		}
	}

	shortHash := createHashAndCheckDB(s.db, l.FULL_LINK)

	sqlQuery = `
	INSERT INTO links (full_link, short_link)
	VALUES ($1, $2)
	RETURNING *;
	`
	s.db.Exec(sqlQuery, payload.FULL_LINK, shortHash)

	shortLink := os.Getenv("DOMAIN") + "/" + shortHash

	return link{
		FULL_LINK:  payload.FULL_LINK,
		SHORT_LINK: shortLink,
	}
}

func (s *linkService) GetFullLink(hash string) string {

	var l link

	sqlQuery := `
	SELECT full_link FROM links
	WHERE short_link = $1;
	`
	err := s.db.QueryRow(sqlQuery, hash).Scan(&l.FULL_LINK)
	if err != nil {
		fmt.Println(err)
	}

	return l.FULL_LINK
}

func createHashAndCheckDB(db *sql.DB, fullLink string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(fullLink), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	linkStr := string(hash[:])
	shortHash := ""
	for i := 7; i < 13; i++ {
		var char = string(linkStr[i])
		if char != "." && char != "/" {
			shortHash += char
		}
	}

	var l link
	sqlQuery := `
	SELECT short_link from links
	WHERE short_link = $1;
	`
	err = db.QueryRow(sqlQuery, shortHash).Scan(&l.SHORT_LINK)
	if err != nil {
		fmt.Println(err)
	}

	if l.SHORT_LINK != "" {
		fmt.Println("Hash already exists in db, creating new short hash")
		return createHashAndCheckDB(db, fullLink)
	}

	return shortHash
}
