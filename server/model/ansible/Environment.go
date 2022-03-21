package db

import (
	"encoding/json"
	"errors"
)

// Environment is used to pass additional arguments, in json form to ansible
type Environment struct {
	ID        int     `db:"id" json:"id"`
	Name      string  `db:"name" json:"name" binding:"required"`
	ProjectID int     `db:"project_id" json:"project_id"`
	Password  *string `db:"password" json:"password"`
	JSON      string  `db:"json" json:"json" binding:"required"`
}

func (env *Environment) Validate() error {
	if env.Name == "" {
		return errors.New("environment name can not be empty")
	}

	if !json.Valid([]byte(env.JSON)) {
		return errors.New("environment must be valid JSON")
	}

	return nil
}
