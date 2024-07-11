package handler

import (
	"lotteryapi/db"
	"lotteryapi/domain"
	"lotteryapi/helpers"
	"net/http"
)

func AnalyzeResults(w http.ResponseWriter, r *http.Request) {
	reqParams, err := helpers.ParseAnalyzeResultsRequestParam(r)
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
			Message:    "error in extracting pdf links",
			ErrorField: "get pdf link",
		}})
		return
	}
	var finalResults domain.AnalyzeLotteryResultResponse
	collection := db.ConnectDB()
	dbResults := db.GetMyAllResults(collection)
	pdfMap := helpers.MapPdfResultToName(dbResults)
	if reqParams.LotteryName != "" {
		//specific lottery's all series
		mapedData := helpers.MapPdfDatas(pdfdatas)
		if _, ok := mapedData[reqParams.LotteryName]; ok {
			finalResults, err = helpers.EvaluateAllLotteries(mapedData[reqParams.LotteryName], reqParams.LotteryCodes, pdfMap)
			if err != nil {
				helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
					Message:    err.Error(),
					ErrorField: "evaluating all result",
				}})
				return
			}
		} else {
			finalResults = domain.AnalyzeLotteryResultResponse{}
		}
	} else {
		//all lottery all series
		finalResults, err = helpers.EvaluateAllLotteries(pdfdatas, reqParams.LotteryCodes, pdfMap)
		if err != nil {
			helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
				Message:    err.Error(),
				ErrorField: "evaluating all result",
			}})
			return
		}
	}
	helpers.Success(w, 200, finalResults)
}
