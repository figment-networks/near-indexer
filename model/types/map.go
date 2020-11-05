package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Map map[string]interface{}

func NewMap() Map {
	return Map{}
}

func (m Map) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Map) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("map requires []byte source")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*m, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("map type assertion failed")
	}

	return nil
}
