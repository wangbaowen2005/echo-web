package api

import (
	// "time"

	// ec "github.com/hobo-go/echo-mw/cache"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/hobo-go/echo-mw/binder"

	"github.com/hobo-go/echo-web/conf"
	"github.com/hobo-go/echo-web/modules/cache"
	"github.com/hobo-go/echo-web/modules/session"
)

//-----
// API Routers
//-----
func Routers() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Customization
	// e.SetLogPrefix("Echo")
	e.SetLogLevel(log.DEBUG)
	if conf.RELEASE_MODE {
		e.SetDebug(false)
	}

	// CORS
	e.Use(mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins: []string{"http://echo.www.localhost:8080", "http://echo.api.localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// CSRF
	e.Use(mw.CSRFWithConfig(mw.CSRFConfig{
		TokenLookup: "form:X-XSRF-TOKEN",
	}))

	// Gzip
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Bind
	e.SetBinder(binder.New())

	// Session
	e.Use(session.Session())

	// Cache
	e.Use(cache.Cache())

	// e.Use(ec.SiteCache(ec.NewMemcachedStore([]string{conf.MEMCACHED_SERVER}, time.Hour), time.Minute))
	// e.Get("/user/:id", ec.CachePage(ec.NewMemcachedStore([]string{conf.MEMCACHED_SERVER}, time.Hour), time.Minute, UserHandler))

	// Auth
	// e.Use(auth.Auth(models.GenerateAnonymousUser))

	// Routers
	e.Get("/", ApiHandler)
	e.Get("/:id", ApiHandler)
	e.Get("/login", UserLoginHandler)
	e.Get("/register", UserRegisterHandler)

	// JWT
	r := e.Group("/restricted")
	r.Use(mw.JWTWithConfig(mw.JWTConfig{
		SigningKey:  []byte("secret"),
		ContextKey:  "_user",
		TokenLookup: "header:" + echo.HeaderAuthorization,
	}))

	// curl http://echo.api.localhost:8080/restricted/user -H "Authorization: Bearer XXX"
	r.Get("/user", UserHandler)

	post := r.Group("/post")
	{
		post.Get("/save", PostSaveHandler)
		post.Get("/id/:id", PostHandler)
		post.Get("/:userId/p/:p/s/:s", PostsHandler)
	}

	return e
}
