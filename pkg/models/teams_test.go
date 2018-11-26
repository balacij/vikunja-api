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

package models

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTeam_Create(t *testing.T) {
	//Dummyteam
	dummyteam := Team{
		Name:        "Testteam293",
		Description: "Lorem Ispum",
	}

	// Doer
	doer, err := GetUserByID(1)
	assert.NoError(t, err)

	// Insert it
	assert.True(t, dummyteam.CanCreate(&doer))
	err = dummyteam.Create(&doer)
	assert.NoError(t, err)

	// Check if it was inserted and we're admin
	tm := Team{ID: dummyteam.ID}
	err = tm.ReadOne()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tm.Members))
	assert.Equal(t, doer.ID, tm.Members[0].User.ID)
	assert.True(t, tm.Members[0].Admin)
	assert.True(t, dummyteam.CanRead(&doer))

	// Get all teams the user is part of
	ts, err := tm.ReadAll("", &doer, 1)
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(ts).Kind(), reflect.Slice)
	s := reflect.ValueOf(ts)
	assert.Equal(t, 2, s.Len())

	// Check inserting it with an empty name
	dummyteam.Name = ""
	err = dummyteam.Create(&doer)
	assert.Error(t, err)
	assert.True(t, IsErrTeamNameCannotBeEmpty(err))

	// update it (still no name, should fail)
	assert.True(t, dummyteam.CanUpdate(&doer))
	err = dummyteam.Update()
	assert.Error(t, err)
	assert.True(t, IsErrTeamNameCannotBeEmpty(err))

	// Update it, this time with a name
	dummyteam.Name = "Lorem"
	err = dummyteam.Update()
	assert.NoError(t, err)

	// Delete it
	assert.True(t, dummyteam.CanDelete(&doer))
	err = dummyteam.Delete()
	assert.NoError(t, err)

	// Try deleting a (now) nonexistant team
	err = dummyteam.Delete()
	assert.Error(t, err)
	assert.True(t, IsErrTeamDoesNotExist(err))

	// Try updating the (now) nonexistant team
	err = dummyteam.Update()
	assert.Error(t, err)
	assert.True(t, IsErrTeamDoesNotExist(err))
}

func TestIsErrInvalidTeamRight(t *testing.T) {
	assert.NoError(t, TeamRightAdmin.isValid())
	assert.NoError(t, TeamRightRead.isValid())
	assert.NoError(t, TeamRightWrite.isValid())

	// Check invalid
	var tr TeamRight
	tr = 938
	err := tr.isValid()
	assert.Error(t, err)
	assert.True(t, IsErrInvalidTeamRight(err))
}