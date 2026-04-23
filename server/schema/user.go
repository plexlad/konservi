package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("username").
			Unique().
			NotEmpty(),
		field.String("email").
			Unique().
			NotEmpty(),
		field.String("password_hash"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("edited_last_at").
			Default(time.Now),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// Entry.author, Project.owner, and ProjectMembership.user are all
		// edge.From(...).Ref("X") — so User must supply the matching edge.To("X").
		edge.To("projects", Project.Type),
		edge.To("entries", Entry.Type),
		edge.To("memberships", ProjectMembership.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "email"),
	}
}

