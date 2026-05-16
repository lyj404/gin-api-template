package idgen

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/lyj404/gin-api-template/config"
)

var node *snowflake.Node

func InitSnowflake() error {
	startTime, err := time.Parse("2006-01-02", config.CfgSnowflake.StartTime)
	if err != nil {
		return fmt.Errorf("snowflake start time parse error: %w", err)
	}
	snowflake.Epoch = startTime.UnixMilli()
	n, err := snowflake.NewNode(int64(config.CfgSnowflake.NodeID))
	if err != nil {
		return fmt.Errorf("snowflake node init error: %w", err)
	}
	node = n
	return nil
}

func NextID() uint64 {
	return uint64(node.Generate().Int64())
}
