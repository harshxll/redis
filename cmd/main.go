package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/mux"
)

var ErrorKeyNotExist = errors.New("key do not exist")
var ErrorKeyAlreadyExist = errors.New("key already exist")
type shard struct{
	sync.RWMutex
	m  map[string]string
}
const nshard = 7
var store []*shard


func init(){
	store = make([]*shard, nshard)
	for i := 0; i < len(store); i++ {
		store[i] = &shard{ m : make(map[string]string) }
	}
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/delete/{key}", keyValueDeleteHandler).Methods("GET")

	log.Println("server started sucessfully")
	http.ListenAndServe(":8080", router)
}



func ELFHash(key string) (int){
	var mask int = ((1 << 4) - 1) << 28
	var hash int = 0

	for i := 0; i < len(key);  i++{
		var c int = int(key[i])
		hash <<= 4
		hash |= c

		var g = hash & mask
		if g != 0 {
			var nhash = hash & (^mask)
			g >>= 8
			nhash ^= g
			hash = nhash
		}
	}

	return hash
}

func(s *shard)Put(key, value string) error{
	s.Lock()
	defer s.Unlock()

	_, ok := s.m[key]
	if ok {
		return ErrorKeyAlreadyExist
	}
	s.m[key] = value

	return nil
}

func (s *shard)Delete(key string) error{
	s.Lock()
	defer s.Unlock()

	_, ok := s.m[key]
	if !ok {
		return ErrorKeyNotExist
	} else { delete(s.m, key)}

	return nil
}

func (s *shard)Get(key string) (string, error){
	s.RLock()
	value, ok := s.m[key]
	s.RUnlock()

	if !ok {
		return "", ErrorKeyNotExist
	}
	return value, nil
}


func Get(key string) (string, error){
	var hash = ELFHash(key)
	hash %= nshard
	value, err := store[hash].Get(key)
	
	return value, err
}

func Put(key, value string) error{
	var hash = ELFHash(key)
	hash %= nshard
	err := store[hash].Put(key, value)
	return err
}

func Delete(key string) error{
	var hash = ELFHash(key)
	hash %= nshard
	err := store[hash].Delete(key)
	return err
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := Get(key)
	if errors.Is(err, ErrorKeyNotExist) {
		http.Error(w, err.Error(), http.StatusConflict )
		return
	}

	if err != nil{
		http.Error(w, 
			"internal server error", 
			http.StatusInternalServerError)
	}

	w.Write([]byte(value))
}


func keyValuePutHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))

	if err != nil {
		http.Error(w, 
			ErrorKeyAlreadyExist.Error(), 
			http.StatusContinue)
		return 
	}

	fmt.Fprintf(w, "key added successfully{key:%s}\n", key)
}


func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} 

	fmt.Fprintf(w, "sucessfully deleted key :%s\n", key)
}