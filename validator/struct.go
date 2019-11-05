package validator

import (
	"errors"
	c "github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/types"
)

func StructsValidation(Structs []c.Struct) error {
	var errorText string
	var requiredFields []string

	errs := map[string]string{}

	structsNames := ExtractStructsNames(Structs)

	for _, Struct := range Structs {
		if len(Struct.Complicated) != 0 {
			for _, comp := range Struct.Complicated {
				errs[fieldsSection] = "comp: " + comp.String()
			}
		}

		requiredFields = StructValidation(Struct)

		if len(requiredFields) != 0 {
			for _, req := range requiredFields {
				if !StructExist(structsNames, req) {
					errs[fieldsSection] = "reqError -> " + req
				}
			}
		}
	}

	if len(errs) != 0 {
		for section, s := range errs {
			errorText += section + ":\n" + s
		}
		return errors.New(errorText)
	}

	return nil
}

func ExtractStructsNames(Structs []c.Struct) (Names []string) {
	for _, Struct := range Structs {
		Names = append(Names, Struct.Name)
	}
	return Names
}

func StructExist(StructsNames []string, StructName string) bool {
	for _, SName := range StructsNames {
		if SName == StructName {
			return true
		}
	}
	return false
}

func StructValidation(Struct c.Struct) (required []string) {
	for _, Field := range Struct.Fields {
		if req := FieldValidation(Field); req != "" {
			required = append(required, req)
		}
	}

	return required
}

func FieldValidation(Field c.Field) string {
	return typeVerify(Field)
}

func typeVerify(Field c.Field) string {
	ftype := Field.Type

	if result, typename := types.IsArray(ftype); result {
		if !types.IsSimpleType(typename) {
			return typename
		}
		return ""
	}

	if result, maptype := types.IsMap(ftype); result {
		if !types.IsSimpleType(maptype.Key) {
			return maptype.Key
		}

		if !types.IsSimpleType(maptype.Value) {
			return maptype.Value
		}
		return ""
	}

	if !types.IsSimpleType(ftype) {
		return ftype
	}

	return ""
}
