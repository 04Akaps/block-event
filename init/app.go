package init

import (
	"github.com/04Akaps/block-event/init/config"
	"github.com/04Akaps/block-event/log"
	"github.com/04Akaps/block-event/repository"
)

type App struct {
	cfg *config.Config

	repository *repository.Repository
}

func StartApp(cfg *config.Config) {
	a := &App{cfg: cfg}

	var err error

	if a.repository, err = repository.NewRepository(cfg); err != nil {
		log.CritLog("Failed To Connect Repository")
	}

}
