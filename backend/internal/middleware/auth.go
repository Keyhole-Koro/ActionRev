package middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"

	"github.com/keyhole-koro/actionrev/internal/domain"
	"github.com/keyhole-koro/actionrev/internal/infra/firebase"
)

type contextKey string

const (
	contextKeyUserID      contextKey = "user_id"
	contextKeyMemberRole  contextKey = "member_role"
)

// AuthInterceptor は Connect RPC リクエストの Authorization ヘッダから
// Firebase ID Token を検証し、UID を context に注入する。
type AuthInterceptor struct {
	auth *firebase.Auth
}

func NewAuthInterceptor(auth *firebase.Auth) connect.UnaryInterceptorFunc {
	i := &AuthInterceptor{auth: auth}
	return connect.UnaryInterceptorFunc(i.WrapUnary)
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		token := extractBearerToken(req.Header().Get("Authorization"))
		if token == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}

		uid, err := i.auth.VerifyIDToken(ctx, token)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		ctx = context.WithValue(ctx, contextKeyUserID, uid)
		return next(ctx, req)
	}
}

func extractBearerToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}

// UserIDFromContext は context から Firebase Auth UID を取得する。
func UserIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(contextKeyUserID).(string)
	return uid, ok
}

// RequireRole はハンドラ内でロール確認を行うヘルパー。
// middleware での確認はルート単位では難しいため、handler 層で明示的に呼ぶ。
func RequireRole(role domain.MemberRole, actual domain.MemberRole) error {
	order := map[domain.MemberRole]int{
		domain.MemberRoleViewer: 1,
		domain.MemberRoleEditor: 2,
		domain.MemberRoleDev:    3,
	}
	if order[actual] < order[role] {
		return connect.NewError(connect.CodePermissionDenied, domain.ErrPermissionDenied)
	}
	return nil
}
