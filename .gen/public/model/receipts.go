//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Receipts struct {
	ID          int64 `sql:"primary_key"`
	UserID      int64
	Description string
	URLSlug     string
	Total       *int32
	CreatedAt   time.Time
}
