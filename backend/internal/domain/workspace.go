package domain

import "time"

type WorkspacePlan string

const (
	WorkspacePlanFree WorkspacePlan = "free"
	WorkspacePlanPro  WorkspacePlan = "pro"
)

type MemberRole string

const (
	MemberRoleEditor MemberRole = "editor"
	MemberRoleViewer MemberRole = "viewer"
	MemberRoleDev    MemberRole = "dev"
)

type Workspace struct {
	WorkspaceID           string
	Name                  string
	OwnerID               string
	Plan                  WorkspacePlan
	StripeCustomerID      string
	StripeSubscriptionID  string
	StorageUsedBytes      int64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type WorkspaceMember struct {
	WorkspaceID string
	UserID      string
	Role        MemberRole
	InvitedAt   time.Time
}

type PlanLimits struct {
	StorageQuotaBytes      int64
	MaxFileSizeBytes       int64
	MaxUploadsPerDay       int64
	MaxMembers             int64
	AllowedExtractionDepths []ExtractionDepth
}
