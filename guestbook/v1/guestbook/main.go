package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/codegangsta/negroni"
    "github.com/gorilla/mux"
    "github.com/xyproto/simpleredis"
)

var (
    masterPool *simpleredis.ConnectionPool
    slavePool  *simpleredis.ConnectionPool
)

func ListRangeHandler(rw http.ResponseWriter, req *http.Request) {
    key := mux.Vars(req)["key"]
    list := simpleredis.NewList(slavePool, key)
    members := list.GetAll()
    membersJSON, err := json.Marshal(members)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }
    rw.Header().Set("Content-Type", "application/json")
    rw.Write(membersJSON)
}

func ListPushHandler(rw http.ResponseWriter, req *http.Request) {
    key := mux.Vars(req)["key"]
    value := mux.Vars(req)["value"]
    list := simpleredis.NewList(masterPool, key)
    list.Add(value)
    rw.Write([]byte("OK"))
}

func InfoHandler(rw http.ResponseWriter, req *http.Request) {
    info := simpleredis.Info(slavePool)
    rw.Header().Set("Content-Type", "text/plain")
    rw.Write([]byte(info))
}

func EnvHandler(rw http.ResponseWriter, req *http.Request) {
    environment := make(map[string]string)
    for _, item := range os.Environ() {
        splits := strings.Split(item, "=")
        key := splits[0]
        val := strings.Join(splits[1:], "=")
        environment[key] = val
    }

    envJSON, err := json.Marshal(environment)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.Write(envJSON)
}

func main() {
    masterPool = simpleredis.NewConnectionPoolHost("redis-master:6379")
    slavePool = simpleredis.NewConnectionPoolHost("redis-slave:6379")

    defer masterPool.Close()
    defer slavePool.Close()

    r := mux.NewRouter()
    r.Path("/lrange/{key}").Methods("GET").HandlerFunc(ListRangeHandler)
    r.Path("/rpush/{key}/{value}").Methods("GET").HandlerFunc(ListPushHandler)
    r.Path("/info").Methods("GET").HandlerFunc(InfoHandler)
    r.Path("/env").Methods("GET").HandlerFunc(EnvHandler)

    n := negroni.Classic()
    n.UseHandler(r)
    n.Run(":3000")
}
