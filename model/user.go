package model

type User struct {
	ID uint64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(user *User) error {
	statement := `INSERT INTO users(name, email, password) VALUES (?, ?, ?);`;

	_, err := db.Exec(statement, user.Name, user.Email, user.Password);
	return err;
}

func GetUser(id string) (User, error) {
	var user User

	statement := `SELECT u.id, u.name, u.email FROM users u WHERE u.id = ?;`;

	rows, err := db.Query(statement, id);
	if err != nil {
		return User{}, err;
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password);
		if err != nil {
			return User{}, err;
		}
	}

	return user, nil;
}

func CheckUser(email string, user *User) bool {
	statement := `SELECT u.id, u.name, u.email, u.password FROM users u WHERE u.email = ?`;
	rows, err := db.Query(statement, email);
	if err != nil { 
		return false ;
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password);
		if err != nil {
			return false;
		}
	}

	return true;
}