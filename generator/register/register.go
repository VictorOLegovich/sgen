package register

type Register struct {
	Objects []RegObject
}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) AddObject(object RegObject) error {
	panic("implement me")
}

func (r *Register) DelObject(objName string) error {
	panic("implement me")
}

func (r *Register) GetObject(objName string) (RegObject, error) {
	panic("implement me")
}

func (r *Register) GetObjects() ([]RegObject, error) {
	panic("implement me")
}
