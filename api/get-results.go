package handler

import (
	"lotteryapi/helpers"
	"net/http"
)

func GetResults(w http.ResponseWriter, r *http.Request) {
	reqBody, err := helpers.ParseGetResultRequestParam(r)
	if err != nil {
		helpers.Fail(w, http.StatusBadRequest, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "parsing body",
		}})
		return
	}
	pdfdatas, err := helpers.ExtractPdfLink(reqBody.SeriesName)
	if err != nil || len(pdfdatas) == 0 {
		helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
			Message:    "error in extracting pdf link",
			ErrorField: "get pdf link",
		}})
		return
	}
	results, err := helpers.ExtractResultsFromLink(pdfdatas[0].Name, pdfdatas[0].Link)
	if err != nil {
		helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
			Message:    err.Error(),
			ErrorField: "pdf result extraction",
		}})
		return
	}
	helpers.Success(w, 200, results)
}
