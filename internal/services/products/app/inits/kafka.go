package inits

import (
	"context"
	"pkg/kafka"
	"pkg/logger"
)

func InitKafka(ctx context.Context, kafka *kafka.KafkaClient, cnf *kafka.KafkaClientsConfig, log logger.Zapper) {

	productConf, ok := (*cnf)["products-producer"]

	if !ok {
		log.Infof(ctx, "Config not found %v", productConf)
		return
	}

	kafka.GetProducer(ctx, productConf.Topic, productConf)
}
