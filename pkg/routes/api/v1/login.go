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
	"code.vikunja.io/api/pkg/routes/crud"
	"crypto/md5"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

// Token represents an authentification token
type Token struct {
	Token string `json:"token"`
}

// Login is the login handler
// @Summary Login
// @Description Logs a user in. Returns a JWT-Token to authenticate further requests.
// @tags user
// @Accept json
// @Produce json
// @Param credentials body models.UserLogin true "The login credentials"
// @Success 200 {object} v1.Token
// @Failure 400 {object} models.Message "Invalid user password model."
// @Failure 403 {object} models.Message "Invalid username or password."
// @Router /login [post]
func Login(c echo.Context) error {
	u := models.UserLogin{}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, models.Message{"Please provide a username and password."})
	}

	// Check user
	user, err := models.CheckUserCredentials(&u)
	if err != nil {
		return crud.HandleHTTPError(err)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	avatar := md5.Sum([]byte(user.Email))
	claims["avatar"] = hex.EncodeToString(avatar[:])

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("service.JWTSecret")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Token{Token: t})
}