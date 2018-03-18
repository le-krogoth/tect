package db
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
	"log"
	"github.com/jinzhu/gorm"
	"area41.io/tect/lib"
	"area41.io/tect/db/model"
)

var db gorm.DB

//
func InitDatabase(activateLog bool) {

	dialect := lib.GetStringFromConfig("db.dialect")
	args := lib.GetStringFromConfig("db.args")

	database, err := gorm.Open(dialect, args)
	if err != nil {
		log.Fatalf("failed to connect database, %s", err)
		panic("failed to connect database")
	}

	db = *database

	db.SingularTable(true)

	if activateLog {

		db.LogMode(true)
	}

	db.AutoMigrate(&model.Proteus{})
	db.AutoMigrate(&model.Log{})
}

//
func CloseDB() {

	db.Close()
}

//
func GetDB() *gorm.DB {

	return &db
}
