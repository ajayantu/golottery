package handler

import (
	"lotteryapi/db"
	"lotteryapi/domain"
	"lotteryapi/helpers"
	"net/http"
)

func CheckResults(w http.ResponseWriter, r *http.Request) {
	reqParams, err := helpers.ParseCheckResultsRequestParam(r)
	if err != nil {
		helpers.Fail(w, http.StatusBadRequest, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "parsing body",
		}})
		return
	}
	pdfdatas, err := helpers.ExtractPdfLink(reqParams.SeriesName)
	if err != nil || len(pdfdatas) == 0 {
		helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "get pdf link",
		}})
		return
	}
	var finalResults domain.CheckLotteryResultResponse

	collection := db.ConnectDB()
	dbResults := db.GetMyAllResults(collection)
	pdfMap := helpers.MapPdfResultToName(dbResults)
	finalResults, err = helpers.EvaluateResultsFromLink(pdfdatas[0].Name, pdfdatas[0].Link, reqParams.LotteryCodes, pdfMap, reqParams.Templating)
	if err != nil {
		helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "evaluating result",
		}})
		return
	}
	helpers.Success(w, 200, finalResults)
}
