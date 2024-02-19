package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{data: map[string]dollars{"shoes": 50, "socks": 5}}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)

	//Additional Handlers
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/read", db.read)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	data map[string]dollars
	mu   sync.RWMutex
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {

	//Lock Before Reading
	db.mu.RLock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	db.mu.RUnlock()
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.mu.RLock()
	price, ok := db.data[item]
	db.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}

	fmt.Fprintf(w, "%s\n", price)

}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	updatedPrice := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(updatedPrice, 32)

	if err != nil {
		fmt.Println("Error encountered: ", err)
		fmt.Fprintf(w, "'%s' is not a valid price, please try again.\n", updatedPrice)
		return
	}
	//Lock for reading
	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	//Lock for writing
	db.mu.Lock()
	db.data[item] = dollars(price)
	db.mu.Unlock()

	fmt.Fprintf(w, "Updated the price for %s to %s\n", item, db.data[item])
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(newPrice, 32)

	if err != nil {
		fmt.Println("Error encountered: ", err)
		return
	}
	// Lock for writing
	db.mu.Lock()
	db.data[item] = dollars(price)
	db.mu.Unlock()

	fmt.Fprintf(w, "New item added!\nItem Name: %s \nPrice: %s\n", item, db.data[item])
}

func (db *database) read(w http.ResponseWriter, req *http.Request) {
	db.list(w, req)
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	// Lock for writing
	db.mu.Lock()
	delete(db.data, item)
	db.mu.Unlock()

	fmt.Fprintf(w, "Item Removed!\nItem Name: %s \n", item)
}
