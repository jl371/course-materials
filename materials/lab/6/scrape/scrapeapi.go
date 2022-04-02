package scrape

// scrapeapi.go HAS TEN TODOS - TODO_5-TODO_14 and an OPTIONAL "ADVANCED" ASK

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "os"
    "path/filepath"
    "strconv"
	"github.com/gorilla/mux"
	"regexp"
)


//==========================================================================\\

// Helper function walk function, modfied from Chap 7 BHG to enable passing in of
// additional parameter http responsewriter; also appends items to global Files and 
// if responsewriter is passed, outputs to http 

func walkFn(w http.ResponseWriter) filepath.WalkFunc {
    return func(path string, f os.FileInfo, err error) error {
        w.Header().Set("Content-Type", "application/json")

        for _, r := range regexes {
            if r.MatchString(path) {
                var tfile FileInfo
                dir, filename := filepath.Split(path)
                tfile.Filename = string(filename)
                tfile.Location = string(dir)

                //TODO_5: As it currently stands the same file can be added to the array more than once 
                //TODO_5: Prevent this from happening by checking if the file AND location already exist as a single record

				

                Files = append(Files, tfile)

                if w != nil && len(Files)>0 {

                    //TODO_6: The current key value is the LEN of Files (this terrible); 
                    //TODO_6: Create some variable to track how many files have been added
                    w.Write([]byte(`"`+(strconv.FormatInt(int64(len(Files)), 10))+`":  `))
                    json.NewEncoder(w).Encode(tfile)
                    w.Write([]byte(`,`))

                } 
                
                log.Printf("[+] HIT: %s\n", path)

            }

        }
        return nil
    }

}


func walkFn2(w http.ResponseWriter, query string) filepath.WalkFunc {
    return func(path string, f os.FileInfo, err error) error {
		w.Header().Set("Content-Type", "application/json")

		r := regexp.MustCompile(`(?i)` + query)

        
		if r.MatchString(path) {
			var tfile FileInfo
			dir, filename := filepath.Split(path)
			tfile.Filename = string(filename)
			tfile.Location = string(dir)

			

			Files = append(Files, tfile)

			if w != nil && len(Files)>0 {

				
				w.Write([]byte(`"`+(strconv.FormatInt(int64(len(Files)), 10))+`":  `))
				json.NewEncoder(w).Encode(tfile)
				w.Write([]byte(`,`))

			} 
			
			log.Printf("[+] HIT: %s\n", path)

		}

        
		return nil
	}
}

//==========================================================================\\

func APISTATUS(w http.ResponseWriter, r *http.Request) {

	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{ "status" : "API is up and running ",`))
    var regexstrings []string
    
    for _, regex := range regexes{
        regexstrings = append(regexstrings, regex.String())
    }

    w.Write([]byte(` "regexs" :`))
    json.NewEncoder(w).Encode(regexstrings)
    w.Write([]byte(`}`))
	log.Println(regexes)

}


func MainPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
    w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)


	fmt.Fprintf(w, "<html><body><H1>Welcome to the file page endpoint.</H1></body>")
}


func FindFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
    q, ok := r.URL.Query()["q"]

    w.WriteHeader(http.StatusOK)
    if ok && len(q[0]) > 0 {
        log.Printf("Entering search with query=%s",q[0])

        // ADVANCED: Create a function in scrape.go that returns a list of file locations; call and use the result here
        // e.g., func finder(query string) []string { ... }
		var fnd bool = false
        for _, File := range Files {
		    if File.Filename == q[0] {
                json.NewEncoder(w).Encode(File.Location)
                fnd = true
		    }
        }
        //TODO_9: Handle when no matches exist; print a useful json response to the user; hint you might need a "FOUND variable" to check here ...
		if fnd == false {
			w.Write([]byte(`"File not found"`))
		}

    } else {
        // didn't pass in a search term, show all that you've found
        w.Write([]byte(`"files":`))    
        json.NewEncoder(w).Encode(Files)
    }
}

func IndexFiles(w http.ResponseWriter, r *http.Request) {
    log.Printf("Entering %s end point", r.URL.Path)
    w.Header().Set("Content-Type", "application/json")

    location, locOK := r.URL.Query()["location"]
	regEx, regExOK := r.URL.Query()["regex"]
    
    //TODO_10: Currently there is a huge risk with this code ... namely, we can search from the root /
    //TODO_10: Assume the location passed starts at /home/ (or in Windows pick some "safe?" location)
    //TODO_10: something like ...  rootDir string := "???"
    //TODO_10: create another variable and append location[0] to rootDir (where appropriate) to patch this hole

    if locOK && len(location[0]) > 0 {
        w.WriteHeader(http.StatusOK)

    } else {
        w.WriteHeader(http.StatusFailedDependency)
        w.Write([]byte(`{ "parameters" : {"required": "location",`))    
        w.Write([]byte(`"optional": "regex"},`))    
        w.Write([]byte(`"examples" : { "required": "/indexer?location=/xyz",`))
        w.Write([]byte(`"optional": "/indexer?location=/xyz&regex=(i?).md"}}`))
        return 
    }

    //wrapper to make "nice json"
    w.Write([]byte(`{ `))
    
    if regExOK {
		filepath.Walk(location[0], walkFn2(regEx[0]))
	} else filepath.Walk(location[0], walkFn(w)) {
		log.Panicln(err)
	}

    //wrapper to make "nice json"
    w.Write([]byte(` "status": "completed"} `))

}

func Resetarray(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering reset endpoint")
	w.Header().Set("Content-Type", "application/json")
	resetRegEx()
	Files = nil
}

func Clear(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering reset endpoint")
	w.Header().Set("Content-Type", "application/json")
	clearRegEx()
}


func AddRegEx(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	addRegEx(`(?i)`+params["regex"])

}