package main

import (
	"database/sql" // New import
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.cesar90.com/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // New import
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as developer progresses
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	//Define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "root:test@/snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use the flag.Parse() function to parse the command-file flag
	// This reads in the command-line flag value and assigns it to addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any erros are
	// encountered during parsing the application will be terminated
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stram and uses the default settings
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exist
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initalize a decoder instance
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger)
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Use the http.NewServerMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	// mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that path given to the http.Dir function is relative to the project
	// directy root
	// fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.handle() function to register the file server as the handler for
	// all URL paths that start with "/static". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("GET /{$}", app.home)                      //Restrict this route to exact matches on / only
	// mux.HandleFunc("GET /snippet/view/{id}", app.snippetView) // Add the (id) wildcard segment
	// mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	// Create the new route, which is restricted to POST request only,
	// mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Print a log message to say that the server is starting
	// log.Print("Starting server on :4000")
	// Use the info() method to log the starting server message at info severity
	// (along with the listen address as an attribute)
	logger.Info("starting server", "addr", *addr)

	// log.Printf("starting server on %s", *addr)

	// Use the http.listenAndServe() function to start a new web server. Web pass in
	// two parameters: the TCP network address to listen on (in this case "":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil
	// err := http.ListenAndServe(":4000", mux)
	// err = http.ListenAndServe(*addr, app.routes())

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at the Error, level, and assign it to the ErrorLog Field. If
		// you would prefer to log the server errors at Warn level instead, you,
		// could pass slog.LevelWarn as the final parameter
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()

	// log.Fatal(err)
	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServe() at Error severity (with no additional attributes)
	// and the call os.Exit(1) to termiante the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}

// The OpenDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
