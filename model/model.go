package model

import (
	"lotteryapi/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LotteryResult struct {
	ID            primitive.ObjectID           `json:"_id,omitempty" bson:"_id,omitempty"`
	LotteryName   string                       `json:"lottery_name,omitempty"`
	LotteryDate   string                       `json:"lottery_date,omitempty"`
	LotteryTime   string                       `json:"lottery_time,omitempty"`
	LotteryLink   string                       `json:"lottery_link,omitempty"`
	LotteryPrizes map[string]domain.PrizeCodes `json:"lottery_results,omitempty"`
}
