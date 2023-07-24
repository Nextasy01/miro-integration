package handler

import (
	"fmt"
	"log"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath     = "webdrivers/selenium-server-standalone-3.5.3.jar"
	chromeDriverPath = "webdrivers/chromedriver.exe"
	port             = 4444
)

func StartSelenium(email, password string) (string, error) {
	// opts := []selenium.ServiceOption{
	// 	// selenium.StartFrameBuffer(),
	// 	selenium.ChromeDriver("webdrivers/chromedriver.exe"),
	// 	selenium.Output(os.Stderr),
	// }

	// selenium.SetDebug(false)
	// service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	// if err != nil {
	// 	log.Println("Could not start selenium service")
	// 	return "server error", err
	// }
	// defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Println("Could not start web driver")
		return "server error", err
	}
	defer wd.Quit()

	if err := wd.Get("https://miro.com/login/"); err != nil {
		log.Println("Could not get page")
		return "server error", err
	}

	elem, err := wd.FindElements(selenium.ByCSSSelector, ".signup__input-text")
	if err != nil {
		log.Println("Could not get element")
		return "server error", err
	}

	if err := elem[0].Clear(); err != nil {
		return "server error", err
	}

	err = elem[0].SendKeys(email)
	if err != nil {
		return "server error", err
	}

	if err := elem[1].Clear(); err != nil {
		return "server error", err
	}

	err = elem[1].SendKeys(password)
	if err != nil {
		return "server error", err
	}

	btn, err := wd.FindElement(selenium.ByCSSSelector, ".signup__submit")
	if err != nil {
		return "server error", err
	}
	if err := btn.Click(); err != nil {
		return "server error", err
	}

	token, err := wd.GetCookie("token")
	if err != nil {
		return "server error", err
	}

	return token.Value, err
}
