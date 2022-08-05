package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		// TODO - measure the time using the histogram that you have defined and a prometheus timer

		// TODO - increment counter for request rate

		// TODO - increment the counter vec for HTTP methods

		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		strA := r.URL.Query().Get("a")
		strB := r.URL.Query().Get("b")

		a, errA := strconv.Atoi(strA)
		b, errB := strconv.Atoi(strB)

		if errA != nil || errB != nil {
			// TODO - increment error counter for "/sum"
			return
		}

		fmt.Fprintf(w, fmt.Sprintf("%d + %d = %d", a, b, a+b))
	})

	// TODO BONUS - using the os package configure the port using an environment variable
	//              you will have to provide a default value if this is empty. Then configure the
	//              the port in docker-compose and try to run it again
	fmt.Printf("Server running (port=8080)\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
