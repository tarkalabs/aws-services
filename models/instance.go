package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

//Instance represents an AWS EC2 instance
type Instance struct {
	ID         string `db:"id"`
	Name       string `db:"name"`
	Region     string `db:"region"`
	Attributes string `db:"attributes"`
}

// NewInstance returns a new Instance struct
func NewInstance() *Instance {
	return &Instance{ID: uuid.NewV4().String()}
}

// SaveInstance saves a given instance
func SaveInstance(db *sqlx.DB, instance *Instance) {
	_, err := db.Exec("INSERT into instances(id, name, region, attributes) values ($1, $2, $3, $4)", instance.ID, instance.Name, instance.Region, instance.Attributes)
	failOnError(err, "unable to insert record")
}
