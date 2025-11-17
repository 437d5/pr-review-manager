package routers

import (
	"net/http"

	"github.com/437d5/pr-review-manager/internal/application/http/handlers"
)

func InitRouter(
	teamHandler *handlers.TeamHandler,
	userHandler *handlers.UserHandler,
	prHandler *handlers.PRHandler,
) *http.ServeMux {
	mainRouter := http.NewServeMux()

	prRouter := http.NewServeMux()
	teamRouter := http.NewServeMux()
	userRouter := http.NewServeMux()

	prRouter.HandleFunc("POST /create", prHandler.CreatePR)
	prRouter.HandleFunc("POST /merge", prHandler.Merge)
	prRouter.HandleFunc("POST /reassign", prHandler.Reassign)

	userRouter.HandleFunc("POST /setIsActive", userHandler.SetIsActive)
	userRouter.HandleFunc("GET /getReview", userHandler.GetPRs)

	teamRouter.HandleFunc("POST /add", teamHandler.CreateTeam)
	teamRouter.HandleFunc("GET /get", teamHandler.GetTeam)

	mainRouter.Handle("/team/", http.StripPrefix("/team", teamRouter))
	mainRouter.Handle("/users/", http.StripPrefix("/users", userRouter))
	mainRouter.Handle("/pullRequest/", http.StripPrefix("/pullRequest", prRouter))

	return mainRouter
}
