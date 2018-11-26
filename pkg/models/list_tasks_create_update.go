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
	"github.com/imdario/mergo"
)

// Create is the implementation to create a list task
// @Summary Create a task
// @Description Inserts a task into a list.
// @tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "List ID"
// @Param task body models.ListTask true "The task object"
// @Success 200 {object} models.ListTask "The created task object."
// @Failure 400 {object} models.HTTPError "Invalid task object provided."
// @Failure 403 {object} models.HTTPError "The user does not have access to the list"
// @Failure 500 {object} models.Message "Internal error"
// @Router /lists/{id} [put]
func (i *ListTask) Create(doer *User) (err error) {
	i.ID = 0

	// Check if we have at least a text
	if i.Text == "" {
		return ErrListTaskCannotBeEmpty{}
	}

	// Check if the list exists
	l := &List{ID: i.ListID}
	if err = l.GetSimpleByID(); err != nil {
		return
	}

	u, err := GetUserByID(doer.ID)
	if err != nil {
		return err
	}

	i.CreatedByID = u.ID
	i.CreatedBy = u
	_, err = x.Cols("text", "description", "done", "due_date_unix", "reminder_unix", "created_by_id", "list_id", "created", "updated").Insert(i)
	return err
}

// Update updates a list task
// @Summary Update a task
// @Description Updates a task. This includes marking it as done.
// @tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Param task body models.ListTask true "The task object"
// @Success 200 {object} models.ListTask "The updated task object."
// @Failure 400 {object} models.HTTPError "Invalid task object provided."
// @Failure 403 {object} models.HTTPError "The user does not have access to the task (aka its list)"
// @Failure 500 {object} models.Message "Internal error"
// @Router /tasks/{id} [post]
func (i *ListTask) Update() (err error) {
	// Check if the task exists
	ot, err := GetListTaskByID(i.ID)
	if err != nil {
		return
	}

	// For whatever reason, xorm dont detect if done is updated, so we need to update this every time by hand
	// Which is why we merge the actual task struct with the one we got from the
	// The user struct overrides values in the actual one.
	if err := mergo.Merge(&ot, i, mergo.WithOverride); err != nil {
		return err
	}

	// And because a false is considered to be a null value, we need to explicitly check that case here.
	if i.Done == false {
		ot.Done = false
	}

	_, err = x.ID(i.ID).Cols("text", "description", "done", "due_date_unix", "reminders_unix").Update(ot)
	*i = ot
	return
}