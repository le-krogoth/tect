package model
/*-----------------------------------------------------------------------------
 **
 ** - TECT -
 **
 ** Copyright 2018 by Krogoth and the contributing authors
 **
 ** This program is free software; you can redistribute it and/or modify it
 ** under the terms of the GNU Affero General Public License as published by the
 ** Free Software Foundation, either version 3 of the License, or (at your option)
 ** any later version.
 **
 ** This program is distributed in the hope that it will be useful, but WITHOUT
 ** ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 ** FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 ** for more details.
 **
 ** You should have received a copy of the GNU Affero General Public License
 ** along with this program. If not, see <http://www.gnu.org/licenses/>.
 **
 ** -----------------------------------------------------------------------------*/

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Proteus struct {
	Id    string `gorm:"primary_key"`
	CrtDat	time.Time `sql:"DEFAULT:current_timestamp"`
	LastSeen	time.Time `sql:"DEFAULT:current_timestamp"`
}

func (m *Proteus) BeforeUpdate(scope *gorm.Scope) (err error) {

	scope.SetColumn("LastSeen", time.Now())
	return  nil
}

