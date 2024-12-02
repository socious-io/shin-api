package shortner

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type ShortnerURL struct {
	ID        uuid.UUID `db:"id" json:"id"`
	LongURL   string    `db:"long_url" json:"long_url"`
	ShortID   string    `db:"short_id" json:"short_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ShortnerURLNewForm struct {
	LongURL string `json:"long_url" validate:"required"`
}

func (ShortnerURL) TableName() string {
	return "urls_shortens"
}

func (ShortnerURL) FetchQuery() string {
	return "shortens/fetch"
}

func (s *ShortnerURL) Create() error {
	rows, err := database.Queryx("shortens/create", s.LongURL)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(s); err != nil {
			return err
		}
	}
	return nil
}

func New(url string) (*ShortnerURL, error) {
	s := new(ShortnerURL)
	s.LongURL = url
	if err := s.Create(); err != nil {
		return nil, err
	}
	return s, nil
}

func Fetch(shortID string) (*ShortnerURL, error) {
	s := new(ShortnerURL)
	if err := database.Fetch(s, shortID); err != nil {
		return nil, err
	}
	return s, nil
}

func Routers(router *gin.RouterGroup) {
	router.GET("/:short_id/fetch", func(c *gin.Context) {
		shortID := c.Param("short_id")
		s, err := Fetch(shortID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	router.GET("/:short_id", func(c *gin.Context) {
		shortID := c.Param("short_id")
		s, err := Fetch(shortID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusPermanentRedirect, s.LongURL)
	})

	router.POST("", func(c *gin.Context) {
		form := new(ShortnerURLNewForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s := new(ShortnerURL)
		s.LongURL = form.LongURL
		if err := s.Create(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

}
