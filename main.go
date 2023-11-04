package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kelsman/go-fiber-postgres/models"
	"github.com/kelsman/go-fiber-postgres/storage"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type (
	Repository struct {
		DB *gorm.DB
	}
	Book struct {
		Author    string `json:"book"`
		Title     string `json:"title"`
		Publisher string `json:"publisher"`
	}
)

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Request failed",
		})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not create the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book has been added"})
	return nil
}

func (r *Repository) GetBooks(ctx *fiber.Ctx) error {
	bookModels := &[]models.Book{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the books"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books fetched successfully", "data": bookModels})
	return nil

}

func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := models.Book{}
	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "id cannot be empty"})
	}

	err := r.DB.Delete(&bookModel, id).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "could  not find book"})
		return err
	}
	return nil
}
func (r *Repository) GetBookByID(ctx *fiber.Ctx) error {
	bookModel := models.Book{}
	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "id cannot be empty"})
	}

	err := r.DB.First(&bookModel, id).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "could  not find book"})
		return err
	}
	ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "could  not find book", "data": bookModel})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	//api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)

}
func main() {
	// Load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbConfig := storage.Config{
		Host:     os.Getenv("PGHOST"),
		Password: os.Getenv("PGPASSWORD"),
		User:     os.Getenv("PGUSER"),
		Port:     os.Getenv("PG_PORT"),
		DbName:   os.Getenv("PGDATABASE"),
		SSLMode:  os.Getenv("PG_SSL"),
	}
	fmt.Println("%v", dbConfig)
	db, err := storage.NewConnection(&dbConfig)
	if err != nil {
		log.Fatal("could not load database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate database")
	}
	r := Repository{
		DB: db,
	}

	app := fiber.New()
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

	r.SetupRoutes(app)

}
