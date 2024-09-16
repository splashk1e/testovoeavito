package bootstrap

type Application struct {
	Env      *Env
	Postgres *PostgresClient
}

func App() Application {
	app := Application{}
	app.Env = NewEnv()
	app.Postgres = NewPostgresDb(app.Env)
	return app
}
