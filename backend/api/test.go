package api

import (
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, req *http.Request){ 
	if req.Method != http.MethodGet {
		fmt.Fprint(w,"Error")
		return
	}

	fmt.Fprint(w,"WHAT?")
}