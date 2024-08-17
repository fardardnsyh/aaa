package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Burakbgmk/go-tbc-bot/internal/ai"
	"github.com/Burakbgmk/go-tbc-bot/internal/util"
	"github.com/tmc/langchaingo/vectorstores/chroma"

	"github.com/Burakbgmk/go-tbc-bot/public/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/styles", "styles")
	e.Static("/uploads", "uploads")
	e.GET("/", s.IndexHandler)
	e.POST("/upload", s.UploadHandler)
	e.POST("/querry", s.QueryHandler)

	return e
}

func (s *Server) IndexHandler(c echo.Context) error {
	component := view.Index()

	return util.Render(c, component)
}

var store chroma.Store

func (s *Server) QueryHandler(c echo.Context) error {
	fmt.Printf("store: %v\n", store)
	body, error := io.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if error != nil {
		panic(error)
	}
	querry := string(body)

	fmt.Println(querry)
	answer, err := ai.QueryToVectorDB(querry, c, filename)
	if err != nil {
		return err
	}
	val := answer["text"].(string)

	return util.Render(c, view.SuccesfulPrompt(val))
}

var filename string

func (s *Server) UploadHandler(c echo.Context) error {
	fmt.Println("Uploading handler enters")

	file, err := c.FormFile("file-input")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create("./uploads/file.pdf")
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	filename = file.Filename

	store, err := ai.InsertToVectorDb(c, filename)
	if err != nil {
		return err
	}
	fmt.Printf("store: %v\n", store)

	return util.Render(c, view.SuccesfulUpload(filename, "/uploads/file.pdf"))

}
