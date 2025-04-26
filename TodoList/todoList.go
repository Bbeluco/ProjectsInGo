package main

import (
	"encoding/csv"
	"os"
	"fmt"
	"time"
)

func openFileInAppendMode() (file *os.File, err error){
	file, err = os.OpenFile("database.csv", os.O_APPEND|os.O_WRONLY, 0644);
	
	if err != nil {
		fmt.Println("Error while trying to open file", err)
	}

	return file, err
}

func add() {
	appendFile, err := openFileInAppendMode();
	if err != nil {
		fmt.Println("Error while trying read database", err);
	}

	defer appendFile.Close()

	writer := csv.NewWriter(appendFile);
	
	defer writer.Flush();

	userTask := os.Args[2]

	row := []string {"4", userTask, time.Now().Format("09-07-2017")}
	err = writer.Write(row);

	if err != nil {
		fmt.Println("Error while trying to add data do database");
	}
	fmt.Println("New row inserted\n")
	fmt.Println(row);
}

func list() {
	fmt.Println("Hi, here's your TODO list!")
	file, err := os.Open("database.csv");
	if err != nil {
		fmt.Println("Error while opening database");
	}

	defer file.Close()

	reader := csv.NewReader(file);
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error while reading records");
	}

	for _, eachrecord := range records {
		fmt.Println(eachrecord)
	}
}

func main() {
	switch option := os.Args[1]; option {
		case "add":
			add();
		case "list":
			list();
		default:
			fmt.Println("Option not found, please try again")
	}
}