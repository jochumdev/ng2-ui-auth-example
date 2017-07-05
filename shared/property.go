package shared

import "errors"

const (
	PropertyTypeString = "string"
	PropertyTypeInt    = "int"
	PropertyTypeBool   = "bool"
)

var (
	ErrPropertyConversion = errors.New("Property conversion error")
)

type Property struct {
	/**
	 * The Name of the property.
	 */
	Key string

	/**
	 * Value store use AsXXX below to get the Value back in your type.
	 */
	Value interface{}

	/**
	 * The type of the property one of Property* Consts below
	 */
	Type string

	/**
	 * The properties View/Edit Permissions
	 */
	PermissionView string
	PermissionEdit string
}

func (p *Property) AsString() (string, error) {
	if p.Type != PropertyTypeString {
		return "", ErrPropertyConversion
	}

	return p.Value.(string), nil
}

func (p *Property) AsInt() (int, error) {
	if p.Type != PropertyTypeInt {
		return 0, ErrPropertyConversion
	}

	return p.Value.(int), nil
}

func (p *Property) AsBool() (bool, error) {
	if p.Type != PropertyTypeBool {
		return false, ErrPropertyConversion
	}

	return p.Value.(bool), nil
}
