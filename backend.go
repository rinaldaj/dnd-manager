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
	ProficiencyB	int
	Proficienies	[]string//THis is a list of strings specifying proficiency
	Clothes	Armor //THis is the item that is currently equiped
	Stable	bool// this is only used when Health == 0 it is to determine if we must make death saves
	DeathFails int //If this is ever over three player is dead
	Alignment	string
	Class	string //Currently doesn't do anything
}




func main(){

fmt.Println("Hello World")
}
