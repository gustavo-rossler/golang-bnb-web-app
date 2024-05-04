package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gustavo-rossler/golang-bnb-web-app/pkg/config"
	"github.com/gustavo-rossler/golang-bnb-web-app/pkg/handlers"
	"github.com/gustavo-rossler/golang-bnb-web-app/pkg/render"
)

const portNumber = ":8089"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("could not create the template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	log.Printf("Staring application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
