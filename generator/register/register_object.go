package register

type RegObject struct {
	Name              string
	DataFilesState    []string
	StorageFilesState []string
}

func NewRegObject() *RegObject {
	return &RegObject{}
}

func (obj *RegObject) AddToDataFilesState(file string) error {
	panic("implement me")
}

func (obj *RegObject) RemoveInDataFilesState(file string) error {
	panic("implement me")
}

func (obj *RegObject) GetDataFilesState(file string) ([]string, error) {
	panic("implement me")
}

func (obj *RegObject) AddToStorageFilesState(file string) error {
	panic("implement me")
}

func (obj *RegObject) RemoveInStorageFilesState(file string) error {
	panic("implement me")
}

func (obj *RegObject) GetStorageFilesState(file string) ([]string, error) {
	panic("implement me")
}
