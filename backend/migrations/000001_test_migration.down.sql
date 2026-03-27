-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd