package snowflake

import (
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
	"myBulebell/pkg/conf"
)

var node *snowflake.Snowflake

func Init() {
	var err error
	node, err = snowflake.NewSnowflake(
		conf.AppConf.DataCenterId,
		conf.AppConf.MachineId,
	)
	if err != nil {
		panic(err)
	}
}

func GenID() int64 {
	return node.NextVal()
}
