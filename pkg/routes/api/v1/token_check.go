//  Vikunja is a todo-list application to facilitate your life.
//  Copyright 2018 Vikunja and contributors. All rights reserved.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <https://www.gnu.org/licenses/>.

package v1

import (
	"code.vikunja.io/api/pkg/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// CheckToken checks prints a message if the token is valid or not. Currently only used for testing pourposes.
func CheckToken(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)

	fmt.Println(user.Valid)

	return c.JSON(418, models.Message{"🍵"})
}