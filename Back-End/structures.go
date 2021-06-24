package main

type Person struct {
	Type string `json:"type,omitempty"`
	Uid  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
	Date string `json:"Date,omitempty"`
}

type Product struct {
	Type  string `json:"type,omitempty"`
	Uid   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Price string `json:"price,omitempty"`
	Date  string `json:"Date,omitempty"`
}

type Transaction struct {
	Type     string    `json:"type,omitempty"`
	Uid      string    `json:"uid,omitempty"`
	Buyer    Person    `json:"buyer,omitempty"`
	Ip       string    `json:"ip,omitempty"`
	Device   string    `json:"device,omitempty"`
	Products []Product `json:"products,omitempty"`
	Date     string    `json:"Date,omitempty"`
}
