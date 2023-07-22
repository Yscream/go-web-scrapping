package dou

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Yscream/go-webScrapping/model"
	"github.com/go-rod/rod"
	"github.com/xuri/excelize/v2"
)

const (
	dou = "https://jobs.dou.ua/companies/?name=%D0%94%D0%BD%D1%96%D0%BF%D1%80%D0%BE"
)

func WebScrapingDOU() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(dou).MustWaitLoad()

	companies := assembleCompanyModels(page)
	for _, d := range companies {
		fmt.Println(d)
	}

	writeResutlXls(companies)
}

func assembleCompanyModels(page *rod.Page) []*model.Company {
	assemblyCompanies := make([]*model.Company, 0)

	scrappedCompanyNames, links := scarpingCompanyNamesAndDOULinks(page)
	scrappedCompanyDescr := scarpingCompanyDescripitons(page)

	log.Println("starting assemble company")
	for i, name := range scrappedCompanyNames {
		log.Println("Current URL:", page.MustInfo().URL)
		log.Println("check:", links[i], "and", name)
		navPage := page.MustNavigate(links[i]).MustWaitLoad()
		assemblyCompany := &model.Company{
			Name:             name,
			Type:             "",
			Url:              scrapingCompanySites(navPage),
			Desctiption:      scrappedCompanyDescr[i],
			EmployeeQuantity: scarpingQuantityOfCompanyEmployees(navPage),
			Address:          scarpingCompanyOfficeAddresses(navPage, links[i]),
		}
		page.MustNavigateBack().MustWaitLoad()
		fmt.Println("!!!!!", assemblyCompany)
		assemblyCompanies = append(assemblyCompanies, assemblyCompany)
	}
	log.Println("finished assemble company")

	return assemblyCompanies
}

func scarpingCompanyNamesAndDOULinks(page *rod.Page) ([]string, []string) {
	log.Println("starting render list of company")

	moreCompaniesButton := page.MustElement(".more-btn a")
	for {
		log.Println("click more button")

		checkDisplay := moreCompaniesButton.MustEval(`() => this.style.display`).String()
		if checkDisplay == "none" {
			break
		}

		moreCompaniesButton.MustClick()
		for {
			class, err := moreCompaniesButton.Attribute("class")
			if err != nil {
				log.Fatal(err)
			}

			if *class != "__loading" {
				break
			}
		}
	}

	companies := make([]string, 0)
	newCompanies := page.MustElements(".company .ovh .h2 .cn-a")
	for _, company := range newCompanies {
		companies = append(companies, company.MustText())
	}
	log.Println("finished render list of company")

	//i return 2 value and the seconde one is also slice of strings. it is links. i understand that looks dirty, but in my opinion
	//this is better than return two same companies but in different types
	return companies, scarpingCompanyDOULinks(newCompanies)
}

func scarpingCompanyDescripitons(page *rod.Page) []string {
	descrs := []string{}

	log.Println("starting collect info about descrp")
	descrElems := page.MustElements(".company .ovh .descr")
	
	for i := range descrElems {
		descrs = append(descrs, descrElems[i].MustText())
	}
	log.Println("finished collect info about descrp")

	return descrs
}

func scarpingCompanyDOULinks(companies rod.Elements) []string {
	links := []string{}

	for _, company := range companies {
		href, err := company.Attribute("href")
		if err != nil {
			log.Fatal(err)
		}

		str := *href
		links = append(links, str)
	}
	log.Println("done scarp dou links")

	return links
}

func scrapingCompanySites(page *rod.Page) string {
	if !isPageExists(page) {
		log.Println("page doesnt exist. returns empty string")
		return ""
	}
	log.Println("page exists")

	log.Println("starting collect info about sites")
	companyInfo := page.MustElement(".company-info")
	if !companyInfo.MustHas(".site a") {
		return ""
	}

	site := companyInfo.MustElement(".site a")

	href, err := site.Attribute("href")
	if err != nil {
		log.Fatal(err)
	}

	str := *href
	log.Println("finished collect info about sites")
	return str
}

func scarpingQuantityOfCompanyEmployees(page *rod.Page) string {
	//remark: due to the fact that info/text about quantity of employees even doesn't contain inside any div or other tags/attributes i need to fetch all html page
	// and find this between two div with classes "g-h2" and "offices"
	if !isPageExists(page) {
		log.Println("page doesnt exist. returns empty string")
		return ""
	}
	log.Println("page exists")

	log.Println("started collect info about employees")
	html := page.MustHTML()

	re := regexp.MustCompile(`<h1 class="g-h2">[^<]*</h1>\s*([^<]*)\s*<div class="offices">`)
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		log.Println("didn't find anything")
		return ""
	}
	log.Println("finshed collect info about employees")

	return match[1]
}

func scarpingCompanyOfficeAddresses(page *rod.Page, link string) string {
	suffix := "offices/"
	link += suffix

	page.MustNavigate(link).MustWaitLoad()

	log.Println("starting collect address")
	cities := page.MustElements(".city")
	fmt.Println("allo chto proishodit")

	for _, city := range cities {
		if city.MustHas("#dnepr") {
			if city.MustHas(".contacts") {
				if city.MustHas(".address") {
					return city.MustElement(".contacts .address").MustText()
				}
			}
		}
	}
	log.Println("finshed collect address")

	return ""
}

func isPageExists(page *rod.Page) bool {
	log.Println("checking page existence")
	s, err := page.HTML()
	if err != nil {
		log.Fatal(err)
	}

	return !strings.Contains(s, "Такої сторінки немає")
}

func writeResutlXls(companies []*model.Company) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	f.NewSheet("Sheet1")

	for i, company := range companies {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i+1), company.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i+1), company.Type)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i+1), company.Url)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i+1), company.Desctiption)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%v", i+1), company.EmployeeQuantity)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%v", i+1), company.Address)
	}

	log.Println("generating xlsx...")
	err := f.SaveAs("./dou.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")
}
