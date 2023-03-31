package collection

type Collecter interface {
	Collect() error
	Id() string
}
