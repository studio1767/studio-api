package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Action int

const (
	CREATE Action = iota
	READ
	UPDATE
	DELETE
)

func Authorize(ctx context.Context, object string, action Action) error {
	// extract the email and groups from the context
	email := EmailFromContext(ctx)
	if email == "" {
		return status.New(codes.PermissionDenied, "no subject provided").Err()
	}
	groups := GroupsFromContext(ctx)
	if len(groups) == 0 {
		return status.New(codes.PermissionDenied, "no group membership provided").Err()
	}

	// check for "admins"
	if groups["admins"] {
		if action == DELETE || action == CREATE || action == UPDATE || action == READ {
			return nil
		}
	}

	// check for "operators"
	if groups["operators"] {
		if action == CREATE || action == UPDATE || action == READ {
			return nil
		}
	}

	// check for "users"
	if groups["users"] {
		if action == READ {
			return nil
		}
	}

	// fallback
	return status.New(codes.PermissionDenied, "not authorized").Err()
}

func EmailFromContext(ctx context.Context) string {
	return ctx.Value(emailContextKey{}).(string)
}

func GroupsFromContext(ctx context.Context) map[string]bool {

	// get the groups embedded in the certificate
	groups := ctx.Value(groupsContextKey{}).(map[string]bool)

	// get the groups from the external source
	email := ctx.Value(emailContextKey{}).(string)
	getter := ctx.Value(getterContextKey{}).(GroupGetter)

	lgroups, err := getter.GroupsForUser(email)
	if err != nil {
		return groups
	}

	// merge the two maps
	for gname, _ := range lgroups {
		groups[gname] = true
	}

	return groups
}
