package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type Person struct {
	ent.Schema
}

func (Person) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("project_id", uuid.UUID{}),
		field.String("external_id").
			Optional().
			Comment("External identifier (EXID tag)"),
		field.String("primary_name").
			NotEmpty(),
		field.JSON("names", []map[string]interface{}{}).
			Optional().
			Comment("All name variants with LANG and TRAN"),
		field.String("preferred_name").
			Optional().
			Comment("Rufname/call name"),
		field.String("maiden_name").
			Optional(),
		field.Enum("sex").
			Values("M", "F", "U", "O", "X"),
		field.Time("birth_date").
			Optional(),
		field.String("birth_date_approx").
			Optional().
			Comment("For BET/AND dates"),
		field.String("birth_place").
			Optional(),
		field.JSON("birth_place_jurisdiction", []string{}).
			Optional(),
		field.Time("death_date").
			Optional(),
		field.String("death_place").
			Optional(),
		field.String("death_cause").
			Optional(),
		field.JSON("death_details", []map[string]interface{}{}).
			Optional(),
		field.String("burial_place").
			Optional(),
		field.Time("burial_date").
			Optional(),
		field.String("occupation").
			Optional(),
		field.JSON("occupation_details", []map[string]interface{}{}).
			Optional(),
		field.String("education").
			Optional(),
		field.String("religion").
			Optional(),
		field.String("nationality").
			Optional(),
		field.String("ethnicity").
			Optional(),
		field.JSON("aliases", []string{}).
			Optional(),
		field.JSON("nicknames", []string{}).
			Optional(),
		field.JSON("titles", []string{}).
			Optional(),
		field.JSON("height", map[string]interface{}{}).
			Optional(),
		field.JSON("weight", map[string]interface{}{}).
			Optional(),
		field.JSON("blood_type", map[string]interface{}{}).
			Optional(),
		field.JSON("medical_conditions", []map[string]interface{}{}).
			Optional(),
		field.JSON("genetic_markers", []map[string]interface{}{}).
			Optional(),
		field.JSON("dna_matches", []map[string]interface{}{}).
			Optional(),
		field.JSON("photos", []map[string]interface{}{}).
			Optional(),
		field.JSON("documents", []map[string]interface{}{}).
			Optional(),
		field.JSON("events", []map[string]interface{}{}).
			Optional(),
		field.Text("notes").
			Optional(),
		field.JSON("sources", []map[string]interface{}{}).
			Optional(),
		field.JSON("shared_notes", []map[string]interface{}{}).
			Optional(),
		field.JSON("media", []map[string]interface{}{}).
			Optional(),
		field.Enum("privacy").
			Values("public", "confidential", "private").
			Default("public").
			Comment("RESN tag support"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		// FIX: marked Optional so inserts don't panic when not provided
		field.UUID("father_id", uuid.UUID{}).
			Optional(),
		field.UUID("mother_id", uuid.UUID{}).
			Optional(),
		field.UUID("created_by", uuid.UUID{}).
			Optional().
			Comment("UUID of User who created this record"),
		field.UUID("updated_by", uuid.UUID{}).
			Optional().
			Comment("UUID of User who last updated this record"),
	}
}

func (Person) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("people").
			Field("project_id").
			Required().
			Unique(),

		// FIX: declare father/mother edges with explicit Field() so ent
		// generates father_id / mother_id columns (required for the indexes below)
		edge.To("father", Person.Type).
			Field("father_id").
			Unique(),
		edge.To("mother", Person.Type).
			Field("mother_id").
			Unique(),
		edge.To("children", Person.Type),

		// Spouse relationship (M2M self-referential)
		edge.To("spouses", Person.Type),

		// Entries authored by this person
		edge.To("entries", Entry.Type),

		// Entries where this person is tagged/linked
		// FIX: Ref must match the To edge name on Entry ("linked_people")
		edge.From("linked_entries", Entry.Type).
			Ref("linked_people"),

		// FIX: self-referential M2M edges must be declared only as To;
		// ent handles the inverse automatically for symmetric relations.
		// Using separate named edges avoids the ambiguous From+To pattern.
		edge.To("relatives", Person.Type),
		edge.To("associations", Person.Type),

		// Family edges — inverse refs must match To edges declared on Family
		edge.From("families_as_husband", Family.Type).
			Ref("husband"),
		edge.From("families_as_wife", Family.Type).
			Ref("wife"),
		edge.From("families_as_child", Family.Type).
			Ref("children"),
	}
}

func (Person) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "primary_name"),
		index.Fields("id").Unique(),
		index.Fields("external_id"),
		index.Fields("birth_date"),
		index.Fields("death_date"),
		// FIX: father_id / mother_id now exist as edge-backed fields
		index.Fields("father_id", "mother_id"),
		index.Fields("privacy"),
	}
}

