package order

type Order struct {
	id    int
	items []Item
}

type Item struct {
	name     string
	quantity int
	value    int
}
