package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"html/template"
)

type Object interface{
	getName() string
	getWeight() float32
	getValue()	float32
	getDescription()	string
	getQuantity()	float32
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
	//How many one has
	Quantity	float32
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
func (i Item) getQuantity() float32{
	return i.Quantity
}

type Weapon struct{
	Base	Item
	Damage	string //Specifies number of sides the die has a spell with a save against it has this value
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
func (i Weapon) getQuantity() float32{
	return i.Base.Quantity
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
func (i Armor) getQuantity() float32{
	return i.Base.Quantity
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
	Speed	int
	Name	string
}

var dbPass string

func getInventory(player string, db *sql.DB) []Object{
	//Gets all of the inventory for a the given playername and returns it
	var ret []Object
	query := fmt.Sprintf("SELECT name,weight,value,description,quantity,damage,dist,ammo,ac,modifier FROM item where owner=%q",player)
	results,err := db.Query(query)
	if err != nil {
		return ret
	}
	for results.Next(){
		var nuItem Item
		var dam *string
		var dist *int
		var ammo *string
		var ac *int
		var mod *string
		if err := results.Scan(&nuItem.Name,&nuItem.Weight,&nuItem.Value,&nuItem.Description,&nuItem.Quantity,&dam,&dist,&ammo,&ac,&mod); err != nil{
			continue
		}
		if dist != nil {
			var weaponBox Weapon;
			weaponBox = Weapon{nuItem,*dam,*dist,*ammo,*mod}
			ret = append(ret,weaponBox)
			continue
		}
		if ac != nil {
			var armorBox Armor;
			armorBox = Armor{nuItem,*ac,*mod}
			ret = append(ret,armorBox)
			continue
		}
			ret = append(ret,nuItem)
	}
	return ret
}


func getPlayer(player string,db *sql.DB) Player{
	//Get's a specific palyer object from the database
	query := fmt.Sprintf("SELECT health,maxHealth,strength,dexterity,Intelligence,wisdom,proficiencies,clothes,deathFails,alignment,level,name,race,speed,charisma FROM player where name=%q;",player)
	res,err := db.Query(query)
	var nuPlayer Player;
	var prof *string;
	var cloth *string;
	if err != nil {
		return Player{}
	}
	res.Next()
	if err = res.Scan(&nuPlayer.Health,&nuPlayer.MaxHealth,&nuPlayer.Strength,&nuPlayer.Dexterity,&nuPlayer.Intelligence,&nuPlayer.Wisdom,&prof,&cloth,&nuPlayer.DeathFails,&nuPlayer.Alignment,&nuPlayer.Level,&nuPlayer.Name,&nuPlayer.Race,&nuPlayer.Speed,&nuPlayer.Charisma); err != nil {
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
		nuPlayer.Proficienies = strings.Split(*prof,",")
	}
	return nuPlayer
}


func routeSelectHandler(w http.ResponseWriter, r *http.Request){
	//This determines where to redirect to depending on the input page
	DB,err := sql.Open("mysql",dbPass)
	if err != nil {
		fmt.Fprintf(w,"ERROR: COULD NOT TOUCH DB")
		return
	}
	defer DB.Close()
	cur :=getPlayer(r.FormValue("name"),DB)
	if cur.Name == ""{
		http.Redirect(w,r,"/makecharacter.html",http.StatusFound)
	} else {
		http.Redirect(w,r,fmt.Sprintf("/viewCharacter?name=%s",cur.Name),http.StatusFound)
	}
}

func charHandler(w http.ResponseWriter, r *http.Request){
	//This makes the character in the database
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
	fmt.Sscanf(r.FormValue("Speed"),"%d",&plas.Speed)
	plas.DeathFails = -1
	plas.Alignment = r.FormValue("Alignment")
	plas.Class = r.FormValue("Class")
	plas.Name = r.FormValue("name")
	plas.Race = r.FormValue("race")
	query := fmt.Sprintf("INSERT INTO player(health,maxHealth,strength,dexterity,Intelligence,wisdom,deathFails,alignment,level,name,race,speed,charisma) VALUES(%d,%d,%d,%d,%d,%d,%d,%q,%d,%q,%q,%d,%d);",plas.Health,plas.MaxHealth,plas.Strength,plas.Dexterity,plas.Intelligence,plas.Wisdom,plas.DeathFails,plas.Alignment,plas.Level,plas.Name,plas.Race,plas.Speed,plas.Charisma);
	_,err = DB.Query(query)
	if err != nil{
		fmt.Fprintf(w,"<!DOCTYPE HTML><html><head><title>dnd-manager</title></head><body>%q couldn't be created <br> <a href=\"./\"> Return to home? </a></body></html>",plas.Name)
	return
	}
	http.Redirect(w,r,fmt.Sprintf("/viewCharacter?name=%s",plas.Name),http.StatusFound)
}

func finalHandler(w http.ResponseWriter,r *http.Request){
	//Handles serving of the main page
	nome := r.FormValue("name")
	if nome == "" {
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}
	DB,err := sql.Open("mysql",dbPass)
	if err != nil {
		fmt.Fprintf(w,"<!DOCTYPE HTML><html><head><title>dnd-manager</title> Something is wrong with the Database </body></html>")
		return
	}
	defer DB.Close()
	plas := getPlayer(nome,DB)
	if plas.Name != nome {
		fmt.Fprintf(w,"<!DOCTYPE HTML><html><head><title>dnd-manager</title> Player %q not found</body></html>",nome)
		return
	}
	top,err := template.ParseFiles("./ftop.html")
	if err != nil {
		fmt.Fprintf(w,"<!DOCTYPE HTML><html><head><title>dnd-manager</title> Template Couldn't be Parsed</body></html>")
		return

	}
	_ = top.Execute(w,plas)
	weap,_ := template.ParseFiles("./weaponrack.html")
	fmt.Fprintf(w,"<br><table><tr><th>Weapon</th><th>Damage</th><th>Ammo</th><th>Modifier</th><th>Description</th></tr>")
	for _,i := range plas.Inventory {
		switch v:= i.(type) {
			case Weapon:
				ammCount := 0.0
				for _,j := range plas.Inventory{
					if j.getName() == v.Ammo{
						ammCount+= float64(j.getQuantity())
					}
					v.Ammo = fmt.Sprintf("%s / %f",v.Ammo,ammCount)
				}
				weap.Execute(w,v)
			default:
				continue
		}
	}
	fmt.Fprintf(w,"</table>")
	fmt.Fprintf(w,"<br><table><tr><th>Namer</th><th>Quantity</th><th>Weight</th><th>Description</th><th>Description</th></tr>")
}


func main(){
	dbPass = "root:@/dnd"
	port := ":8080"
	http.HandleFunc("/routelogin",routeSelectHandler)
	http.HandleFunc("/makechar",charHandler)
	http.HandleFunc("/viewCharacter",finalHandler)
	http.Handle("/",http.FileServer(http.Dir("./static")))
	fmt.Printf("Listening on port %s",port)
	http.ListenAndServe(port,nil)
}
