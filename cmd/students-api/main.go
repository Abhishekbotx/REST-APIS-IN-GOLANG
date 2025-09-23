package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhishekbotx/golang-restapi/internal/config"
	"github.com/abhishekbotx/golang-restapi/internal/http/handlers/student"
	"github.com/abhishekbotx/golang-restapi/internal/storage/sqlite"
)

func main() {
	// 🔹 1. Load configuration (from YAML/env) so we know which port/address to bind to
	cfg := config.Mustload()

	// database setup
	storage, err := sqlite.New(cfg) //✨here storage is of type sqlite.Sqlite struct
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialised", slog.String("env", cfg.Env), slog.String("version", "1.25.1"))

	// 🔹 2. Setup the HTTP router (ServeMux is Go's default lightweight router)
	router := http.NewServeMux()

	// Register a GET handler at /hi
	// Note: "GET /hi" means this route only matches GET requests to /hi
	router.HandleFunc("POST /api/students", student.New(storage)) //✨ pass any struct which implements storage interface , as new method of student.go requires storage interface
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage)) //✨ pass any struct which implements storage interface , as new method of student.go requires storage interface
	router.HandleFunc("GET /api/allstudents", student.GetStudents(storage)) //✨ pass any struct which implements storage interface , as new method of student.go requires storage interface

	// 🔹 3. Setup the HTTP server with address from config
	server := http.Server{
		Addr:    cfg.Addr, // comes from config.yaml
		Handler: router,   // attach our router
	}

	// 🔹 4. Create a channel to listen for OS signals (like Ctrl+C or `kill` command)
	done := make(chan os.Signal, 1)

	// Notify the channel when we receive termination signals (graceful shutdown)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 🔹 5. Start the server in a goroutine
	// Why? Because ListenAndServe is blocking, and we still want to listen for shutdown signals.
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// Fatal here will crash immediately — good because failing to start server is critical
			log.Fatal("failed to start server")
		}
	}()

	// At this point, the server is running in background

	// Log info: Server started successfully
	slog.Info("server started ", slog.String("address", cfg.Addr))
	fmt.Println("server started successfully")

	// 🔹 6. Block until we receive something on `done` channel
	// This done line.no 60 keeps main() alive until a shutdown signal is received.
	<-done
	//if its recieved contine the main goroutine

	// 🔹 7. Start graceful shutdown
	slog.Info("shutting down the server")

	// Create a context with timeout (5 seconds)
	// Why? We want to give active requests time to finish before forcefully killing the server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //context.background is empty starting point
	defer cancel()                                                          // ensures cleanup when main exits

	// Gracefully shutdown server (🍁stop accepting  new requests, wait for ongoing requests to finish)
	// Gracefully shutdown the server
	// --------------------------------------------------------
	// What server.Shutdown(ctx) does under the hood:
	// 1. Stops listening on the TCP socket → no new clients can connect.
	// 2. Closes idle connections immediately.
	// 3. Keeps active connections open (the ones already being served).
	// 4. Waits until either:
	//      - All active requests finish, OR
	//      - The provided context (ctx) times out.
	//
	// That’s why we passed:
	//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//
	// → If requests don’t finish in 5s, the server forces shutdown.

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error:", err.Error()))
	}

	// Final log
	slog.Info("server shutdown successfully")
}



/*
How Dependency Injection and Interfaces Work in main.go

1. You create a concrete implementation (`Sqlite`) of the `Storage` interface:
   storage, err := sqlite.New(cfg) ln:24
   // Here, `storage` is of type `*sqlite.Sqlite` which implements the `Storage` interface.

2. You pass this implementation to the handler when setting up your routes:
   router.HandleFunc("POST /api/students", student.New(storage)) ln:36
   // The handler only requires something that implements the `Storage` interface.

3. Above two points is called "wiring":
   - You connect the concrete implementation (`Sqlite`) to the interface (`Storage`) that the handler expects.
   - The handler does not care about the actual type, only that it fulfills the contract (interface).

4. Benefits:
   - The handler only knows about the interface, so it’s reusable and testable.
   - The main function chooses which implementation to use and connects it to the handler.
   - This pattern is called dependency injection and is a key part of writing flexible, maintainable code.

Summary:
- Interface defines the contract.
- Concrete struct implements the contract.
- Handler depends on the interface, not the concrete type.
- Main wires everything together by passing the implementation to the handler.
*/


/*

explain depenedency injection plz

Dependency injection is a design pattern where you provide (inject) the dependencies a component needs from the outside, rather than having it create them itself.

In your code:
Your handler (student.New) needs a way to store students, but it doesn’t care how—it just needs something that implements the Storage interface.
In main.go, you create a concrete implementation (Sqlite) and pass it to the handler.
This means the handler is dependent on the interface, not the concrete type.
Why is this useful?
Flexibility: You can swap out Sqlite for another storage backend (like Postgres) without changing the handler.
Testability: You can inject a mock implementation for testing.
Loose coupling: Your code is easier to maintain and extend.
Summary:
Dependency injection lets you build reusable, testable, and flexible code by passing dependencies (like database access) into components, rather than hard-coding them.
*/