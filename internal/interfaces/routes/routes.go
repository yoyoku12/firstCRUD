package routes

import (
	"URL_SHORT/internal/interfaces/controllers"
	"URL_SHORT/internal/interfaces/middleware"
	"net/http"
)

func RegisterRoutes(urlController *controllers.URLController, authController *controllers.AuthController, authMiddleware *middleware.AuthMiddleware) {
	http.Handle("/longToShort", authMiddleware.Middleware(http.HandlerFunc(urlController.LongToShort)))
	http.Handle("/", http.HandlerFunc(urlController.LongToShort))
	http.HandleFunc("/register", authController.Register)
	http.HandleFunc("/login", authController.Login)
}
