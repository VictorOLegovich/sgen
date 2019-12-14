package data

type Data struct {
	ID     int
	Family string
	Num
	Jet
}

type Num struct {
	ID     int
	DataID int
	Jet
}

type Jet struct {
	ID     int
	DataID int
	NumID  int
	Frazer int
}
