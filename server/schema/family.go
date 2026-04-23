package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type Family struct {
	ent.Schema
}

func (Family) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("project_id", uuid.UUID{}),
		field.String("uid").
			Optional().
			Comment("GEDCOM UID - durable identifier across cycles"),
		field.String("external_id").
			Optional().
			Comment("GEDCOM EXID - external reference"),
		field.UUID("husband_id", uuid.UUID{}).
			Optional(),
		field.UUID("wife_id", uuid.UUID{}).
			Optional(),
		field.JSON("spouse_ids", []uuid.UUID{}).
			Optional().
			Comment("Additional spouses for polygamous families"),
		field.Time("marriage_date").
			Optional(),
		field.String("marriage_place").
			Optional(),
		field.String("marriage_calendar").
			Optional(),
		field.Time("marriage_banns_date").
			Optional().
			Comment("MARB - public notice of intent to marry"),
		field.Time("divorce_date").
			Optional(),
		field.Time("annulment_date").
			Optional(),
		field.String("marriage_type").
			Optional(),
		field.JSON("children_ids", []uuid.UUID{}).
			Optional(),
		field.Text("notes").
			Optional(),
		field.JSON("sources", []map[string]interface{}{}).
			Optional(),
		field.JSON("shared_notes", []map[string]interface{}{}).
			Optional(),
		field.Enum("privacy").
			Values("public", "confidential", "private").
			Default("public").
			Comment("RESN tag support for sensitive data"),
		field.UUID("created_by", uuid.UUID{}).
			Optional().
			Comment("UUID of User who created this record"),
		field.UUID("updated_by", uuid.UUID{}).
			Optional().
			Comment("UUID of User who last updated this record"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Family) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("families").
			Field("project_id").
			Required().
			Unique(),

		// FIX: use To edges with explicit Field() so ent owns the FK column
		// and Person's From edges can reference them with matching Ref() names.
		edge.To("husband", Person.Type).
			Field("husband_id").
			Unique(),
		edge.To("wife", Person.Type).
			Field("wife_id").
			Unique(),

		edge.To("spouses", Person.Type).
			Comment("Additional spouses linked via ASSO or multiple FAMs"),
		edge.To("children", Person.Type),
		edge.To("entries", Entry.Type),
	}
}

func (Family) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "husband_id", "wife_id"),
		// FIX: uid unique index will fail in most DBs when multiple NULLs exist.
		// Removed the unique constraint; enforce uniqueness in application logic
		// or use a partial index via a raw Atlas migration if your DB supports it.
		index.Fields("uid"),
		index.Fields("external_id"),
		index.Fields("marriage_date"),
		index.Fields("privacy"),
	}
}

