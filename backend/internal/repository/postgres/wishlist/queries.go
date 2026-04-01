package wishlist

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createWishlistQuery() squirrel.InsertBuilder {
	return qb.Insert("wishlists").
		Columns("id", "participant_id", "visibility", "created_at", "updated_at")
}

func createWishlistItemQuery() squirrel.InsertBuilder {
	return qb.Insert("wishlist_items").
		Columns("id", "wishlist_id", "title", "link", "image_url", "comment", "created_at")
}

func getWishlistByParticipantQuery(participantID string) squirrel.SelectBuilder {
	return qb.Select("id", "participant_id", "visibility", "created_at", "updated_at").
		From("wishlists").
		Where(squirrel.Eq{"participant_id": participantID})
}

func getWishlistItemsQuery(wishlistID string) squirrel.SelectBuilder {
	return qb.Select("id", "wishlist_id", "title", "link", "image_url", "comment", "created_at").
		From("wishlist_items").
		Where(squirrel.Eq{"wishlist_id": wishlistID})
}
func updateWishlistItemQuery(itemID string) squirrel.UpdateBuilder {
	return qb.Update("wishlist_items").
		Where(squirrel.Eq{"id": itemID})
}

func deleteWishlistItemQuery(itemID string) squirrel.DeleteBuilder {
	return qb.Delete("wishlist_items").
		Where(squirrel.Eq{"id": itemID})
}
