package helpers

import (
	"fmt"
	"lotteryapi/domain"
	"net/http"
	"strings"
)

func ParseGetResultRequestParam(r *http.Request) (domain.GetLotteryResultRequest, error) {
	reqParams := domain.GetLotteryResultRequest{}
	seriesName := r.URL.Query().Get("series_name")
	if seriesName == "" {
		return domain.GetLotteryResultRequest{}, fmt.Errorf("lottery name required")
	}
	reqParams.SeriesName = seriesName
	return reqParams, nil
}
func ParseCheckResultsRequestParam(r *http.Request) (domain.CheckLotteryResultRequest, error) {
	reqParams := domain.CheckLotteryResultRequest{}
	seriesName := r.URL.Query().Get("series_name")
	if seriesName == "" {
		return domain.CheckLotteryResultRequest{}, fmt.Errorf("series name required")
	}
	reqParams.SeriesName = seriesName

	queryValues := r.URL.Query()
	lotteryCodesQuery := queryValues["lottery_codes"]
	if len(lotteryCodesQuery) > 0 {
		lotteryCodes := strings.Split(lotteryCodesQuery[0], ",")
		reqParams.LotteryCodes = lotteryCodes
	}
	return reqParams, nil
}
func ParseAnalyzeResultsRequestParam(r *http.Request) (domain.AnalyzeLotteryResultRequest, error) {
	reqParams := domain.AnalyzeLotteryResultRequest{}
	lotteryName := r.URL.Query().Get("lottery_name")

	reqParams.LotteryName = lotteryName
	queryValues := r.URL.Query()

	lotteryCodesQuery := queryValues["lottery_codes"]
	if len(lotteryCodesQuery) > 0 {
		lotteryCodes := strings.Split(lotteryCodesQuery[0], ",")
		reqParams.LotteryCodes = lotteryCodes
	}
	return reqParams, nil
}
