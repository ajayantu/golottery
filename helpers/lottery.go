package helpers

import (
	"fmt"
	"lotteryapi/domain"
	"strings"
)

func ExtractResultsFromLink(seriesName string, pdfLink string) (domain.GetLotteryResultRespose, error) {
	pdfText, err := ExtractTextFromPDF(pdfLink)
	if err != nil {
		return domain.GetLotteryResultRespose{}, fmt.Errorf("error in extracting text from pdf")
	}
	results, err := ExtractResults(seriesName, pdfLink, pdfText)
	if err != nil {
		return domain.GetLotteryResultRespose{}, fmt.Errorf("error in extracting text from pdf")
	}
	return results, nil
}

func EvaluateResultsFromLink(seriesName string, pdfLink string, lotteryCodes []string, pdfMap map[string]domain.GetLotteryResultRespose) (domain.CheckLotteryResultResponse, error) {
	dbResults := pdfMap[seriesName]
	var results domain.GetLotteryResultRespose
	var err error
	if dbResults.LotteryName != "" {
		results = dbResults
	} else {
		results, err = ExtractResultsFromLink(seriesName, pdfLink)
		if err != nil {
			return domain.CheckLotteryResultResponse{}, err
		}
	}

	//evaluate result
	evaluationResults := []domain.EvaluateResultsResponse{}
	for key, item := range results.LotteryResults {
		for _, lotteryCode := range lotteryCodes {
			if StringInSlice(lotteryCode, item.PrizeCodes) {

				evaluationResults = append(evaluationResults,
					domain.EvaluateResultsResponse{
						PrizePosition: key,
						PrizeMoney:    item.PrizeMoney,
						WinnerCode:    lotteryCode,
					},
				)
			}
		}
	}

	//create final response
	finalResponse := domain.CheckLotteryResultResponse{
		SeriesName:  seriesName,
		LotteryDate: results.LotteryDate,
		LotteryTime: results.LotteryTime,
		SeriesLink:  pdfLink,
		Results:     evaluationResults,
	}
	return finalResponse, nil
}

func EvaluateAllLotteries(pdfDatas []domain.PdfData, lotteryCodes []string, pdfMap map[string]domain.GetLotteryResultRespose) (domain.AnalyzeLotteryResultResponse, error) {
	var finalResults domain.AnalyzeLotteryResultResponse

	//create the lottery codes struct
	lotteryCodesMap := make(map[string][]string)
	for _, item := range lotteryCodes {
		switch item[0] {
		case 'F':
			lotteryCodesMap["FIFTY-FIFTY"] = append(lotteryCodesMap["FIFTY-FIFTY"], item)
		case 'S':
			lotteryCodesMap["STHREE-SAKTHI"] = append(lotteryCodesMap["STHREE-SAKTHI"], item)
		case 'W':
			lotteryCodesMap["WIN-WIN"] = append(lotteryCodesMap["WIN-WIN"], item)
		case 'A':
			lotteryCodesMap["AKSHAYA"] = append(lotteryCodesMap["AKSHAYA"], item)
		case 'K':
			lotteryCodesMap["KARUNYA(KR"] = append(lotteryCodesMap["KARUNYA(KR"], item)
		case 'N':
			lotteryCodesMap["NIRMAL"] = append(lotteryCodesMap["NIRMAL"], item)
		case 'P':
			lotteryCodesMap["KARUNYA PLUS"] = append(lotteryCodesMap["KARUNYA PLUS"], item)
		default:
			lotteryCodesMap["All"] = append(lotteryCodesMap["All"], item)
		}
	}
	//loop through pdfdatas
	for _, item := range pdfDatas {
		for key, val := range lotteryCodesMap {
			if strings.Contains(item.Name, key) || key == "All" {
				//check if item/lotterypdf contains  lotterycodes from map
				results, err := EvaluateResultsFromLink(item.Name, item.Link, val, pdfMap)
				if err != nil {
					return domain.AnalyzeLotteryResultResponse{}, err
				}
				if len(results.Results) > 0 {
					finalResults.Results = append(finalResults.Results, results)
				}

			}
		}
	}
	return finalResults, nil
}
func StringInSlice(a string, list []string) bool {
	for _, item := range list {
		updatedCode := a
		if len(item) == 4 {
			updatedCode = a[len(a)-4:]
		}
		if item == updatedCode {
			return true
		}
	}
	return false
}
func MapPdfDatas(pdfDatas []domain.PdfData) map[string][]domain.PdfData {
	mapedData := make(map[string][]domain.PdfData)
	for _, pdfdata := range pdfDatas {
		switch {
		case strings.Contains(pdfdata.Name, "AKSHAYA"):
			mapedData["AK"] = append(mapedData["AK"], pdfdata)
		case strings.Contains(pdfdata.Name, "KARUNYA(KR"):
			mapedData["KR"] = append(mapedData["KR"], pdfdata)
		case strings.Contains(pdfdata.Name, "NIRMAL"):
			mapedData["NR"] = append(mapedData["NR"], pdfdata)
		case strings.Contains(pdfdata.Name, "KARUNYA PLUS"):
			mapedData["KN"] = append(mapedData["KN"], pdfdata)
		case strings.Contains(pdfdata.Name, "FIFTY-FIFTY"):
			mapedData["FF"] = append(mapedData["FF"], pdfdata)
		case strings.Contains(pdfdata.Name, "STHREE-SAKTHI"):
			mapedData["SS"] = append(mapedData["SS"], pdfdata)
		case strings.Contains(pdfdata.Name, "WIN-WIN"):
			mapedData["W"] = append(mapedData["W"], pdfdata)
		}
	}
	return mapedData
}
func MapPdfResultToName(pdfResults []domain.GetLotteryResultRespose) map[string]domain.GetLotteryResultRespose {
	pdfMap := make(map[string]domain.GetLotteryResultRespose)
	for _, item := range pdfResults {
		pdfMap[item.LotteryName] = item
	}
	return pdfMap
}
