package templates

type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var Users []Data = []Data{
	{Name: "Maxim", Age: 22},
	{Name: "Anatoly", Age: 82},
	{Name: "Ida", Age: 79},
}
