package utils

import (
	"errors"
	"log"

	"github.com/sclevine/agouti"
)

func StartChromeSession(url string, implicitwait int) (*agouti.Page, error) {
	go func() {
		log.Println("Starting browsing session...")
	}()
	//driver := agouti.ChromeDriver()
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{"--no-sandbox", "--ignore-certificate-errors"}), //[]string{"--headless", "--disable-gpu", "--no-sandbox"}), //[]string{"--headless", "--disable-gpu", "--no-sandbox"}
	)
	if err := driver.Start(); err != nil {
		log.Fatal(err)
	}
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatal(err)
	}
	page.Session().SetImplicitWait(implicitwait)
	if err := page.Navigate(url); err != nil {
		return nil, err
	}
	return page, nil
}

func FindXpathElemOnPage(page *agouti.Page, patterntofind, waytofind string) (*agouti.Selection, error) {
	var selection *agouti.Selection
	switch waytofind {
	case "simple":
		selection = page.Find(patterntofind)
	case "label":
		selection = page.FindByLabel(patterntofind)
	}

	if selection == nil {
		return nil, errors.New("This element was not found\n")
	}
	return selection, nil
}

func FillForm(element *agouti.Selection, stuffToFillWith string) error {
	if err := element.Fill(stuffToFillWith); err != nil {
		return err
	}
	return nil
}

func ClickButton(buttonPattern string) error {

	return nil
}

func ErrHandler(err error, message string) {
	if err != nil {
		log.Fatalf("%s %v\n", message, err)
	}
}
