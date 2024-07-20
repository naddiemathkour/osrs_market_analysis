package models

type Item struct {
	ID       interface{} `json:"id"`
	Name     interface{} `json:"name"`
	Members  interface{} `json:"members"`
	Highalch interface{} `json:"highalch"`
	Lowalch  interface{} `json:"lowalch"`
	Value    interface{} `json:"value"`
	Buylimit interface{} `json:"limit"`
	Icon     interface{} `json:"icon"`
	Examine  interface{} `json:"examine"`
}
