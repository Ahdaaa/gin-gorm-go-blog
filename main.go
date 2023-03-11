package main

/* baca filza :
- udah berhasil buat user + handle password enkripsi dan email duplikat
- udah berhasil buat method login (mencocokkan apakah email ada/tidak, password benar/tidak)
- udah berhasil generate token ketika selesai login.
- udah berhasil validate token ketika mengakses url secured
- udah berhasil ambil id di token, terus berhasil buat blog sesuai dengan id yang di token!!!
- kurang komen, like + operasi delete akun + lihat detail blog + ganti detail user
*/

import (
	"log"
	"oprec/go-blog/config"
	"oprec/go-blog/controller"
	"oprec/go-blog/middleware"
	"oprec/go-blog/repository"
	"oprec/go-blog/routes"
	"oprec/go-blog/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	db := config.SetupDatabaseConnection()

	userRepository := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepository)

	userController := controller.NewUserController(userService)

	// mereka saling dependen
	defer config.CloseDatabaseConnection(db)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRoutes(server, userController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}

/*
format register:
{
	"name":"bebek123",
	"email":"bebekgoreng@gmail.com",
	"password":"Anaana123"
}
format login:
{
    "Email":"ahda123@gmail.com",
    "Password":"Anaana123"
}

format upload blog:
{
    "Judul":"blog`",
    "tgl_post":"23 Mei 2003",
	"Isi":"saya adalah goblog"
}

format komen:
{
	"komen":"blablabla"
}

format ganti nama:
{
	"name":"smthing"
}

*/
