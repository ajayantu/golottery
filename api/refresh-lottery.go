package handler

import (
	"lotteryapi/db"
	"lotteryapi/domain"
	"lotteryapi/helpers"
	"net/http"
)

func RefreshResults(w http.ResponseWriter, r *http.Request) {

	pdfdatas, err := helpers.ExtractPdfLink("")
	if err != nil || len(pdfdatas) == 0 {
		helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
			Message:    "error in extracting pdf link",
			ErrorField: "get pdf link",
		}})
		return
	}
	collection := db.ConnectDB()
	myresults := db.GetLatestResult(collection)

	var results []domain.GetLotteryResultRespose
	for _, item := range pdfdatas {
		isGreater, err := helpers.CompareDates(item.Date, myresults.LotteryDate)
		if err != nil {
			helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
				Message:    err.Error(),
				ErrorField: "date parsing",
			}})
			return
		}
		if isGreater || myresults.LotteryName == "" {
			result, err := helpers.ExtractResultsFromLink(item.Name, item.Link)
			if err != nil {
				helpers.Fail(w, http.StatusInternalServerError, []helpers.FailStruct{{
					Message:    err.Error(),
					ErrorField: "pdf result extraction",
				}})
				return
			}
			db.CreateResult(collection, result)
		}
	}

	helpers.Success(w, 200, results)

}
