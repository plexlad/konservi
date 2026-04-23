package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type Entry struct {
	ent.Schema
}

func (Entry) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("project_id", uuid.UUID{}),
		field.UUID("author_id", uuid.UUID{}),
		field.Text("content").
			NotEmpty(),
		field.Time("entry_date").
			Default(time.Now),
		field.Time("created_at").
			Default(time.Now),
		field.Time("edited_last_at").
			Default(time.Now),
		field.JSON("linked_person_ids", []uuid.UUID{}).
			Optional(),
	}
}

func (Entry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("entries").
			Field("project_id").
			Required().
			Unique(),
		edge.From("author", User.Type).
			Ref("entries").
			Field("author_id").
			Required().
			Unique(),
		edge.To("linked_people", Person.Type),
	}
}

func (Entry) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "entry_date"),
		index.Fields("author_id"),
	}
}
