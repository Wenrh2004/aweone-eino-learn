package third

import (
	"github.com/spf13/viper"

	"github.com/Wenrh2004/travel_assistant/pkg/third/amap"
)

func NewAmapClient(conf *viper.Viper) *amap.Client {
	return amap.NewClient(conf.GetString("app.amap.key"), conf.GetDuration("app.amap.timeout"))
}
