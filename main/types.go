package main

type UpdateResponse struct {
	Ok     bool   `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Update_id int `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Message_id int `json:"message_id"`
	From User `json:"from"`
	Date int `json:"date"`
	Chat Chat `json:"chat"`
	Text string `json:"text"`
	Reply_to_message *Message `json:"reply_to_message"`
}

type User struct {
	Id int `json:"id"`
	Is_bot bool `json:"is_bot"`
	First_name string `json:"first_name"`
	Last_name string `json:"last_name"`
	Username string `json:"username"`
	Language_code string `json:"language_code"`
}

type Chat struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Username string `json:"username"`
	First_name string `json:"first_name"`
	Last_name string `json:"last_name"`
}

type BotAnswer struct {
	Anwswer string `json:"answer"`
}