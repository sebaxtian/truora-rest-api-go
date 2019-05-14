package infoserverctrls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	connectiondb "github.com/sebaxtian/truora-rest-api-go/db"

	"github.com/likexian/whois-go"

	infoserver "github.com/sebaxtian/truora-rest-api-go/structs"
)

// GetInfoServer an InfoServer with domain
func GetInfoServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := r.FormValue("domain")
		fmt.Println("GetInfoServer Domain: ", domain)

		// Call ssllabs api get data for domain
		respssllabs, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=​" + domain)
		if err != nil {
			// log.Fatal(err)
			// handle error
			fmt.Println("ERROR GetInfoServer Domain: ", err)
		}
		// its ok read data sslabs api
		datasslabs, err := ioutil.ReadAll(respssllabs.Body)
		respssllabs.Body.Close()
		if err != nil {
			// log.Fatal(err)
			// handle error
			fmt.Println("ERROR GetInfoServer Domain: ", err)
		}

		// Call index for domain
		respsdomain, err := http.Get("http://​" + domain)
		if err != nil {
			// log.Fatal(err)
			// handle error
			fmt.Println("ERROR GetIndex Domain: ", err)
		}
		// its ok read data index domain
		indexdomain, err := ioutil.ReadAll(respsdomain.Body)
		respsdomain.Body.Close()
		if err != nil {
			// log.Fatal(err)
			// handle error
			fmt.Println("ERROR GetIndex Domain: ", err)
		}
		// fmt.Println("Index Domain: ", indexdomain)
		// transform data sslabs api to string
		strindexdomain := string(indexdomain)
		// fmt.Println("Index Domain: ", strindexdomain)
		// Search Logo
		re := regexp.MustCompile(`href="[:\/\/A-Za-z0-9_\-\/.]*png" {1}`)
		regexstr := re.FindString(strindexdomain)
		logo := "without logo"
		if regexstr != "" {
			// fmt.Println("regexstr: ", regexstr)
			logo = strings.Split(strings.Split(regexstr, " ")[0], "\"")[1]
			// fmt.Println("Logo Index Domain: ", logo)
		}

		// its ok get InfoServer
		infoServer := getInfoServer(datasslabs)
		// Add information about server
		infoServer.Logo = logo

		// DB Operations
		db := connectiondb.DBConnect()
		infoServerDB := connectiondb.GetInfoServer(domain, db)
		// fmt.Println("infoServerDB: ", infoServerDB.SslGrade)
		// When infoServerDB not exist then create it
		if infoServerDB.SslGrade == "" {
			present := time.Now()
			lastUpdated := present.Format(time.RFC3339)
			infoServer.LastUpdated = lastUpdated
			connectiondb.CreateInfoServer(domain, infoServer, db)
			infoServerDB = connectiondb.GetInfoServer(domain, db)
		} else {
			// When infoServerDB exist then update it
			// Validate the previous state
			t := time.Now()
			present, _ := time.Parse(time.RFC3339, t.Format(time.RFC3339))
			past, _ := time.Parse(time.RFC3339, infoServerDB.LastUpdated)
			duration := present.Sub(past)
			fmt.Println("Present: ", present)
			fmt.Println("Past: ", past)
			fmt.Println("Duration: ", duration)
			fmt.Println("Duration.Hours: ", int(duration.Hours()))
			// fmt.Println("Duration.Minutes: ", int(duration.Minutes()))

			// Only or tests
			// if duration.Minutes() >= 0 {
			// Only update past one hour and grade ssl changed
			if duration.Hours() >= 0 {
				// Past one hour, validate grade ssl if changed
				// Only or tests
				// infoServer.SslGrade = "C"
				if infoServer.SslGrade != infoServerDB.PreviousSslGrade {
					infoServer.ServersChanged = true
					infoServer.PreviousSslGrade = infoServerDB.SslGrade
					present := time.Now()
					lastUpdated := present.Format(time.RFC3339)
					infoServer.LastUpdated = lastUpdated
				} else {
					infoServer.ServersChanged = false
					infoServer.PreviousSslGrade = infoServerDB.PreviousSslGrade
					infoServer.LastUpdated = infoServerDB.LastUpdated
				}
				connectiondb.UpdateInfoServer(domain, infoServer, db)
			}
		}
		infoServerDB = connectiondb.GetInfoServer(domain, db)
		// fmt.Println("InfoServer: ", infoServer)
		json.NewEncoder(w).Encode(infoServerDB)
	}
}

func getInfoServer(datassllabs []byte) infoserver.InfoServer {
	// transform data sslabs api to string
	// strsslabs := string(datassllabs)
	// fmt.Println("GetInfoServer Info Domain: ", strsslabs)

	// Decode data ssllabs
	var ssllabsjson interface{}
	json.Unmarshal(datassllabs, &ssllabsjson)
	// fmt.Println("GetInfoServer JSON Data Domain: ", ssllabsjson)
	// Map JSON data ssllabs
	mapssllabs := ssllabsjson.(map[string]interface{})
	// fmt.Println("GetInfoServer Ssllabs Domain: ", mapssllabs)

	// Create InfoServer
	infoServer := infoserver.InfoServer{}

	// Validate response status for api ssllabs
	if mapssllabs["status"] == "READY" {
		// Add information about server
		infoServer.IsDown = false

		// Ssl grade
		mapsslgrade := make(map[string]int)
		mapsslgrade["A"] = 80
		mapsslgrade["B"] = 65
		mapsslgrade["C"] = 50
		mapsslgrade["D"] = 35
		mapsslgrade["E"] = 20
		mapsslgrade["F"] = 19
		minusSslgrade := "A"

		// Map JSON data sslabs.endpoints
		mapssllabsendpoints := mapssllabs["endpoints"].([]interface{})
		// fmt.Println("GetInfoServer Endpoints Domain: ", mapssllabsendpoints[0])

		// Create data Server array
		servers := []infoserver.Server{}

		// Iterate for each endpoint
		for _, value := range mapssllabsendpoints {
			endpointI := value.(map[string]interface{})
			// fmt.Println("endpointI: ", endpointI)
			// fmt.Println("endpointI[ipAddress]: ", endpointI["ipAddress"])

			// Get info about endpoint
			ipAddress := endpointI["ipAddress"].(string)
			sslGrade := endpointI["grade"].(string)

			// Validate the minus ssl grade
			if mapsslgrade[sslGrade] <= mapsslgrade[minusSslgrade] {
				minusSslgrade = sslGrade
			}

			// Whois server
			whoisResult, err := whois.Whois(ipAddress)
			if err != nil {
				// log.Fatal(err)
				// handle error
				fmt.Println("ERROR Whois Server Endpoint: ", err)
				continue
			}
			// fmt.Println("Whois Server Endpoint: ", whoisResult)

			// Apply regular expression for extract information about server

			// Search Address
			re := regexp.MustCompile(`Address:[\w .]*`)
			serverAddress := strings.Trim(strings.Split(re.FindString(whoisResult), ":")[1], " ")
			// fmt.Println("Address: ", serverAddress)
			// Search OrgName
			re = regexp.MustCompile(`OrgName:[\w .]*`)
			serverOrgName := strings.Trim(strings.Split(re.FindString(whoisResult), ":")[1], " ")
			// fmt.Println("OrgName: ", serverOrgName)
			// Search Country
			re = regexp.MustCompile(`Country:[\w .]*`)
			serverCountry := strings.Trim(strings.Split(re.FindString(whoisResult), ":")[1], " ")
			// fmt.Println("Country: ", serverCountry)

			serverI := infoserver.Server{
				IPAddress: ipAddress,
				Address:   serverAddress,
				SslGrade:  sslGrade,
				Country:   serverCountry,
				Owner:     serverOrgName,
			}
			servers = append(servers, serverI)
		}

		// Add servers information
		infoServer.Servers = servers
		// Add information about server
		infoServer.SslGrade = minusSslgrade
	} else {
		// Add information about server
		infoServer.IsDown = true
		// status for api ssllabs not ready
		fmt.Println("Status for api ssllabs not ready Domain: ", mapssllabs["status"])
	}

	return infoServer
}
