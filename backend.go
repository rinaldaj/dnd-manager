package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Object interface{
	getName() string
	getWeight() float32
	getValue()	float32
	getDescription()	string
}


type Item struct {
	//WHat the item is called
	Name	string
	//Weight is weight in lbs
	Weight	float32
	//value is in gold coins
	Value	float32
	//The description of an item
	Description	string
}
func (i Item) getName() string{
	return i.Name
}
func (i Item) getDescription() string{
	return i.Description
}
func (i Item) getValue() float32{
	return i.Value
}
func (i Item) getWeight() float32{
	return i.Weight
}

type Weapon struct{
	Base	Item
	Damage	int //Specifies number of sides the die has a spell with a save against it has this value
	Range	int //in feet
	Ammo	string //The type of ammunition this is tracked seperatly in the case of spells this is a spell slot
	Mod	string //This is the trait that is used as the modifier
}
func (i Weapon) getName() string{
	return i.Base.Name
}
func (i Weapon) getDescription() string{
	return i.Base.Description
}
func (i Weapon) getValue() float32{
	return i.Base.Value
}
func (i Weapon) getWeight() float32{
	return i.Base.Weight
}

type Armor struct{
	Base	Item
	AC	int
	Mod	string
}

func (i Armor) getName() string{
	return i.Base.Name
}
func (i Armor) getDescription() string{
	return i.Base.Description
}
func (i Armor) getValue() float32{
	return i.Base.Value
}
func (i Armor) getWeight() float32{
	return i.Base.Weight
}

type Player struct {
	Inventory	[]Object
	Health	int
	MaxHealth	int
	Strength	int
	Dexterity	int
	Intelligence	int
	Wisdom	int
	Charisma	int
	Proficienies	[]string//THis is a list of strings specifying proficiency
	Clothes	Armor //THis is the item that is currently equiped
	DeathFails int //If this is ever over three player is dead, if it is set to a negative value then the player is stable
	Alignment	string
	Class	string //Currently doesn't do anything
	Race	string //Currently useless
	Level	int //The level of the character (does not support multiclassing)
	Name	string
}

var dbPass string

func getInventory(player string, db *sql.DB) []Object{
	var ret []Object
	query := fmt.Sprintf("SELECT name,weight,value,description,quantity,damage,dist,ammo,ac,modifier FROM item where owner=%q",player)
	results,err := db.Query(query)
	if err != nil {
		return ret
	}
	for results.Next(){
		var nuItem Item
		var dam *int
		var dist *int
		var ammo *string
		var ac *int
		var mod *string
		var quant int
		if err := results.Scan(&nuItem.Name,&nuItem.Weight,&nuItem.Value,&nuItem.Description,&quant,&dam,&dist,&ammo,&ac,&mod); err != nil{
			continue
		}
		if dist != nil {
			var weaponBox Weapon;
			weaponBox = Weapon{nuItem,*dam,*dist,*ammo,*mod}
			for i := 0;i<quant;i++{
			ret = append(ret,weaponBox)
			}
			continue
		}
		if ac != nil {
			var armorBox Armor;
			armorBox = Armor{nuItem,*ac,*mod}
			for i := 0;i<quant;i++{
			ret = append(ret,armorBox)
			}
			continue
		}
		for i := 0;i<quant;i++{
			ret = append(ret,nuItem)
		}
	}
	return ret
}


func getPlayer(player string,db *sql.DB) Player{
	query := fmt.Sprintf("SELECT health,maxHealth,strength,dexterity,Intelligence,wisdom,proficiencies,clothes,deathFails,alignment,level,name FROM player where name=%q;",player)
	res,err := db.Query(query)
	var nuPlayer Player;
	var prof *string;
	var cloth *string;
	if err != nil {
		return Player{}
	}
	res.Next()
	if err = res.Scan(&nuPlayer.Health,&nuPlayer.MaxHealth,&nuPlayer.Strength,&nuPlayer.Dexterity,&nuPlayer.Intelligence,&nuPlayer.Wisdom,&prof,&cloth,&nuPlayer.DeathFails,&nuPlayer.Alignment,&nuPlayer.Level,&nuPlayer.Name); err != nil {
		return Player{}
}
	inv := getInventory(player,db)
	nuPlayer.Inventory = inv
	for _,i := range inv {
		switch v := i.(type){
			case Armor:
				if cloth != nil && i.getName() == *cloth{
					nuPlayer.Clothes = v
					break
				}
		}
	}
	if prof != nil{
		nuPlayer.Proficienies = strings.Split(*prof," ")
	}
	return nuPlayer
}


func routeSelectHandler(w http.ResponseWriter, r *http.Request){
	DB,err := sql.Open("mysql",dbPass)
	if err != nil {
		fmt.Fprintf(w,"ERROR: COULD NOT TOUCH DB")
		return
	}
	defer DB.Close()
	cur :=getPlayer(r.FormValue("name"),DB)
	if cur.Name == ""{
		http.Redirect(w,r,"/makecharacter.html",http.StatusSeeOther)
	} else {
		http.Redirect(w,r,"/viewCharacter",http.StatusSeeOther)
	}
}

func charHandler(w http.ResponseWriter, r *http.Request){
	DB,err := sql.Open("mysql",dbPass)
	if err != nil {
		fmt.Fprintf(w,"ERROR: COULD NOT TOUCH DB")
		return
	}
	plas := Player{}
	fmt.Sscanf(r.FormValue("HP"),"%d",&plas.Health)
	fmt.Sscanf(r.FormValue("HP"),"%d",&plas.MaxHealth)
	fmt.Sscanf(r.FormValue("Strength"),"%d",&plas.Strength)
	fmt.Sscanf(r.FormValue("Dexterity"),"%d",&plas.Dexterity)
	fmt.Sscanf(r.FormValue("Intelligence"),"%d",&plas.Intelligence)
	fmt.Sscanf(r.FormValue("Wisdom"),"%d",&plas.Wisdom)
	fmt.Sscanf(r.FormValue("Charisma"),"%d",&plas.Charisma)
	fmt.Sscanf(r.FormValue("Level"),"%d",&plas.Level)
	plas.DeathFails = -1
	plas.Alignment = r.FormValue("Alignment")
	plas.Class = r.FormValue("Class")
	plas.Name = r.FormValue("name")
//Player{empt,r.FormValue("HP"),r.FormValue("HP"),r.FormValue("Strength"),r.FormValue("Dexterity"),r.FormValue("Intelligence"),r.FormValue("Wisdom"),r.FormValue("Charisma"),empts,Armor{},-1,r.FormValue("Alignment"),r.FormValue("Class"),r.FormValue("race"),r.FormValue("Level"),r.FormValue("name")}
	query := fmt.Sprintf("INSERT INTO player(health,maxHealth,strength,dexterity,Intelligence,wisdom,deathFails,alignment,level,name) VALUES(%d,%d,%d,%d,%d,%d,%d,%q,%d,%q);",plas.Health,plas.MaxHealth,plas.Strength,plas.Dexterity,plas.Intelligence,plas.Wisdom,plas.DeathFails,plas.Alignment,plas.Level,plas.Name);
	_,err = DB.Query(query)
	if err != nil{
		fmt.Fprintf(w,"%q couldn't be created <br> <a href=\"./\"> Return to home? </a>",plas.Name)
	return
	}
	http.Redirect(w,r,"/viewCharacter",http.StatusSeeOther)
}


func main(){
	dbPass = "root:@/dnd"
	port := ":8080"
	http.HandleFunc("/routelogin",routeSelectHandler)
	http.HandleFunc("/makechar",charHandler)
	http.Handle("/",http.FileServer(http.Dir("./static")))
	fmt.Printf("Listening on port %s",port)
	http.ListenAndServe(port,nil)
}
