package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"math/rand"
)

func main() {

	updatesChannel := make(chan Update)

	go func() {
		offset := 322245050
		access_token := "406250013:AAEBBjxkedB_tQi5JQzXmOV-vVg4xEDRSlg"
		for {
			time.Sleep(time.Second)
			request, err := http.NewRequest("GET", "https://api.telegram.org/bot"+access_token+"/"+"getUpdates?offset="+strconv.Itoa(offset), nil)
			if err != nil {
				fmt.Println(err)
				continue
			}

			http_client := &http.Client{}
			response, err := http_client.Do(request)
			if err != nil {
				fmt.Println(err)
				continue
			}

			buf := make([]byte, 32768)
			for {
				n, err := response.Body.Read(buf)

				if err == nil || fmt.Sprint(err) == "EOF" {
					if n == 0 {
						break
					}
					var updateResponse UpdateResponse
					json.Unmarshal(buf[:n], &updateResponse)

					if updateResponse.Ok {
						updates := updateResponse.Result
						for i := 0; i < len(updates); i++ {
							update := updates[i]
							offset = update.Update_id + 1
							updatesChannel <- update
						}
					}
				} else {
					fmt.Println(err)
					break;
				}
			}
		}
	}()


	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for {
		update := <-updatesChannel



		if r1.Intn(101) > 10 {
			continue
		}

		var text string
		phrase := r1.Intn(3)
		if phrase == 0{
			text = "Эх, Ангелиночка..."
		} else if phrase == 1{
			text = "Эх, Таечка..."
		} else {
			text = "ФУ, ДИАНА"
		}

		_, err := executeApiRequest("sendMessage", map[string]string{
			"chat_id":             strconv.Itoa(update.Message.Chat.Id),
			"text":                text,
			"reply_to_message_id": strconv.Itoa(update.Message.Message_id)})

		if err == nil {

		} else {
			fmt.Println(err)
		}
	}
}

func executeApiRequest(methodName string, params map[string]string) ([]byte, error) {
	access_token := "406250013:AAEBBjxkedB_tQi5JQzXmOV-vVg4xEDRSlg"

	paramsString := ""
	if params != nil {
		for k, v := range params {
			paramsString += k + "=" + v + "&"
		}
	}

	resp, err := http.Get("https://api.telegram.org/bot" + access_token + "/" + methodName + "?" + paramsString)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
