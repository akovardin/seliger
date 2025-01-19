package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/template"
	"go.uber.org/fx"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"gohome.4gophers.ru/kovardin/seliger/app/config"
	"gohome.4gophers.ru/kovardin/seliger/app/handlers"
	"gohome.4gophers.ru/kovardin/seliger/app/settings"
	_ "gohome.4gophers.ru/kovardin/seliger/migrations"
	"gohome.4gophers.ru/kovardin/seliger/static"
)

func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

var (
	cfg string
)

func main() {
	cfg = os.Getenv("CONFIG")

	fx.New(
		handlers.Module,
		// tasks.Module,

		fx.Provide(func(conf config.Database) *pocketbase.PocketBase {

			app := pocketbase.NewWithConfig(pocketbase.Config{
				DBConnect: func(dbPath string) (*dbx.DB, error) {
					if strings.Contains(dbPath, "data.db") {
						return dbx.Open("libsql", conf.Url+"?authToken="+conf.Token)
					}

					return core.DefaultDBConnect(dbPath)
				},
			})

			return app
		}),
		fx.Provide(template.NewRegistry),
		fx.Provide(settings.New),
		fx.Provide(
			func() (config.Config, error) {
				return config.New(cfg)
			},
		),
		fx.Invoke(
			migration,
		),
		fx.Invoke(
			routing,
		),
	).Run()
}

func routing(
	app *pocketbase.PocketBase,
	lc fx.Lifecycle,
	settings *settings.Settings,
	home *handlers.Home,
	ads *handlers.Ads,
) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/", home.Home)

		// listing
		e.Router.GET("/ads/:publisher", ads.List)

		// static
		e.Router.GET("/static/{path...}", apis.Static(static.FS, false))

		return e.Next()

	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Start(); err != nil {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

// go run cmd/depot/main.go migrate collections --dir ./data/
func migration(app *pocketbase.PocketBase) {
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})
}
