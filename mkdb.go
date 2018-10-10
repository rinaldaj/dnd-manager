package main

import (
	 "fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


func main(){
	db,err := sql.Open("mysql","root:@/")
	fmt.Println(err)
	query := "CREATE DATABASE dnd;"
	_,err = db.Query(query)
	fmt.Println(err)
	db.Close();
	db,err = sql.Open("mysql","root@/dnd")
	defer db.Close();
	fmt.Println(err)
	_,err = db.Query("CREATE TABLE player(health INT NOT NULL,maxHealth INT NOT NULL,strength INT NOT NULL,dexterity INT NOT NULL,Intelligence INT NOT NULL,wisdom INT NOT NULL,proficiencies TEXT NOT NULL,clothes NVARCHAR(200) NOT NULL,deathFails INT NOT NULL,alignment VARCHAR(30) NOT NULL,level INT NOT NULL,name NVARCHAR(200) NOT NULL, race NVARCHAR(200) NOT NULL,speed INT NOT NULL,charisma INT NOT NULL,constitution INT NOT NULL,PRIMARY KEY(name));")
	fmt.Println(err)
	_,err = db.Query("CREATE TABLE item(name VARCHAR(200) NOT NULL,weight DECIMAL(40,20) NOT NULL,value DECIMAL(40,20) NOT NULL,description TEXT NOT NULL,quantity INT NOT NULL,damage VARCHAR(10),dist INT,ammo VARCHAR(200),ac INT, modifier VARCHAR(20), owner NVARCHAR(200) NOT NULL,FOREIGN KEY(owner) REFERENCES player(name));")

}
