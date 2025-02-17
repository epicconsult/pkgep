package pkgep

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigName("th-admin-divisions")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs/masterdata/")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("cannot read data in th-admin-divisions.json: %s", err.Error()))
	}

	fmt.Println("viper successfully init")
}
