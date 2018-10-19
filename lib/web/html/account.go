package html

import (
	"net/http"

	"github.com/alexedwards/scs"
	"github.com/anyandrea/smartdev/lib/database/weatherdb"
	"github.com/anyandrea/smartdev/lib/web"
)

func Account(db weatherdb.WeatherDB, sm *scs.Manager) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		session := sm.Load(req)
		userId, _ := session.GetInt("user_id")
		if userId == 0 {
			http.Redirect(rw, req, "/login", http.StatusFound)
			return
		}

		user, err := db.GetUserById(userId)
		if err != nil {
			Error(sm, rw, req, err)
			return
		}

		// TODO: display account page for logged-in user
		page := &Page{
			Title:  "Newsfeed - Account",
			Active: "account",
			User:   user.Name,
		}
		web.Render().HTML(rw, http.StatusOK, "account", page)
	}
}
