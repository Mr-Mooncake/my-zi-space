package main

import (
	"net/http"
	"fmt"
	"log"
)

 func main(){
	
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}	