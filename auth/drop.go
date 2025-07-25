package auth

// Copyright (C) Ava Glass <SuperNinja_4965>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

import (
	"log"
	"net/http"
)

// dropHandler immediately closes the connection without response
func DropHandler(w http.ResponseWriter, r *http.Request) {
	// Get the underlying connection and close it immediately
	if hj, ok := w.(http.Hijacker); ok {
		conn, _, err := hj.Hijack()
		if err == nil {
			if cerr := conn.Close(); cerr != nil {
				log.Printf("Error closing connection: %v", cerr)
			}
			return
		}
	}
	// Fallback if hijacking fails - return 404
	http.NotFound(w, r)
}
