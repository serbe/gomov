package main

func main() {
	cfg, err := getConfig()
	if err != nil {
		return
	}
	initDB(
		cfg.Base.Host,
		cfg.Base.Dbname,
		cfg.Base.User,
		cfg.Base.Password,
		cfg.Base.LogSQL,
	)
	initServer(cfg.Web.Host+":"+cfg.Web.Port, cfg.Web.Log)
}
