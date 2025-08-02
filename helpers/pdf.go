package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"os"

	"lotteryapi/domain"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/romanpickl/pdf"
)

func ExtractTextFromPDF(fileUrl string) (string, error) {
	// Fetch the PDF content
	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("error in fetching pdf")
	}
	defer resp.Body.Close()

	var textContent strings.Builder
	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading the pdf response")
	}
	r, err := pdf.NewReader(bytes.NewReader(pdfBytes), int64(len(pdfBytes)))
	if err != nil {
		return "", fmt.Errorf("error in creating pdf reader")
	}

	// Extract text from each page and split into rows
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		rows, err := p.GetTextByRow()
		if err != nil {
			fmt.Println(err)
		}
		for _, row := range rows {
			for _, word := range row.Content {
				textContent.WriteString(word.S)
				textContent.WriteString(" ")
			}
		}
	}
	return textContent.String(), nil
}

func ExtractResults(seriesName string, seriesLink string, pdfText string) (domain.GetLotteryResultRespose, error) {
	extractResultString := `((?:\d+\S{2}|Cons)\sPrize)>(.*?)(<|EOF)`
	splitResultString := `([A-Z]{1,2} \d+)|(\d{4})`
	extractResultRegex := regexp.MustCompile(extractResultString)
	splitResultRegex := regexp.MustCompile(splitResultString)

	resultsMap := make(map[string]domain.PrizeCodes)

	cleanText := CleanPdfText(pdfText)

	//extract prizes
	prizesMap := make(map[string]string)
	extractPrizeMoneyString := `(((?:\d+\S{2}|Cons)\sPrize)(?: Rs|-Rs) :(\d+))`
	extractPrizeMoneyRegex := regexp.MustCompile(extractPrizeMoneyString)
	prizes := extractPrizeMoneyRegex.FindAllStringSubmatch(pdfText, -1)

	for _, prize := range prizes {
		if len(prize) < 4 {
			for k := range prizesMap {
				delete(prizesMap, k)
			}
			break
		}
		prizesMap[prize[2]] = prize[3]
	}

	//extract results and add prizes
	matches := extractResultRegex.FindAllStringSubmatch(cleanText, -1)
	for _, item := range matches {
		if len(item) < 4 {
			return domain.GetLotteryResultRespose{}, fmt.Errorf("error in extracting lottery results")
		}
		splitResult := splitResultRegex.FindAllString(item[2], -1)
		resultsMap[item[1]] = domain.PrizeCodes{PrizeMoney: prizesMap[item[1]], PrizeCodes: splitResult}
	}

	//extract date and time
	dateTimeString := `DRAW held on:-  (\d{1,2}\/\d{1,2}\/\d{4}),(\d{1,2}:\d{2} [AP]M)`
	dateTimeRegex := regexp.MustCompile(dateTimeString)
	dateTime := dateTimeRegex.FindAllStringSubmatch(pdfText, -1)

	//add to domain struct
	parsedDate, err := time.Parse("02/01/2006", dateTime[0][1])
	if err != nil {
		return domain.GetLotteryResultRespose{}, fmt.Errorf("error in parsing date")
	}
	if len(dateTime) >= 1 && len(dateTime[0]) >= 3 {
		finalResults := domain.GetLotteryResultRespose{LotteryName: seriesName, LotteryLink: seriesLink, LotteryDate: parsedDate, LotteryTime: dateTime[0][2], LotteryResults: resultsMap}
		return finalResults, nil
	}
	finalResults := domain.GetLotteryResultRespose{LotteryName: seriesName, LotteryLink: seriesLink, LotteryTime: "", LotteryResults: resultsMap}
	return finalResults, nil
}

func CleanPdfText(pdfText string) string {
	cleanText := pdfText
	const (
		removeExtraSpace   = `\s{2,}`
		footerRegex        = `Page \d+ IT Support : NIC Kerala \d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}`
		forTheTicketEnding = `FOR THE TICKETS ENDING WITH THE FOLLOWING NUMBERS`
		thePrizeWinners    = `(?s)The prize winners.*`
		keralState         = `KERALA STATE LOTTERIES(.*?)THIRUVANANTHAPURAM`
		indexNumbers       = `\d+\)`
		winnerPlace        = `\(\S+\)`
		prizeMoney         = `(\s|-)Rs\s:\d+\/-`
		paraEnclose        = `((?:\d+\S{2}|Cons)\sPrize)`
	)
	cleanTextRegex := []string{removeExtraSpace, footerRegex, forTheTicketEnding, thePrizeWinners, keralState, indexNumbers, winnerPlace, prizeMoney, removeExtraSpace}

	for _, regexString := range cleanTextRegex {
		regexPattern := regexp.MustCompile(regexString)
		cleanText = regexPattern.ReplaceAllString(cleanText, " ")
	}
	regexPattern := regexp.MustCompile(paraEnclose)
	cleanText = regexPattern.ReplaceAllString(cleanText, "<$1>")
	cleanText += "EOF"
	return cleanText
}

func ExtractPdfLink(seriesName string) ([]domain.PdfData, error) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	var pdfdatas []domain.PdfData

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		href := e.ChildAttr("td a", "href")
		name := e.ChildText("td:nth-child(1)")
		date := e.ChildText("td:nth-child(2)")
		if name != "" && href != "" && date != "" {
			if seriesName == "" {
				pdfdatas = append(pdfdatas, domain.PdfData{Name: name, Link: href, Date: date})
			} else {
				if seriesName == name {
					pdfdatas = append(pdfdatas, domain.PdfData{Name: name, Link: href, Date: date})
					return
				}
			}
		}

	})
	err := c.Visit(os.Getenv("LOTTERY_URl"))
	if err != nil {
		return []domain.PdfData{}, fmt.Errorf("error in scraping the results")
	}
	return pdfdatas, nil
}
