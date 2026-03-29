package domain

import "time"

type NodeLevel int

const (
	NodeLevelDomain  NodeLevel = 0
	NodeLevelConcept NodeLevel = 1
	NodeLevelMeasure NodeLevel = 2
	NodeLevelDetail  NodeLevel = 3
)

type NodeCategory string

const (
	NodeCategoryConcept  NodeCategory = "concept"
	NodeCategoryEntity   NodeCategory = "entity"
	NodeCategoryClaim    NodeCategory = "claim"
	NodeCategoryEvidence NodeCategory = "evidence"
	NodeCategoryCounter  NodeCategory = "counter"
)

type NodeEntityType string

const (
	NodeEntityTypeOrganization NodeEntityType = "organization"
	NodeEntityTypePerson       NodeEntityType = "person"
	NodeEntityTypeMetric       NodeEntityType = "metric"
	NodeEntityTypeDate         NodeEntityType = "date"
)

type Node struct {
	DocumentID    string
	NodeID        string
	Label         string
	Level         NodeLevel
	Category      NodeCategory
	EntityType    NodeEntityType // category == entity のみ
	Description   string
	SummaryHTML   string // iframe 向け HTML サマリ
	SourceChunkID string
	Confidence    float64
	CreatedAt     time.Time
}

type EdgeType string

const (
	EdgeTypeHierarchical EdgeType = "hierarchical"
	EdgeTypeSupports     EdgeType = "supports"
	EdgeTypeContradicts  EdgeType = "contradicts"
	EdgeTypeRelatedTo    EdgeType = "related_to"
	EdgeTypeMeasuredBy   EdgeType = "measured_by"
	EdgeTypeInvolves     EdgeType = "involves"
	EdgeTypeCauses       EdgeType = "causes"
	EdgeTypeExemplifies  EdgeType = "exemplifies"
)

type Edge struct {
	DocumentID    string
	EdgeID        string
	SourceNodeID  string
	TargetNodeID  string
	EdgeType      EdgeType
	Description   string
	Weight        float64
	SourceChunkID string
	CreatedAt     time.Time
}
