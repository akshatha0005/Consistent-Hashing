package main
import (
"fmt"
"net/http"
"github.com/julienschmidt/httprouter"
"encoding/json"
)
var Hashmap =map[string]string{}
var Hashmap1 =map[string]string{}
var Hashmap2 =map[string]string{}
type Response struct{
Key string `json:"key"`	
Value string `json:"value"`	
}	
func Getkeyid(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
	id:=p.ByName("keyid")
	r:= Response{
		Key:id,
		Value:Hashmap[id],
	}
	reply, _ := json.Marshal(r)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", reply)	
}

func Getkeys(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
	var r Response
	var x []Response
	for key, value := range Hashmap {
    r = Response{
    Key: key, 
    Value:value,
}
x = append(x,r)
}
	
	reply, _ := json.Marshal(x)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", reply)	
}
func Putkey(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
id:=p.ByName("keyid")
value:=p.ByName("value")
Hashmap[id]=value
rw.WriteHeader(200)	
}

func Getkeyid1(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    id:=p.ByName("keyid")
    r:= Response{
        Key:id,
        Value:Hashmap1[id],
    }
    reply, _ := json.Marshal(r)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", reply)    
}

func Getkeys1(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    var r Response
    var x []Response
    for key, value := range Hashmap1 {
    r = Response{
    Key: key, 
    Value:value,
}
x = append(x,r)
}
    
    reply, _ := json.Marshal(x)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", reply)    
}


func Putkey1(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
id:=p.ByName("keyid")
value:=p.ByName("value")
Hashmap1[id]=value
rw.WriteHeader(200) 
}

func Getkeyid2(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    id:=p.ByName("keyid")
    r:= Response{
        Key:id,
        Value:Hashmap2[id],
    }
    reply, _ := json.Marshal(r)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", reply)    
}

func Getkeys2(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    var r Response
    var x []Response
    for key, value := range Hashmap2 {
    r = Response{
    Key: key, 
    Value:value,
}
x = append(x,r)
}
    
    reply, _ := json.Marshal(x)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", reply)    
}


func Putkey2(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
id:=p.ByName("keyid")
value:=p.ByName("value")
Hashmap2[id]=value
rw.WriteHeader(200) 
}
func main() {
    go func(){
    mux := httprouter.New()
    mux.GET("/keys/:keyid", Getkeyid)
    mux.GET("/keys", Getkeys)
    mux.PUT("/keys/:keyid/:value", Putkey)
        server := http.Server{
        Addr:        "0.0.0.0:3000",
        Handler: mux,
    }
    server.ListenAndServe()

}()

go func(){
   
    mux1 := httprouter.New()
    mux1.GET("/keys/:keyid", Getkeyid1)
    mux1.GET("/keys", Getkeys1)
    mux1.PUT("/keys/:keyid/:value", Putkey1)
        server1 := http.Server{
        Addr:        "0.0.0.0:3001",
        Handler: mux1,
    }
    server1.ListenAndServe()
}()

   
    mux2 := httprouter.New()
    mux2.GET("/keys/:keyid", Getkeyid2)
    mux2.GET("/keys", Getkeys2)
    mux2.PUT("/keys/:keyid/:value", Putkey2)
        server2 := http.Server{
        Addr:        "0.0.0.0:3002",
        Handler: mux2,
    }
    server2.ListenAndServe()
 

}