package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// Function to specify the table name otherwise it uses the name albums
func (Album) TableName() string {
	return "album"
}

var db *gorm.DB

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/recordings?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DBUSER"), os.Getenv("DBPASS"))

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Album{})

	// Operations
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)

	alb2, err := albumByID(albID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album added: %v\n", alb2)
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album
	if err := db.Where("artist = ?", name).Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(id int64) (Album, error) {
	var alb Album
	if err := db.First(&alb, id).Error; err != nil {
		return alb, fmt.Errorf("albumByID %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	if err := db.Create(&alb).Error; err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return alb.ID, nil
}
