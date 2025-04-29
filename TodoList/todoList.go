package main

import (
	"encoding/csv"
	"os"
	"fmt"
	"time"
	"strconv"
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

	records, err := fileReader();
	if err != nil {
		fmt.Println("Error while trying to read current records");
	}

	nextIndex, err := strconv.Atoi(records[len(records) - 1:][0][0]);

	defer appendFile.Close()

	writer := csv.NewWriter(appendFile);
	
	defer writer.Flush();

	userTask := os.Args[2]

	row := []string {strconv.Itoa(nextIndex + 1), userTask, time.Now().Format("01-02-2006")}
	err = writer.Write(row);

	if err != nil {
		fmt.Println("Error while trying to add data do database");
	}
	fmt.Println("New row inserted\n")
	fmt.Println(row);
}

func fileReader() (records [][]string, err error){
	file, err := os.Open("database.csv");
	if err != nil {
		fmt.Println("Error while opening database");
	}

	defer file.Close()

	reader := csv.NewReader(file);
	records, err = reader.ReadAll()

	if err != nil {
		fmt.Println("Error while reading records");
		return nil, err
	}

	return records, nil
}

func list() {
	fmt.Println("Hi, here's your TODO list!")

	records, err := fileReader();

	if err != nil {
		return;
	}

	for _, eachrecord := range records {
		fmt.Println(eachrecord)
	}
}

func help() {
	//colors
	var Reset = "\033[0m"
	var Red = "\033[31m" 
	var Green = "\033[32m" 
	var Yellow = "\033[33m" 

	fmt.Println("How to usage: [go run . COMMAND]")
	fmt.Println("Available commands:")
	fmt.Println();
	fmt.Printf(Green + "list - " + Reset)
	fmt.Println("See all the list itens")
	fmt.Printf(Yellow + "add 'Activity HERE'- " + Reset)
	fmt.Println("Insert new todo activity");
	fmt.Printf(Red + "delete 'ID HERE' - " + Reset)
	fmt.Println("Remove an entry from the list")
}

func main() {
	switch option := os.Args[1]; option {
		case "add":
			add();
		case "list":
			list();
		case "help":
			help();
		default:
			fmt.Println("Option not found, please try again")
	}
}