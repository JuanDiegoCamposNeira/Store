package main

type Day struct {
	Type   string   `json:"type,omitempty"`
	Buyers []Person `json:"buyers,omitempty"`
}

type Person struct {
	Type string `json:"type,omitempty"`
	Uid  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type Product struct {
	Type  string `json:"type,omitempty"`
	Uid   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Price string `json:"price,omitempty"`
}

type Transaction struct {
	Type     string    `json:"type,omitempty"`
	Uid      string    `json:"uid,omitempty"`
	Buyer    Person    `json:"buyer,omitempty"`
	Ip       string    `json:"ip,omitempty"`
	Device   string    `json:"device,omitempty"`
	Products []Product `json:"products,omitempty"`
}

type Suggestion struct {
	Product  Product
	Quantity int
}
