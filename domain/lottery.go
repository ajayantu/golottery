package domain

import "time"

type GetLotteryResultRequest struct {
	SeriesName string
}
type PrizeCodes struct {
	PrizeMoney string   `json:"price_money"`
	PrizeCodes []string `json:"prize_codes"`
}
type GetLotteryResultRespose struct {
	LotteryName    string                `json:"lottery_name"`
	LotteryLink    string                `json:"lottery_link"`
	LotteryDate    time.Time             `json:"lottery_date"`
	LotteryTime    string                `json:"lottery_time"`
	LotteryResults map[string]PrizeCodes `json:"lottery_results"`
}
type CheckLotteryResultRequest struct {
	LotteryDateRange []string
	IsAdvanced       bool
	LotteryName      string
	SeriesName       string
	LotteryCodes     []string
}
type AnalyzeLotteryResultRequest struct {
	LotteryDateRange []string
	LotteryName      string
	LotteryCodes     []string
}
type EvaluateResultsResponse struct {
	SeriesName    string `json:"series_name,omitempty"`
	LotteryDate   string `json:"series_date,omitempty"`
	LotteryTime   string `json:"series_time,omitempty"`
	PrizePosition string `json:"prize_position"`
	PrizeMoney    string `json:"prize_money"`
	WinnerCode    string `json:"winner_code"`
}
type CheckLotteryResultResponse struct {
	SeriesName  string                    `json:"series_name,omitempty"`
	LotteryDate time.Time                 `json:"series_date,omitempty"`
	LotteryTime string                    `json:"series_time,omitempty"`
	SeriesLink  string                    `json:"series_link,omitempty"`
	Results     []EvaluateResultsResponse `json:"results"`
}
type AnalyzeLotteryResultResponse struct {
	LotteryName string                       `json:"lottery_name,omitempty"`
	Results     []CheckLotteryResultResponse `json:"results"`
}
type GetLotteriesRequest struct {
	LotteryName string
}
type PdfData struct {
	Name string
	Date string
	Link string
}
type HelloResponse struct {
	Message string `json:"message"`
}
