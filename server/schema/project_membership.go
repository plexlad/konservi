package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ProjectMembership struct {
	ent.Schema
}

func (ProjectMembership) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("project_id", uuid.UUID{}),
		field.Enum("role").
			Values("owner", "admin", "editor", "viewer"),
		field.Time("joined_at").
			Default(time.Now).
			Immutable(),
	}
}

func (ProjectMembership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("memberships").
			Field("user_id").
			Required().
			Unique(),
		edge.From("project", Project.Type).
			Ref("memberships").
			Field("project_id").
			Required().
			Unique(),
	}
}

func (ProjectMembership) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "project_id").
			Unique(),
	}
}
