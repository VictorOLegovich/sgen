package register

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const format = "reg"

type Register struct {
	objects []RegObject
	changed string
	inited  bool
	Deleted map[string]string
}

type RegObject struct {
	Package  string
	Entistor map[string]string //the connection between the entity and the storage
}

func NewRegister() (*Register, error) {
	objects, e := getObjects()
	if e != nil {
		println(e.Error())
		return &Register{}, errors.New("Не удалось получить объекты ")
	}

	return &Register{objects, "", true, map[string]string{}}, nil
}

func (r *Register) AddObject(object RegObject) {
	obj, err := r.GetObject(object.Package)
	r.changed = object.Package

	//Значит, этого объекта не существует в регистре и мы можем добавить его без проверок.
	if err != nil {
		r.objects = append(r.objects, object)
		return
	}
	//Записываем удалённые
	r.Deleted = r.hasDeletedStructs(obj, object)
	//Перезаписываем стейт объектов
	_ = r.delObject(object.Package)
	r.objects = append(r.objects, object)
}

func (r *Register) delObject(objName string) error {
	remove := func(s []RegObject, i int) []RegObject {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}

	for key, object := range r.objects {
		if object.Package == objName {
			r.objects = remove(r.objects, key)
			return nil
		}
	}

	return errors.New("Not exist ")
}

func (r *Register) GetObject(objName string) (RegObject, error) {
	for _, obj := range r.objects {
		if obj.Package == objName {
			return obj, nil
		}
	}
	return RegObject{}, errors.New("Object with this package name  (" + objName + ") is not exist! ")
}

func (r *Register) Save() error {
	if !r.inited {
		return errors.New("Регистр не был инициализирован! ")
	}
	var filePath, src string

	if object, err := r.GetObject(r.changed); err == nil {
		filePath, _ = filepath.Abs("../register/register/" + object.Package + "." + format)
		for s, path := range object.Entistor {
			src += s + ":" + path + "\n"
		}
		if err = write(filePath, src); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (r *Register) hasDeletedStructs(existingObject, addedObject RegObject) map[string]string {
	for structure := range addedObject.Entistor {
		for s := range existingObject.Entistor {
			if structure == s {
				delete(existingObject.Entistor, s)
			}
		}
	}
	return existingObject.Entistor
}

func getObjects() (objects []RegObject, e error) {
	var object RegObject

	fp, _ := filepath.Abs("../register/register")
	files, e := ioutil.ReadDir(fp)
	if e != nil {
		return objects, e
	}

	for _, file := range files {
		ff := strings.Split(file.Name(), ".")
		if len(ff) > 0 {
			if ff[1] == format {
				object.Package = ff[0]
				object.Entistor = read(filepath.Join(fp, file.Name()))
				objects = append(objects, object)
			}
		}
	}
	return objects, e
}

func read(filePath string) map[string]string {
	var (
		entistor        = map[string]string{}
		lines, elements []string
		src             []byte
	)

	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		println(err.Error())
	}

	lines = strings.Split(string(src), "\n")

	for _, line := range lines {
		elements = strings.Split(line, ":")
		if len(elements) > 1 {
			entistor[elements[0]] = elements[1]
		}
	}
	return entistor
}

func write(filePath, src string) error {
	file, _ := os.Create(filePath)
	_, _ = file.WriteString(src)
	_ = file.Close()

	return nil
}
