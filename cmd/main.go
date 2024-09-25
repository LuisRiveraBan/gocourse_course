package main

import (
	"context"
	"fmt"
	course "github.com/LuisRiveraBan/gocourse_course/internal"
	"github.com/LuisRiveraBan/gocourse_course/pkg/bootstrap"
	"github.com/LuisRiveraBan/gocourse_course/pkg/handler"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func main() {

	// Load environment variables.
	_ = godotenv.Load()

	// Initialize the logger.
	l := bootstrap.InitLogger()

	// Connect to the database.
	db, err := bootstrap.ConnectToDatabase()
	if err != nil {
		l.Fatal("Error connecting to the database", err)
	}

	// Define the course endpoints and handlers.
	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")

	if pagLimDef == "" {
		l.Fatal("REQUESTER VARIABLE")
	}

	// Define the user endpoints and handlers.
	ctx := context.Background()
	courseRepo := course.NewRepository(l, db)
	courseSrv := course.NewService(l, courseRepo)
	h := handler.NewCourseHTTPServer(ctx, course.MakeEndpoints(courseSrv, course.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")

	address := fmt.Sprintf("127.0.0.1:%s", port)

	// Configuración del server HTTP.
	srv := &http.Server{
		Handler: accesControl(h),
		Addr:    address,
		//Tiempo de escritura
		WriteTimeout: 15 * time.Second,
		//Tiempo de lectura
		ReadTimeout: 15 * time.Second,
	}

	// Inicia el server.
	errCh := make(chan error)
	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	// Gracefully shut down the server.
	err = <-errCh
	if err != nil {
		l.Fatal("ListenAndServe: ", err)
	}

}

// Acceso controlado para los recursos RESTful.
func accesControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
