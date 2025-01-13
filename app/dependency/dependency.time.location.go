package dependency

import (
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
)

func ProvideTimeZoneLocation(cfg config.Cfg) (*time.Location, error) {
	return time.LoadLocation(cfg.App.Timezone)
}
