package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"math/rand"
	"net/url"
	"strings"
)

func main() {

	updatesChannel := make(chan Update)

	go func() {

		dat, err := ioutil.ReadFile("offset")
		check(err)
		offset, err := strconv.Atoi(fmt.Sprint(string(dat)))
		check(err)

		access_token := "406250013:AAEBBjxkedB_tQi5JQzXmOV-vVg4xEDRSlg"
		for {
			time.Sleep(time.Second)
			request, err := http.NewRequest("GET", "https://api.telegram.org/bot"+access_token+"/"+"getUpdates?offset="+strconv.Itoa(offset)+"&timeout=10", nil)
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

			var updateResponse UpdateResponse

			result := make([]byte, 16777216)
			cur := 0;
			buf := make([]byte, 32768)

			n, err := response.Body.Read(buf)
			for ; !(err != nil && n == 0); {

				for i := cur; i < cur+n; i++ {
					result[i] = buf[i-cur]
				}
				cur += n;
				n, err = response.Body.Read(buf)
			}

			json.Unmarshal(result[:cur], &updateResponse)
			if updateResponse.Ok {
				updates := updateResponse.Result
				for i := 0; i < len(updates); i++ {
					update := updates[i]
					offset = update.Update_id + 1

					d1 := []byte(strconv.Itoa(offset))
					err := ioutil.WriteFile("offset", d1, 0644)
					check(err)

					updatesChannel <- update
				}
			}

			defer response.Body.Close()
		}
	}()

	ml := make(map[string][]string)




	for {
		update := <-updatesChannel

		if update.Message.Reply_to_message != nil {
			tokens := strings.Split(strings.ToLower(update.Message.Reply_to_message.Text), " ")
			for i := 0; i < len(tokens); i++ {
				if ml[tokens[i]] == nil {
					ml[tokens[i]] = []string{update.Message.Text}
				} else {
					ml[tokens[i]] = append(ml[tokens[i]], update.Message.Text)
				}
			}
		}

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		if r1.Intn(101) > 10 {
			continue
		}


		text := "+"

		tokens := strings.Split(strings.ToLower(update.Message.Text), " ")
		token := tokens[0]
		for i := 0; i < len(tokens) ; i++ {
			if (ml[tokens[i]] != nil) {
				token = tokens[i]
			}
		}
		if (ml[token] != nil) {
			s2 := rand.NewSource(time.Now().UnixNano())
			r2 := rand.New(s2)
			text = ml[token][r2.Intn(len(ml[token]))]
		}

		_, err := executeApiRequest("sendMessage", url.Values{
			"chat_id":             {strconv.Itoa(update.Message.Chat.Id)},
			"text":                {text},
			"reply_to_message_id": {strconv.Itoa(update.Message.Message_id)}})

		if err == nil {
			//fmt.Printf("%s", resp)
		} else {
			fmt.Println(err)
		}
	}

}

func executeApiRequest(methodName string, params url.Values) ([]byte, error) {
	access_token := "406250013:AAEBBjxkedB_tQi5JQzXmOV-vVg4xEDRSlg"

	resp, err := http.Get("https://api.telegram.org/bot" + access_token + "/" + methodName + "?" + params.Encode())
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
