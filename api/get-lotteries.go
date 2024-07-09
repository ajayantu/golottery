package handler

import (
	"lotteryapi/helpers"
	"net/http"
)

func GetLotteries(w http.ResponseWriter, r *http.Request) {
	reqParam, err := helpers.ParseGetLotteriesRequestParam(r)
	if err != nil {
		helpers.Fail(w, http.StatusBadRequest, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "parsing body",
		}})
		return
	}
	pdfdatas, err := helpers.ExtractPdfLink("")
	if err != nil || len(pdfdatas) == 0 {
		helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "get pdf link",
		}})
		return
	}
	finalData := pdfdatas
	if reqParam.LotteryName != "" {
		mapedData := helpers.MapPdfDatas(pdfdatas)
		finalData = mapedData[reqParam.LotteryName]
	}
	helpers.Success(w, 200, finalData)
}
