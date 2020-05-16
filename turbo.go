package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Query to create the table in DB
const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS inventory
(
    Id int AUTO_INCREMENT PRIMARY KEY,
    
    name VARCHAR(50) NOT NULL,
    description VARCHAR(50) NOT NULL
)`

// Query to get all users
const getAllPartsQuery = `
SELECT * FROM inventory`

type Inventory struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var parts []Inventory

func getPartsHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		db := connectDB()
		err := db.Ping()
		if err != nil {
			fmt.Println("Coudnt ping to DB")
			panic(err.Error())
		}
		rows := getAllPartsfromDb(db)

		resultsParts := []Inventory{}
		setDataRow := Inventory{}

		for rows.Next() {
			rows.Scan(&setDataRow.Id, &setDataRow.Name, &setDataRow.Description)
			resultsParts = append(resultsParts, setDataRow)

		}
		if err := rows.Err(); err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		b, _ := json.Marshal(resultsParts)
		c.Writer.Write([]byte(string(b)))
	}
	return gin.HandlerFunc(fn)
}

func getSpecificPartHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.Writer.Write([]byte("Inserted"))
	}
	return gin.HandlerFunc(fn)
}

func addPartHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		name := c.Request.URL.Query().Get("name")
		description := c.Request.URL.Query().Get("description")

		db := connectDB()
		err := db.Ping()
		if err != nil {
			fmt.Println("Coudn't ping to DB")
			panic(err.Error())
		}
		var partDetails Inventory
		partDetails.Name = name
		partDetails.Description = description

		insertPartInDB(partDetails, db)
		c.Writer.Write([]byte("Inserted"))
	}
	return gin.HandlerFunc(fn)
}

func getAllPartsfromDb(db *sql.DB) *sql.Rows {
	bs, err := db.Query(getAllPartsQuery)
	if err != nil {
		panic(err.Error())
	}
	return bs
}

func insertPartInDB(partDetails Inventory, db *sql.DB) bool {
	stmt, err := db.Prepare("INSERT into inventory SET name=?,description=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, queryError := stmt.Exec(partDetails.Name, partDetails.Description)
	if queryError != nil {
		fmt.Println(queryError)
		return false
	}
	return true
}

func main() {
	// Open the connection
	db := connectDB()
	// perform a db.Query insert

	create, err := db.Query(tableCreationQuery)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	defer create.Close()

	part1 := Inventory{Name: "Turbo", Description: "G turbo", Id: 1}
	part2 := Inventory{Name: "Pipe", Description: "Silicon Pipe", Id: 2}
	part3 := Inventory{Name: "Bolt", Description: "Standard Bolt", Id: 3}

	// Insert parts in dB
	b := insertPartInDB(part1, db)
	if !b {
		fmt.Println("Data coudn't be inserted")
	}
	b = insertPartInDB(part2, db)
	if !b {
		fmt.Println("Data coudn't be inserted")
	}
	b = insertPartInDB(part3, db)
	if !b {
		fmt.Println("Data coudn't be inserted")
	}
	r := gin.Default()
	// Routes consist of a path and a handler function.
	r.LoadHTMLGlob("static\\*")
	r.GET("/", getHomeHandler)

	r.GET("/parts", getPartsHandler())
	r.GET("/partsbyid", getSpecificPartHandler())
	r.GET("/addpart", addPartHandler())

	// Bind to a port and pass our router in
	//log.Fatal(http.ListenAndServe(":8000", r))

	r.Run("0.0.0.0:8000")
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:puta007**@tcp(127.0.0.1:3306)/turbodb")
	if err != nil {
		panic(err.Error())
		defer db.Close()
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Coudn't ping to DB")
		panic(err.Error())
	}
	fmt.Println("Database connected.")
	return db
}

func getHomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{"Greeting": "Howdee  Punk", "PageTitle": "Turbo"})
}
