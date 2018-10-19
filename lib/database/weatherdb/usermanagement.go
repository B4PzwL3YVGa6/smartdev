package weatherdb

func (db *weatherDB) GetUsers() ([]User, error) {
	rows, err := db.Query(`
		select
			u.pk_user_id,
			u.password,
			u.name,
			u.email,
			u.role,
			u.active
		from user u
		order by u.email asc, u.name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		// get user
		var user User
		if err := rows.Scan(&user.Id, &user.Password, &user.Name, &user.Email, &user.Role, &user.Active); err != nil {
			return nil, err
		}

		// get users subscriptions
		var err error
		user.Subscriptions, err = db.GetSubscriptionsByUserId(user.Id)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (db *weatherDB) GetUserById(userId int) (User, error) {
	stmt, err := db.Prepare(`
		select
			u.pk_user_id,
			u.password,
			u.name,
			u.email,
			u.role,
			u.active
		from user u
		where u.pk_user_id = ?
		order by u.email asc, u.name asc`)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	var user User
	for rows.Next() {
		// get user
		if err := rows.Scan(&user.Id, &user.Password, &user.Name, &user.Email, &user.Role, &user.Active); err != nil {
			return User{}, err
		}

		// get users subscriptions
		var err error
		user.Subscriptions, err = db.GetSubscriptionsByUserId(user.Id)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}

func (db *weatherDB) GetUserByEmail(email string) (User, error) {
	stmt, err := db.Prepare(`
		select
			u.pk_user_id,
			u.password,
			u.name,
			u.email,
			u.role,
			u.active
		from user u
		where u.email = ?
		order by u.email asc, u.name asc`)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(email)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	var user User
	for rows.Next() {
		// get user
		if err := rows.Scan(&user.Id, &user.Password, &user.Name, &user.Email, &user.Role, &user.Active); err != nil {
			return User{}, err
		}

		// get users subscriptions
		var err error
		user.Subscriptions, err = db.GetSubscriptionsByUserId(user.Id)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}

func (db *weatherDB) GetSubscriptionsByUserId(userId int) ([]Subscription, error) {
	stmt, err := db.Prepare(`
		select
			s.fk_sensor_id,
			s.fk_user_id
		from subscription s
		where s.fk_user_id = ?
		order by s.fk_user_id asc, s.fk_sensor_id asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		// get subscription
		var subscription Subscription
		if err := rows.Scan(&subscription.SensorId, &subscription.UserId); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
