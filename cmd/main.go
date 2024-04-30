package main

import (
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":8080"))
}

func newServer() *handler.Server {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbDsn := os.Getenv("DATABASE_URL")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	var usecase usecase.UsecaseInterface = usecase.NewUsecase(usecase.NewUsecaseOptions{
		Repository: repo,
	})

	opts := handler.NewServerOptions{
		Usecase: usecase,
	}
	return handler.NewServer(opts)
}
