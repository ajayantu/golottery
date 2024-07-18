package helpers

import (
	"fmt"
	"lotteryapi/domain"
	"regexp"
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
	lotteryMap := map[byte]string{
		'F': "FIFTY-FIFTY",
		'S': "STHREE-SAKTHI",
		'W': "WIN-WIN",
		'A': "AKSHAYA",
		'K': "KARUNYA(KR",
		'N': "NIRMAL",
		'P': "KARUNYA PLUS",
	}
	lotteryCodesMap := make(map[string][]string)
	for _, item := range lotteryCodes {
		trimedItem := strings.TrimSpace(item)
		if trimedItem != "" {
			if _, ok := lotteryMap[trimedItem[0]]; ok {
				lotteryCodesMap[lotteryMap[trimedItem[0]]] = append(lotteryCodesMap[lotteryMap[trimedItem[0]]], trimedItem)
			} else {
				lotteryCodesMap["All"] = append(lotteryCodesMap["All"], trimedItem)
			}
		}
	}
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

func StringInSlice(inputCode string, list []string) bool {
	trimmedCode := strings.TrimSpace(inputCode)
	if !MatchFormat(trimmedCode) {
		return false
	}
	for _, item := range list {
		updatedCode := trimmedCode
		if len(item) == 4 && len(inputCode) >= 4 {
			updatedCode = inputCode[len(inputCode)-4:]
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
func MatchFormat(input string) bool {
	format1 := regexp.MustCompile(`^[A-Z]{2} \d{6}$`)
	format2 := regexp.MustCompile(`^\d{4}$`)
	return format1.MatchString(input) || format2.MatchString(input)
}
