package database

import (
	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/ozzadar/world_server/common"
)

var (
	db *gorm.DB
)


/* DATABASE FUNCTIONALITY */
// connect to the db
func InitDB() {

	common.Log("Initializing Database connection.")

	mysqlhost := "127.0.0.1:3306"
	mysqluser := "root"
	mysqlpass := "thesecretsauce"

    dbm, err := gorm.Open("mysql", mysqluser+ ":" + mysqlpass + 
    		"@(" + mysqlhost + ")/world_server?charset=utf8&parseTime=True")

    if(err != nil){
        panic("Unable to connect to the database")
    } else {
    	common.Log("Database connection established.")
    }

    db = &dbm
    dbm.DB().Ping()
    dbm.DB().SetMaxIdleConns(10)
    dbm.DB().SetMaxOpenConns(100)
    db.LogMode(true)
 
    if !dbm.HasTable(&User{}){
        dbm.CreateTable(&User{})
    }
}

//Register new user
func RegisterUser(newuser *User){

	common.Log(newuser)

	common.Log(db.NewRecord(newuser))
	db.Create(newuser)
	common.Log(db.NewRecord(newuser))

}

//Check credentials for login
func Login(username, password string) bool {
	newUser := &User{}

	db.Where("username = ?", username).First(&newUser)

	return newUser.Password == password
}