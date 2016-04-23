package database

import (
	"github.com/pagarme/teleport/action"
	"github.com/pagarme/teleport/batcher/ddldiff"
)

// Define a class attribute
type Attribute struct {
	Name         string `json:"attr_name"`
	Num          int    `json:"attr_num"`
	TypeName     string `json:"type_name"`
	TypeSchema   string `json:"type_schema"`
	TypeOid      string `json:"type_oid"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	Class        *Class
}

// Implements Diffable
func (post *Attribute) Diff(other ddldiff.Diffable, context ddldiff.Context) []action.Action {
	actions := make([]action.Action, 0)

	// r = Tables
	if post.Class.RelationKind == "r" {
		if other == nil {
			actions = append(actions, &action.CreateColumn{
				context.Schema,
				post.Class.RelationName,
				action.Column{
					post.Name,
					post.TypeName,
					post.TypeSchema,
				},
			})
		} else {
			pre := other.(*Attribute)

			if pre.Name != post.Name || pre.TypeOid != post.TypeOid {
				actions = append(actions, &action.AlterColumn{
					context.Schema,
					post.Class.RelationName,
					action.Column{
						pre.Name,
						pre.TypeName,
						pre.TypeSchema,
					},
					action.Column{
						post.Name,
						post.TypeName,
						post.TypeSchema,
					},
				})
			}
		}
	}

	return actions
}

func (a *Attribute) Children() []ddldiff.Diffable {
	return []ddldiff.Diffable{}
}

func (a *Attribute) Drop(context ddldiff.Context) []action.Action {
	return []action.Action{
		&action.DropColumn{
			context.Schema,
			a.Class.RelationName,
			action.Column{
				a.Name,
				a.TypeName,
				a.TypeSchema,
			},
		},
	}
}

func (a *Attribute) IsEqual(other ddldiff.Diffable) bool {
	if other == nil {
		return false
	}

	otherAttr := other.(*Attribute)
	return (a.Num == otherAttr.Num)
}
