package members

type Member struct {
	ID       string `db:"user_id"`
	KORName  string `db:"kor_user_name"`
	ENGName  string `db:"eng_user_name"`
	PASSWD   string `db:"user_passwd"`
	OWNERIDX string `db:"owner_idx"`
}
