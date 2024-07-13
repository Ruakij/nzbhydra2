package main

import (
	"encoding/base64"
	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
	"os"
	"tawesoft.co.uk/go/dialog"
	"theotherp/base"
	"time"
)

const IcoBytesB64 = "AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAACMuAAAjLgAAAAAAAAAAAAAAAAD8DxMO/i05J/4sOyL+NVAf/kZxH/5VjyP+WJwg/0iJEf4zYwn+FygK/gEBAf4AAAD+AAAA/gAAAP4AAAD8AAAA/gICAv8gKBz/Qlc1/0RfMP82SyT/LUAd/zJLH/9DbST/R44P/1CaFP8qQxj/AAAB/wAAAP8AAAD/AAAA/gAAAP4AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/BAUD/yM/Dv9NnQn/Z6cy/xciD/8AAAD/AAAA/wAAAP4AAAD+AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8KDwb/RIwE/3C7MP85UyX/AAAA/wAAAP8AAAD+AAAA/gAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8GCQX/KksM/1CkAv9xuzH/QV0r/wAAAP8AAAD/AAEA/gAAAP4AAAD/AAAA/wAAAP8AAAD/AAAA/wIDAv8qPSH/Wo03/220Mv9gsBj/ebw+/y9CIP8AAAD/AAAA/wAAAP4AAAD+AAAA/wAAAP8AAAD/AAAA/wAAAP84STL/jbx4/4m6bv9dgkT/XJgr/3yyTv8THA3/AAAA/wAAAP8AAAD+AAAA/gAAAP8ICQn/QUVB/xkaGf8RFBD/i6qB/5e8if+Lqn3/OEwv/3u1Uv9egET/AQIB/wAAAP8AAAD/AAAA/gABAP4AAAD/ERMR/5qnl/+erJr/j6GK/6bAnv+cv5D/fJF0/zdMLf+NwHH/Kzkj/wAAAP8AAAD/AAAA/wAAAP4AAQD+AAAA/wAAAP88QDv/k5+P/4iThf81OjP/k6yL/2Z1Yv82SDD/hat0/w4RDf8AAAD/AAAA/wAAAP8AAAD+AAAA/gAAAP8AAAD/AAAA/wsLC/9cYlv/Ghsa/4CQfP9gal3/ISoe/46ygP8sNij/HCYX/xkkFP8CAwL/AAEA/gAAAP4AAAD/AAAA/wAAAP8SExL/n62b/11lXP+XpJP/WF5W/wcIB/93k27/j7l//3upaP9WekX/DBEK/wAAAP4AAAD+AAAA/wAAAP8AAAD/CAgI/36Je/+1x7D/pbOh/yQlI/8DAwP/eo91/3CGaf9DWDz/Hyka/wECAf8AAAD+AAAA/gAAAP8AAAD/AAAA/wAAAP8NDg3/LTEs/xkbGf8AAAD/BwgH/5SlkP9BR0D/AAAA/wAAAP8AAAD/AAAA/gAAAP4AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wMDA/+Cjn//hZKB/0lQRv8JCQn/AAAA/wAAAP4AAAD8AAAA/gAAAP4AAAD+AAAA/gAAAP4AAAD+AAAA/gAAAP4AAAD+Ky8q/o6bi/5xfG3+CgoK/gAAAP4AAAD8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="

var shownAlerts = []string{}

func main() {
	base.Exit = exit
	systray.Run(onReady, onExit)
}

func StartupErrorHandler(message string) {
	base.Logf(logrus.ErrorLevel, message)
	//The PortAlreadyInUseException is logged twice, but we want to show the message only once
	for _, shownMessage := range shownAlerts {
		if shownMessage == message {
			return
		}
	}
	dialog.Alert(message)
	shownAlerts = append(shownAlerts, message)
}

func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("NZBHydra2")
	systray.SetTooltip("NZBHydra2")

	menuItemOpenWebUI := systray.AddMenuItem("Open web UI", "")
	go func() {
		<-menuItemOpenWebUI.ClickedCh
		if base.Uri == "" {
			base.Logf(logrus.ErrorLevel, "NZBHydra2 URI could not be determined")
			dialog.Alert("NZBHydra2 URI could not be determined")
			return
		}
		base.OpenBrowser(base.Uri)
	}()

	menuItemRestart := systray.AddMenuItem("Restart", "")
	go func() {
		<-menuItemRestart.ClickedCh
		if base.Uri == "" {
			base.Logf(logrus.ErrorLevel, "NZBHydra2 URI could not be determined")
			dialog.Alert("NZBHydra2 URI could not be determined")
			return
		}
		base.Logf(logrus.InfoLevel, "Sending restart command to main process")
		_, _ = base.ExecuteGetRequest(base.Uri + "internalapi/control/restart?internalApiKey=" + base.GetInternalApiKey())
	}()

	menuItemShutdown := systray.AddMenuItem("Shutdown", "")
	go func() {
		<-menuItemShutdown.ClickedCh
		if base.Uri != "" {
			base.Logf(logrus.InfoLevel, "Sending shutdown command to main process")
			resp, _ := base.ExecuteGetRequest(base.Uri + "internalapi/control/shutdown?internalApiKey=" + base.GetInternalApiKey())
			if resp.StatusCode != 200 {
				base.Logf(logrus.WarnLevel, "Shutdown command could not be sent to main process - shutting down wrapper")
				//Try shutting down wrapper, child process is hopefully killed gracefully
				exit(0)
			}
			base.Logf(logrus.InfoLevel, "Successfully sent shutdown command to main process")
			return
		}
		base.Logf(logrus.ErrorLevel, "NZBHydra2 URI could not be determined - shutting down wrapper")
		exit(0)

	}()

	base.Entrypoint(true, false, StartupErrorHandler)
}

func onExit() {
}
func exit(code int) {
	systray.Quit()
	//Give systray time to quit
	time.Sleep(1000 * time.Millisecond)
	os.Exit(code)
}

func getIcon() []byte {
	decodedBytes, err := base64.StdEncoding.DecodeString(IcoBytesB64)
	base.LogFatalIfError(err)
	return decodedBytes
}
