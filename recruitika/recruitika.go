package recruitika

import (
	"fmt"
	"strings"

	"github.com/Yscream/go-webScrapping/model"
	"github.com/go-rod/rod"
)

const (
	recruitika = "https://recruitika.com/companies/"
)

func WebScrapingRecruitika() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(recruitika).MustWaitLoad()
	assembleCompanyModels(page)
}

func assembleCompanyModels(page *rod.Page) []*model.Company {
	assemblyCompanies := make([]*model.Company, 0)
	companies := sortCompaniesByCityDnipro(page)
	fmt.Println(companies)
	// scrappedCompanyNames := scarpingCompanyNames(companies)
	// fmt.Println(scrappedCompanyNames)
	// companies := make([]*model.Company, 0)
	// fmt.Println(len(scrappedCompanyNames))

	// for _, name := range scrappedCompanyNames {
	// 	assemblyCompany := &model.Company{
	// 		Name:             name,
	// 		Type:             "",
	// 		Url:              "",
	// 		Desctiption:      "",
	// 		EmployeeQuantity: "",
	// 		Address:          "",
	// 	}
	// 	assemblyCompanies = append(assemblyCompanies, assemblyCompany)
	// 	// arrow := page.MustElement(".companies-list .paginator .simple-arrow.next")
	// 	// linkPage := arrow.MustAttribute("href")
	// 	// fmt.Println(*linkPage)
	// 	// page.MustNavigate("https://recruitika.com" + *linkPage)

	// }

	// div.MustInput("Дніпро")
	// fmt.Println(div.MustText())

	// a := []string{}
	// cities := page.MustElements(".company-item .characteristics-item .item-value span")

	//
	// fmt.Println(*link)

	// for {
	// 	// for i, city := range cities {
	// 	// 	fmt.Println(i)
	// 	// 	if strings.Contains(city.MustText(), "Дніпро") {
	// 	// 		city.MustText()
	// 	// 		fmt.Println(i, city.MustText())
	// 	// 		a = append(a, city.MustText())
	// 	// 	}
	// 	// }

	// 	// arrow.MustClick()
	// }
	return assemblyCompanies
}

func sortCompaniesByCityDnipro(page *rod.Page) rod.Elements {
	dniproCompanies := rod.Elements{}

main:
	for {
		elems := page.MustElements(".company-item")
		for i, elem := range elems {
			fmt.Println(page.MustInfo().URL)
			cities := elem.MustElements(".characteristics-item .item-value")
			for _, city := range cities {
				if strings.Contains(city.MustText(), "Дніпро") {
					fmt.Println(i)
					dniproCompanies = append(dniproCompanies, elem)
					continue
				}
			}
			if i == len(elems)-1 {
				arrow := page.MustElement(".companies-list .paginator .simple-arrow.next")
				if arrow.MustHas(".disabled") || arrow.MustHas(".disabled.simple-arrow.next") {
					break main
				}
				fmt.Println("privet", arrow)
				linkPage := arrow.MustAttribute("href")
				fmt.Println("privet", *linkPage)
				page.MustNavigate("https://recruitika.com" + *linkPage)
			}
		}
	}

	return dniproCompanies
}

func scarpingCompanyNames(elems rod.Elements) []string {
	companies := make([]string, 0)
	for _, elem := range elems {
		company := elem.MustElement(".image-name-block .name-rating .name")
		companies = append(companies, company.MustText())
	}
	// newCompanies := page.MustElements(".company-item .image-name-block .name-rating .name")
	// for _, company := range newCompanies {
	// 	
	// }

	return companies
}

// func scarpingCompanyTypes(page *rod.Page) [] {

// }

// // var companies []Company
// c := colly.NewCollector()

// c.OnHTML(".company-item", func(h *colly.HTMLElement) {
// 	// var test []string
// 	div := h.DOM
// 	divCompany := div.Find(".image-name-block").Find(".name-rating").Find(".name").Text()
// 	fmt.Println(divCompany)
// 	// divCompany := div.Find(".company-item")
// 	divCharacteristics := div.Find(".characteristics-block").Text()
// 	// divTeam := divCharacteristics.Find(".team").Text()
// 	// // divItemLabel := divCharacteristics.Find(".characteristics-item").Text()
// 	fmt.Println(divCharacteristics)
// 	// a := strings.Split(divTeam, "Команда:")
// 	// fmt.Println(a)

// 	// test = append(test, divItemLabel)

// 	// splitted :=  strings.Split(divItemLabel, ":")
// 	// fmt.Println(splitted)
// 	// joined := strings.Join(splitted, ",")
// 	// switch {
// 	// case divItemLabel == "Тип:":
// 	// 	fmt.Println(1)
// 	// case divItemLabel == "Теги:":
// 	// 	fmt.Println(2)
// 	// case divItemLabel == "Сайт:":
// 	// 	fmt.Println(3)
// 	// case divItemLabel == "Офіси:":
// 	// 	fmt.Println(4)
// 	// }

// 	// divName := divCompany.Find(".name").Text()
// 	// fmt.Println(divName)

// 	// paginator := div.Find(".paginator")
// 	// href, _ := paginator.Find(".simple-arrow.next").Attr("href")

// 	// baseURL := "https://recruitika.com"
// 	// fullURL := baseURL + href

// 	// fmt.Println("Ссылка:", fullURL)

// 	// err := c.Visit(fullURL)
// 	// if err != nil {
// 	// 	log.Println("Ошибка при переходе по ссылке:", err)
// 	// }
// })

// c.Visit("https://recruitika.com/companies/")
