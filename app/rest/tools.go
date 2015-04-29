// General tools for rest functions
package rest

import (
	"fmt"
	"net/http"
	"appengine"
	"appengine/user"
)

// Returns current username, or empty string if none
func getUserEmail(r *http.Request) (string, error) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		return "", fmt.Errorf("No user found.")
	}
	return u.Email, nil
}

// Get a login URL
func getLoginURL(r *http.Request) (string, error) {
	c := appengine.NewContext(r)
	// TODO: better redirect path
	dest := r.URL.String()
	return user.LoginURL(c, dest)
}

// Check if current user is an admin
func isAdmin(r *http.Request) bool {
	c := appengine.NewContext(r)
	return user.IsAdmin(c)
}

// Wrapper function to require a logged-in user.
func authRequired(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := getUserEmail(r); err != nil {
			loginURL, err := getLoginURL(r)
			if err != nil {
				http.Error(w, "Access Denied", 403)
				return
			}
			http.Redirect(w, r, loginURL, http.StatusFound)
			return
		}
		f(w, r)
	}
}

// Wrapper function for admin required.
func adminRequired(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdmin(r) {
			http.Error(w, "Access Denied", 403)
			return
		}
		f(w, r)
	}
}
