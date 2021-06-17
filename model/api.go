package model

import (
	"github.com/wonderivan/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)



func CallWechatApi (method string) []byte {
	url := wxConfig.ApiUrl + method
	var resp *http.Response
	var body []byte
	var err error
	Loop1:
		for {
			resp, err = http.Get(url)
			if err != nil {
				logger.Info(err)
			} else {
				body, err = ioutil.ReadAll(resp.Body)
				if resp.StatusCode == 200 {
					break Loop1
				}
				logger.Info(err)
			}
			logger.Info(err)
			time.Sleep( 1 * time.Second)
		}
	return body
}

func CallWechatApiPost ( method string, msgText string ) {
	url := wxConfig.ApiUrl + method
	_, err := http.Post(url,"application/x-www-form-urlencoded",strings.NewReader(msgText))
	if err != nil {
		logger.Info(err)
	}
}

