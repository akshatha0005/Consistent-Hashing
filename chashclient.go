package main
import (
"fmt"
"crypto/md5"
"math"
"sort"
"net/http"
"errors"
"io/ioutil"
"encoding/json"
"github.com/julienschmidt/httprouter"
)
type keyval struct{
key string
value string
}

type Hashsort []uint32
func (h Hashsort) Len() int           { return len(h) }
func (h Hashsort) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Hashsort) Less(i, j int) bool { return h[i] < h[j] }

type Chash struct{
	circle map[uint32]string
	sortedKeys []uint32
	nodes []string
	weights map[string]int
}

type Response struct{
Key string `json:"key"`	
Value string `json:"value"`	
}

	
//creating circle of servers	
func New(nodes []string) *Chash {
	hashcircle := &Chash{
		circle:       make(map[uint32]string),
		sortedKeys: make([]uint32, 0),
		nodes:      nodes,
		weights:    make(map[string]int),
	}
	hashcircle.generateCircle()
	return hashcircle
}

//consisting hashing circle generator
func (h *Chash) generateCircle() {
	totalWeight := 0
	for _, node := range h.nodes {
		if weight, ok := h.weights[node]; ok {
			totalWeight += weight
		} else {
			totalWeight += 1
		}
	}

	for _, node := range h.nodes {
		weight := 1

		if _, ok := h.weights[node]; ok {
			weight = h.weights[node]
		}

		factor := math.Floor(float64(40*len(h.nodes)*weight) / float64(totalWeight))

		for j := 0; j < int(factor); j++ {
			nodeKey := fmt.Sprintf("%s-%d", node, j)
			bKey := hashDigest(nodeKey)

			for i := 0; i < 3; i++ {
				key := hashVal(bKey[i*4 : i*4+4])
				h.circle[key] = node
				h.sortedKeys = append(h.sortedKeys, key)
			}
		}
	}

	sort.Sort(Hashsort(h.sortedKeys))
}

//md5 hash generator(hashing algorithm)
func hashDigest(key string) []byte {
	m := md5.New()
	m.Write([]byte(key))
	return m.Sum(nil)
}


//select node for key based on hash value of key
func (h *Chash) GetNode(stringKey string) (node string, ok bool) {
	pos, ok := h.GetNodePos(stringKey)
	if !ok {
		return "", false
	}
	return h.circle[h.sortedKeys[pos]], true
}

//function returning node position for given key
func (h *Chash) GetNodePos(stringKey string) (pos int, ok bool) {
	if len(h.circle) == 0 {
		return 0, false
	}
	key := h.GenKey(stringKey)
	nodes := h.sortedKeys
	pos = sort.Search(len(nodes), func(i int) bool { return nodes[i] > key })
	if pos == len(nodes) {
		return 1, true
	} else {
		return pos, true
	}
}

//hash generator for key
func (h *Chash) GenKey(key string) uint32 {
	bKey := hashDigest(key)
	return hashVal(bKey[0:4])
}

//function for PUT call
func putkey(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
var url string
id:=p.ByName("keyid")
value:=p.ByName("value")
ring := New(servers)
server1,z := ring.GetNode(id)
if(z==true){
url =server1+"/"+"keys"+"/"+id+"/"+value}
req, err := http.NewRequest("PUT", url, nil)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Server not running")
    }
resp.Body.Close()

	// any status code 200..299 is "success", so fail on anything else
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(errors.New(http.StatusText(resp.StatusCode)))
	}

	fmt.Println("The key " + id + " is being placed into " + server1)
}

//function for GET call
func getkey(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
var url string
var u Response
id := p.ByName("keyid")
ring := New(servers)
server1,z := ring.GetNode(id)
if(z==true){
url =server1+"/"+"keys"+"/"+id
}
resp, err := http.Get(url)
 if err != nil {
 	fmt.Println("Server not running")
 }
 if resp.StatusCode >= 400 {
		// handle 404 
		fmt.Println("key not found")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("empty body, key not found")
	}
err = json.Unmarshal(body, &u)
    if err != nil {
       fmt.Println("cannot unmarshal json")
    }
	fmt.Println(u.Key,"==>",u.Value)
}

//hash value generator
func hashVal(bKey []byte) uint32 {

	return ((uint32(bKey[3]) << 24) |
		(uint32(bKey[2]) << 16) |
		(uint32(bKey[1]) << 8) |
		(uint32(bKey[0])))
}

var(
servers = []string{"http://localhost:3001","http://localhost:3000","http://localhost:3002"}
)

func main(){
mux := httprouter.New()
    mux.GET("/keys/:keyid", getkey)
    mux.PUT("/keys/:keyid/:value",putkey)
        server := http.Server{
        Addr:        "0.0.0.0:5000",
        Handler: mux,
    }
    server.ListenAndServe()

}