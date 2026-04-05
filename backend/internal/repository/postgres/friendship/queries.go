package friendship

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createFriendshipQuery() squirrel.InsertBuilder {
	return qb.Insert("friendships").
		Columns("requester_id", "addressee_id", "status")
}

func getFriendshipByIDQuery(id uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "requester_id", "addressee_id", "status", "created_at", "updated_at").
		From("friendships").
		Where(squirrel.Eq{"id": id})
}

func getFriendshipByUsersQuery(requesterID, addresseeID uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "requester_id", "addressee_id", "status", "created_at", "updated_at").
		From("friendships").
		Where(squirrel.Or{
			squirrel.And{squirrel.Eq{"requester_id": requesterID}, squirrel.Eq{"addressee_id": addresseeID}},
			squirrel.And{squirrel.Eq{"requester_id": addresseeID}, squirrel.Eq{"addressee_id": requesterID}},
		})
}

func getFriendsQuery(userID uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "requester_id", "addressee_id", "status", "created_at", "updated_at").
		From("friendships").
		Where(squirrel.And{
			squirrel.Eq{"status": "accepted"},
			squirrel.Or{
				squirrel.Eq{"requester_id": userID},
				squirrel.Eq{"addressee_id": userID},
			},
		})
}

func getPendingRequestsQuery(userID uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "requester_id", "addressee_id", "status", "created_at", "updated_at").
		From("friendships").
		Where(squirrel.Eq{"addressee_id": userID, "status": "pending"})
}

func updateFriendshipStatusQuery(id uuid.UUID, status string) squirrel.UpdateBuilder {
	return qb.Update("friendships").
		Set("status", status).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"id": id})
}

func deleteFriendshipQuery(id uuid.UUID) squirrel.DeleteBuilder {
	return qb.Delete("friendships").
		Where(squirrel.Eq{"id": id})
}
