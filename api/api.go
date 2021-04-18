package api

import (
	"net/http"

	"github.com/fberrez/simple-uptime-backend/backend"
	"github.com/fberrez/simple-uptime-backend/repositories"
	"github.com/fberrez/simple-uptime-backend/services"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

// API represents an API instance.
type API struct {
	// engine is the API engine
	engine *gin.Engine
	// fizz is the gin wrapper used to run the API
	fizz *fizz.Fizz

	// accountService is an instance of the service Account
	accountService *services.Account
	// loginService is an instance of the service Login
	loginService *services.Login
}

// New returns a new instance of API.
func New(backend backend.Backend) *API {
	engine := gin.New()

	engine.Use(logger.SetLogger(), gin.Recovery())

	f := fizz.NewFromEngine(engine)

	// initializes openapi route
	infos := &openapi.Info{
		Title:       "Minimalist.Uptime",
		Description: `A minimalist API used to get services uptime.`,
		Version:     "1.0.0",
	}
	f.GET("/openapi.json", nil, f.OpenAPI(infos, "json"))

	// initializes repositories
	accountRepository := repositories.NewAccountRepository(backend)

	// initializes api
	api := &API{
		engine:         engine,
		fizz:           f,
		accountService: services.NewAccountService(accountRepository),
		loginService:   services.NewLoginService(accountRepository),
	}

	// initializes group /account
	accountGroup := f.Group("/account", "Account", "Everything to know about account management")
	accountGroup.POST("/", []fizz.OperationOption{
		fizz.Summary("Create a new account"),
	}, tonic.Handler(api.createAccount, 200))
	accountGroup.PUT("/:id", []fizz.OperationOption{
		fizz.Summary("Update an existing account"),
	}, tonic.Handler(api.updateAccount, 200))
	accountGroup.DELETE("/:id", []fizz.OperationOption{
		fizz.Summary("Delete an existing account"),
	}, tonic.Handler(api.deleteAccount, 200))

	// initializes login route
	f.POST("/login", []fizz.OperationOption{
		fizz.Summary("Perform an authentication"),
	}, tonic.Handler(api.handleLogin, 200))

	return api
}

// ServeHTTP calls the gin built-in function ServeHTTP().
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.fizz.ServeHTTP(w, r)
}
