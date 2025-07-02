package api

import (
	"RetoTecnicoDragonBall/business-management/application/handler"
	"RetoTecnicoDragonBall/business-management/application/usecase"
	"RetoTecnicoDragonBall/business-management/infrastructure/repository"
	"RetoTecnicoDragonBall/internal/db"
	internal "RetoTecnicoDragonBall/internal/http"
	"RetoTecnicoDragonBall/internal/logs"
	"RetoTecnicoDragonBall/internal/utils"
	"RetoTecnicoDragonBall/internal/utils/helpers"
	"net/http"
	"time"
)

func Start() {

	// Instances Others
	log := logs.NewLog()

	// ClientRest
	clientRest := internal.NewClientRest(log)

	// Postgres
	conPostgresDB := db.NewPostgresConnection(log)
	clientPostgresDB := db.NewClientPostgres(conPostgresDB, log)

	// Repository
	repositoryPoc := repository.NewDragonBallRepository(clientPostgresDB, log)

	// Use Case Instance
	useCasePoc := usecase.NewDragonBallUseCase(repositoryPoc, clientRest, log)

	r := routes(
		&HandlerDto{
			dragonBallHandler: handler.NewDragonBallHandler(useCasePoc, log),
		},
	)
	timeout, _ := time.ParseDuration("300s")

	server := newServer(log)
	s := &http.Server{
		Addr:    ":8080",
		Handler: utils.TimeoutHandlerServer(r, timeout, helpers.GetJsonTimeOut()),
	}
	s.ListenAndServe()

	server.Start(":8080", r)
}
