package database

type User struct {
  Id       	int64
  Username 	string	`sql:"size:255; not null; unique;"`
  Password 	string 	`sql:"size:255; not null;"`
  Role 		string 	`sql:"size:30"` 
  IsEnabled	bool 	`sql:"default:true"`
}