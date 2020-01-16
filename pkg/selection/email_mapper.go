package selection

import (
	"context"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// GenerateUserEmailMap creates a map of user emails given a member list
// to make searching for members with email an O(1) call.
func GenerateUserEmailMap(ctx context.Context, members *rm.MemberList, iamc iam.IAMServiceClient) (map[string]string, error) {
	res := make(map[string]string)
	for _, m := range members.Items {
		user, err := iamc.GetUser(ctx, &common.IDOptions{Id: m.UserId})
		if err != nil {
			return nil, err
		}
		res[user.Email] = user.Id
	}
	return res, nil
}
