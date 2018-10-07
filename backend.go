package main

import (
	"fmt"
	"net/http"
)

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

type Weapon struct{
	Item
	Damage	int //Specifies number of sides the die has a spell with a save against it has this value
	Range	int //in feet
	Ammo	string //The type of ammunition this is tracked seperatly in the case of spells this is a spell slot
	Mod	string //This is the trait that is used as the modifier
}

type Armor struct{
	Item
	AC	int
	Mod	string
}


type Player struct {
	Inventory	[]Item
	Health	int
	MaxHealth	int
	Strength	int
	Dexterity	int
	Intellegence	int
	Wisdom	int
	Charisma	int
	Proficienies	[]string//THis is a list of strings specifying proficiency
	Clothes	Armor //THis is the item that is currently equiped
	DeathFails int //If this is ever over three player is dead, if it is set to a negative value then the player is stable
	Alignment	string
	Class	string //Currently doesn't do anything
	Race	string //Currently useless
	Level	int //The level of the character (does not support multiclassing)
}

func routeSelectHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,":)%q","Lorem ipsum")
}


func main(){
	port := ":8080"
	http.HandleFunc("/routelogin",routeSelectHandler)
	http.Handle("/",http.FileServer(http.Dir("./static")))
	fmt.Printf("Listening on port %s",port)
	http.ListenAndServe(port,nil)
}
