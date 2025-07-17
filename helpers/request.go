package helpers

import (
	"fmt"
	"lotteryapi/domain"
	"net/http"
	"strconv"
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
	var err error

	seriesName := r.URL.Query().Get("series_name")
	if seriesName == "" {
		return domain.CheckLotteryResultRequest{}, fmt.Errorf("series name required")
	}
	reqParams.SeriesName = seriesName

	templating := r.URL.Query().Get("templating")
	if templating != "" {
		reqParams.Templating, err = strconv.ParseBool(templating)
		if err != nil {
			return domain.CheckLotteryResultRequest{}, fmt.Errorf("error in parsing params")
		}
	}

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
	var err error

	lotteryName := r.URL.Query().Get("lottery_name")
	reqParams.LotteryName = lotteryName

	templating := r.URL.Query().Get("templating")
	if templating != "" {
		reqParams.Templating, err = strconv.ParseBool(templating)
		if err != nil {
			return domain.AnalyzeLotteryResultRequest{}, fmt.Errorf("error in parsing params")
		}
	}

	queryValues := r.URL.Query()
	lotteryCodesQuery := queryValues["lottery_codes"]
	if len(lotteryCodesQuery) > 0 {
		lotteryCodes := strings.Split(lotteryCodesQuery[0], ",")
		reqParams.LotteryCodes = lotteryCodes
	}
	return reqParams, nil
}
func ParseGetLotteriesRequestParam(r *http.Request) (domain.GetLotteriesRequest, error) {
	reqParams := domain.GetLotteriesRequest{}
	lotteryName := r.URL.Query().Get("lottery_name")
	reqParams.LotteryName = lotteryName
	return reqParams, nil
}
