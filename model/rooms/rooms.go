package rooms

type Room struct {
	ROOMNO    string `db:"room_no"`
	ROOMTITLE string `db:"title"`
}
