package register

type IRegister interface {
	AddObject(object RegObject) error
	DelObject(objName string) error
	GetObject(objName string) (RegObject, error)
	GetObjects() ([]RegObject, error)
}

type IRegObject interface {
	AddToDataFilesState(file string) error
	RemoveInDataFilesState(file string) error
	GetDataFilesState(file string) ([]string, error)

	AddToStorageFilesState(file string) error
	RemoveInStorageFilesState(file string) error
	GetStorageFilesState(file string) ([]string, error)
}
