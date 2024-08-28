// handlers/simulation_handler.go
package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/html"

	"gopkg.in/yaml.v2"
)

type Simulation struct {
	SimName             string `json:"sim_name"`
	BuildType           string `json:"build_type"`
	SimURL              string `json:"sim_url"`
	Timestamp           string `json:"created_at"`
	SupportedInterfaces string `json:"supportedinterfaces"`
}

type SimulationHandler struct {
	db *sql.DB
}

func NewSimulationHandler(db *sql.DB) *SimulationHandler {
	return &SimulationHandler{
		db: db,
	}
}

type DocumentHandler struct {
	db *sql.DB
}

func NewDocumentHandler(db *sql.DB) *DocumentHandler {
	return &DocumentHandler{
		db: db,
	}
}

type OnboardHandler struct {
	db *sql.DB
}

func NewOnboardHandler(db *sql.DB) *OnboardHandler {
	return &OnboardHandler{
		db: db,
	}
}

type OfferingsHandler struct {
	db *sql.DB
}

func NewofferingsHandler(db *sql.DB) *OfferingsHandler {
	return &OfferingsHandler{
		db: db,
	}
}

//_______________________________________________SIMULATION_CATALOG_______________________________________________________________________________________

//function to insert simulation entries in simulation catalog through post

func (h *SimulationHandler) CreateSimulation(c echo.Context) error {
	simulation := new(Simulation)
	if err := c.Bind(simulation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid simulation data")
	}

	insertStmt := "INSERT INTO simulation_catalog (sim_name, build_type, sim_url,created_at) VALUES ($1, $2, $3, $4)"
	createdat := time.Now()
	created_At := createdat
	_, err := h.db.Exec(insertStmt, simulation.SimName, simulation.BuildType, simulation.SimURL, created_At)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create simulation")
	}

	//return c.JSON(http.StatusCreated, simulation)
	return c.JSONPretty(http.StatusOK, map[string]string{"msg": " Inserted Successfully"}, " ")
}

// Get Simulation Handler

func (h *SimulationHandler) GetSimulations(c echo.Context) error {
	buildType := c.QueryParam("build_type")

	var query string
	var args []interface{}

	if buildType != "" {
		query = "SELECT sim_name, build_type, sim_url, created_at FROM simulation_catalog WHERE build_type = $1 ORDER BY created_at DESC"
		args = []interface{}{buildType}
	} else {
		query = "SELECT sim_name, build_type, sim_url, created_at FROM simulation_catalog ORDER BY created_at DESC"
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	simulations := make([]Simulation, 0)
	for rows.Next() {
		var sim Simulation
		var createdAt sql.NullTime

		if err := rows.Scan(&sim.SimName, &sim.BuildType, &sim.SimURL, &createdAt); err != nil {
			return err
		}
		// sim.Timestamp = createdAt.Format("Monday 02 January, 2006 03:04 PM")
		// // Format createdAt timestamp as a string in the desired format

		// Set the supported interfaces field
		sim.SupportedInterfaces = "cmedit, cmevents"

		if createdAt.Valid {
			sim.Timestamp = createdAt.Time.Format("2006-01-02 03:04")
		} else {
			sim.Timestamp = "" // Set empty string for NULL timestamps
		}

		simulations = append(simulations, sim)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(simulations) == 0 {
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Simulations not found"}, " ")
	}

	response := map[string]interface{}{
		"format":     "table",
		"table_name": "Simulation Catalog",
		"columns": []map[string]string{
			{
				"name":       "Simulation Name",
				"type":       "string",
				"field_name": "sim_name",
			},
			{
				"name":       "Build Type",
				"type":       "string",
				"field_name": "build_type",
			},
			{
				"name":       "Simulation Link",
				"type":       "url",
				"field_name": "sim_url",
			},
			{
				"name":       "Supported Interfaces",
				"type":       "string",
				"field_name": "supportedinterfaces",
			},
			{
				"name":       "Created At",
				"type":       "timestamp",
				"field_name": "created_at",
			},
		},
		"data": simulations,
	}

	return c.JSONPretty(http.StatusOK, response, " ")
}

//struct for installation documentation

type Documentation struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

//function to get installation documentation in the document table

func (h *DocumentHandler) GetProductDocumentationLink(c echo.Context) error {
	var documentation Documentation
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'Product Documentation'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "installation documentation not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve installation documentation")
	}

	//return c.JSONPretty(http.StatusOK, documentation, " ")
	return c.JSONPretty(http.StatusOK, map[string]string{"installation_documentation": documentation.DocumentLink}, " ")
}

//function to update installation documentation in the document table

func (h *DocumentHandler) UpdateProductDocumentationLink(c echo.Context) error {
	documentation := new(Documentation)
	if err := c.Bind(documentation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid documentation data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'Product Documentation'"
	_, err := h.db.Exec(updateStmt, documentation.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update installation documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "installation documentation updated successfully"}, " ")
}

//struct for onboarding document

type Onboarding struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

// function fo getting onboarding documentation link

func (h *OnboardHandler) GetOnboardingDocumentationLink(c echo.Context) error {
	var documentation Documentation
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'User Onboarding Hyperlink'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "User onboarding documentation not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user onboarding documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"user_onboarding_hyperlink": documentation.DocumentLink}, " ")
}

// function to update onboarding documentation link

func (h *OnboardHandler) UpdateOnboardingDocumentationLink(c echo.Context) error {
	onboarding := new(Onboarding)
	if err := c.Bind(onboarding); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid onboarding data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'User Onboarding Hyperlink'"
	_, err := h.db.Exec(updateStmt, onboarding.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user onboarding documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "User onboarding documentation updated successfully"}, " ")
}

// struct for offerings

type Offerings struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

// function for getting homepage link

func (h *OfferingsHandler) GetofferingsDocumentLink(c echo.Context) error {
	var documentation Offerings
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'Offerings Link'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Offerings link not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve Offerings Link")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"offerings_link": documentation.DocumentLink}, " ")
}

// function for updating the homepage link

func (h *OfferingsHandler) UpdateofferingsDocumentationLink(c echo.Context) error {
	offerings := new(Offerings)
	if err := c.Bind(offerings); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid homepage data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'Offerings Link'"
	_, err := h.db.Exec(updateStmt, offerings.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update Homepage Link")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "Offerings Link updated successfully"}, " ")
}

type SimCatalogHandler struct {
	db                *sql.DB
	lastYAMLHash      string
	lastSortedEntries []NameWithVersion
}

func NewSimCatalogHandler(db *sql.DB) *SimCatalogHandler {
	return &SimCatalogHandler{
		db: db,
	}
}

func (h *SimCatalogHandler) StartPeriodicComparison() {
	updateInterval := 2 * time.Minute
	go func() {
		for {
			time.Sleep(updateInterval)
			h.compareAndUpdateData()
		}
	}()
}

type Entry struct {
	APIVersion string   `yaml:"apiVersion"`
	Name       string   `yaml:"name"`
	Version    string   `yaml:"version"`
	URLs       []string `yaml:"urls"`
	Created    string   `yaml:"created"`
}

type NameWithVersion struct {
	NameWithVersion string
	URL             string
	Created         time.Time
}

type Config struct {
	APIVersion string             `yaml:"apiVersion"`
	Entries    map[string][]Entry `yaml:"entries"`
}

func (h *SimCatalogHandler) compareAndUpdateData() {
	yamlURL := "https://arm.seli.gic.ericsson.se/artifactory/proj-eric-oss-restsim-drop-helm-local/index.yaml"

	resp, err := http.Get(yamlURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := yaml.NewDecoder(resp.Body)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		log.Fatal(err)
	}

	entriesData, err := yaml.Marshal(config.Entries)
	if err != nil {
		log.Fatal(err)
	}

	currentHash := computeHash(entriesData)

	var entriesWithNameVersion []NameWithVersion
	for _, entries := range config.Entries {
		for _, entry := range entries {
			for _, url := range entry.URLs {
				createdTime, err := time.Parse(time.RFC3339Nano, entry.Created)
				if err != nil {
					log.Printf("Failed to parse created time: %s\n", err.Error())
					continue
				}
				nameWithVersion := fmt.Sprintf("eric-oss-restsim-release-%s", entry.Version)
				entryWithNameVersion := NameWithVersion{
					NameWithVersion: nameWithVersion,
					URL:             url,
					Created:         createdTime,
				}
				entriesWithNameVersion = append(entriesWithNameVersion, entryWithNameVersion)
			}
		}
	}

	sort.Slice(entriesWithNameVersion, func(i, j int) bool {
		return entriesWithNameVersion[i].Created.Before(entriesWithNameVersion[j].Created)
	})

	if currentHash != h.lastYAMLHash {
		h.lastYAMLHash = currentHash
		deltaEntries := findDelta(h.lastSortedEntries, entriesWithNameVersion)
		h.insertSimulations(deltaEntries)
		h.lastSortedEntries = entriesWithNameVersion
	}
}

func findDelta(oldEntries, newEntries []NameWithVersion) []NameWithVersion {
	var delta []NameWithVersion
	oldMap := make(map[NameWithVersion]bool)
	for _, entry := range oldEntries {
		oldMap[entry] = true
	}
	for _, entry := range newEntries {
		if !oldMap[entry] {
			delta = append(delta, entry)
		}
	}
	return delta
}

func computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (h *SimCatalogHandler) insertSimulations(entries []NameWithVersion) {
	insertStmt := "INSERT INTO simulation_catalog (sim_name, build_type, sim_url, created_at) VALUES ($1, 'dev', $2, $3)"
	for _, entry := range entries {
		simName := entry.NameWithVersion
		url := entry.URL
		createdAt := entry.Created
		fmt.Println(createdAt)
		_, err := h.db.Exec(insertStmt, simName, url, createdAt)
		if err != nil {
			log.Printf("Failed to insert entry: %s\n", err.Error())
		} else {
			log.Printf("Inserted: %s, %s, %s\n", simName, url, createdAt)
		}
	}
}

type DatasetCatalogHandler struct {
	db              *sql.DB
	previousEntries map[string]DataEntry // Store previous entries in memory using dataset names as keys
}

func NewDatasetCatalogHandler(db *sql.DB) *DatasetCatalogHandler {
	return &DatasetCatalogHandler{
		db:              db,
		previousEntries: make(map[string]DataEntry),
	}
}

func (h *DatasetCatalogHandler) StartPeriodicComparison1() {
	updateInterval := 3 * time.Minute
	go func() {
		for {
			time.Sleep(updateInterval)
			h.datasetcompareAndUpdateData()
		}
	}()
}

func (h *DatasetCatalogHandler) datasetcompareAndUpdateData() {
	mainURL := "https://arm1s11-eiffel004.eiffel.gic.ericsson.se:8443/nexus/content/repositories/nss/com/ericsson/nss/RestsimDatasets/"
	mainURLs, _ := extractUrlsFromHTML(mainURL)

	var outputUrls []string

	for _, url := range mainURLs {
		subURLs, _ := extractFilteredUrlsFromHTML(url)
		for _, subURL := range subURLs {
			xmlData, err := fetchXMLData(subURL)
			if err != nil {
				fmt.Println("Error fetching XML data from", subURL, ":", err)
				continue
			}

			var metadata Metadata
			err = xml.Unmarshal(xmlData, &metadata)
			if err != nil {
				fmt.Println("Error parsing XML data from", subURL, ":", err)
				continue
			}

			if metadata.Versioning.Release != "" {
				generatedURL := generateURL(mainURL, metadata.ArtifactID, metadata.Versioning.Release)
				htmlContent, err := fetchHTMLContent(generatedURL)
				if err != nil {
					fmt.Println("Error fetching HTML content from", generatedURL, ":", err)
					continue
				}
				subURLs, err := extractUrlsFromHTMLContent(htmlContent)
				if err != nil {
					fmt.Println("Error extracting URLs from HTML content of", generatedURL, ":", err)
					continue
				}
				for _, url := range subURLs {
					outputURL := url
					outputUrls = append(outputUrls, outputURL)
				}
			}
		}
	}
	var entries []DataEntry
	for _, url := range outputUrls {
		jsonData, err := fetchAndParseJSON(url)
		if err != nil {
			fmt.Println("Error fetching or parsing JSON data from", url, ":", err)
			continue
		}

		var entry DataEntry
		err = json.Unmarshal(jsonData, &entry)
		if err != nil {
			fmt.Println("Error parsing JSON data from", url, ":", err)
			continue
		}

		// Store the raw JSON data in the Metadata field
		entry.Metadata = make(map[string]interface{})
		err = json.Unmarshal(jsonData, &entry.Metadata)
		if err != nil {
			fmt.Println("Error parsing JSON data for metadata from", url, ":", err)
			continue
		}

		entries = append(entries, entry)
	}

	// Marshal the entries slice to JSON
	_, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Print formatted JSON array
	// log.Println(string(jsonData))

	// Compare new entries with previous entries
	for _, entry := range entries {
		// Check if the dataset name already exists in the previousEntries map
		prevEntry, exists := h.previousEntries[entry.DatasetName]
		if !exists || !reflect.DeepEqual(prevEntry, entry) {
			// Dataset entry does not exist or has changed, perform insert or update operation
			h.updateDataset(entry)
			// Update the previousEntries map with the current entry
			h.previousEntries[entry.DatasetName] = entry
		}
	}
}

//function to update the dataset entries

func (h *DatasetCatalogHandler) updateDataset(entry DataEntry) error {
	// Convert netypes and metadata maps to JSON strings
	netypesJSON, err := json.Marshal(entry.Netypes)
	if err != nil {
		log.Fatal(err)
	}

	metadataJSON, err := json.Marshal(entry.Metadata)
	if err != nil {
		log.Fatal(err)
	}

	// Perform the database insert or update operation here
	_, err = h.db.Exec("INSERT INTO dataset (datasetname, cellcount, netypes, networkelements, metadata) VALUES ($1, $2, $3::json, $4, $5::json) ON CONFLICT (datasetname) DO UPDATE SET cellcount = EXCLUDED.cellcount, netypes = EXCLUDED.netypes, networkelements = EXCLUDED.networkelements, metadata = EXCLUDED.metadata",
		entry.DatasetName, entry.CellCount, string(netypesJSON), entry.NetworkElements, string(metadataJSON))
	if err != nil {
		return err
	}
	return nil
}

func fetchAndParseJSON(url string) ([]byte, error) {
	jsonData, err := fetchJSONData(url)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func fetchJSONData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type URLObject struct {
	URL string `json:"url"`
}

func extractUrlsFromHTML(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var htmlContent strings.Builder
	_, err = io.Copy(&htmlContent, resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(htmlContent.String()))
	if err != nil {
		return nil, err
	}

	var tdValues []string
	var extractUrls func(*html.Node)
	extractUrls = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if strings.Contains(attr.Val, "metadata") {
						tdValues = append(tdValues, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractUrls(c)
		}
	}

	extractUrls(doc)
	return tdValues, nil
}

func extractFilteredUrlsFromHTML(url string) ([]string, error) {
	urls, err := extractUrlsFromHTML(url)
	if err != nil {
		return nil, err
	}

	var filteredUrls []string
	for _, u := range urls {
		if strings.HasSuffix(u, "/maven-metadata.xml") && !strings.HasSuffix(u, ".xml.md5") && !strings.HasSuffix(u, ".xml.sha1") {
			filteredUrls = append(filteredUrls, u)
		}
	}

	return filteredUrls, nil
}

type Metadata struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Versioning struct {
		Release     string   `xml:"release"`
		Versions    []string `xml:"versions>version"`
		LastUpdated string   `xml:"lastUpdated"`
	} `xml:"versioning"`
}

func fetchXMLData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return xmlData, nil
}

func generateURL(baseURL, artifactID string, version string) string {
	return fmt.Sprintf("%s%s/%s/", baseURL, artifactID, version)
}

func fetchHTMLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(htmlContent), nil
}

func extractUrlsFromHTMLContent(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var urls []string
	var extractUrls func(*html.Node)
	extractUrls = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if strings.HasSuffix(attr.Val, ".json") {
						urls = append(urls, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractUrls(c)
		}
	}

	extractUrls(doc)
	return urls, nil
}

type DataEntry struct {
	CellCount       string                 `json:"Cellcount"`
	DatasetName     string                 `json:"Datasetname"`
	NetworkElements string                 `json:"NetworkElements"`
	Netypes         map[string]interface{} `json:"Netypes"`
	Metadata        map[string]interface{} `json:"Metadata"`
}

type DataEntry1 struct {
	CellCount       int                    `json:"Cellcount"`
	DatasetName     string                 `json:"Datasetname"`
	NetworkElements int                    `json:"NetworkElements"`
	Netypes         string                 `json:"Netypes"`
	Metadata        map[string]interface{} `json:"Metadata"`
}
type DatasetHandler struct {
	db *sql.DB
}

func NewDatasetHandler(db *sql.DB) *DatasetHandler {
	return &DatasetHandler{
		db: db,
	}
}

//api to get dataset catalog items

func (h *DatasetHandler) GetDatasets(c echo.Context) error {
	// Query the database to retrieve datasets
	rows, err := h.db.Query("SELECT datasetname, cellcount, netypes, networkelements, metadata FROM dataset")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal Server Error",
		})
	}
	defer rows.Close()

	var datasets []DataEntry1
	for rows.Next() {
		var dataset DataEntry1
		var cellCount, netypes, networkElements, metadata sql.NullString

		err := rows.Scan(&dataset.DatasetName, &cellCount, &netypes, &networkElements, &metadata)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		// // Handle NULL values by replacing them with empty strings

		// Convert strings to numbers
		if cellCount.Valid {
			cellCountInt, err := strconv.Atoi(cellCount.String)
			if err == nil {
				dataset.CellCount = cellCountInt
			}
		}

		if networkElements.Valid {
			networkElementsInt, err := strconv.Atoi(networkElements.String)
			if err == nil {
				dataset.NetworkElements = networkElementsInt
			}
		}

		// Parse metadata column as JSON
		if metadata.Valid {
			if err := json.Unmarshal([]byte(metadata.String), &dataset.Metadata); err != nil {
				fmt.Println("Error parsing metadata JSON:", err)
			}
		}

		if netypes.Valid {
			var netypesMap map[string]interface{}
			if err := json.Unmarshal([]byte(netypes.String), &netypesMap); err != nil {
				fmt.Println("Error parsing netypes JSON:", err)
			}

			var netypesStringSlice []string
			for key := range netypesMap {
				netypesStringSlice = append(netypesStringSlice, key)
			}

			netypesFormatted := strings.Join(netypesStringSlice, ", ")
			dataset.Netypes = netypesFormatted
		}

		datasets = append(datasets, dataset)
	}

	response := map[string]interface{}{
		"format":     "table",
		"table_name": "Dataset Catalog",
		"columns": []map[string]string{
			{
				"name":       "Dataset Name",
				"type":       "string",
				"field_name": "Datasetname",
			},
			{
				"name":       " Cell Count",
				"type":       "number",
				"field_name": "Cellcount",
			},
			{
				"name":       "Network Elements",
				"type":       "number",
				"field_name": "NetworkElements",
			},
			{
				"name":       "NE Types",
				"type":       "string",
				"field_name": "Netypes",
			},
			{
				"name":       "Metadata",
				"type":       "json",
				"field_name": "Metadata",
			},
		},
		"data": datasets,
	}
	return c.JSONPretty(http.StatusOK, response, " ")
}


// db instance initialisation for simulation documentation

type SimDocHandler struct {
	db *sql.DB
}

func NewSimDocHandler(db *sql.DB) *SimDocHandler {
	return &SimDocHandler{
		db: db,
	}
}

//struct for simulation document

type SimDoc struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

// function for getting simulation documentation link

func (h *SimDocHandler) GetSimulationDocumentationLink(c echo.Context) error {
	var documentation SimDoc
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'Simulation Document'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Simulation documentation not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user onboarding documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"simulation_document": documentation.DocumentLink}, " ")
}

// function to update simulation documentation link

func (h *SimDocHandler) UpdateSimulationDocumentationLink(c echo.Context) error {
	simdoc := new(SimDoc)
	if err := c.Bind(simdoc); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid simulation data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'Simulation Document'"
	_, err := h.db.Exec(updateStmt, simdoc.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update simulation documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "Simulation documentation updated successfully"}, " ")
}

// db instance initialisation for dataset documentation

type DataDocHandler struct {
	db *sql.DB
}

func NewDataDocHandler(db *sql.DB) *DataDocHandler {
	return &DataDocHandler{
		db: db,
	}
}

//struct for dataset document

type DataDoc struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

// function fo getting dataset documentation link

func (h *DataDocHandler) GetDatasetDocumentationLink(c echo.Context) error {
	var documentation DataDoc
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'Dataset Document'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Dataset documentation not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user onboarding documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"dataset_document": documentation.DocumentLink}, " ")
}

// function to update simulation documentation link

func (h *DataDocHandler) UpdateDatasetDocumentationLink(c echo.Context) error {
	datadoc := new(DataDoc)
	if err := c.Bind(datadoc); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid onboarding data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'Dataset Document'"
	_, err := h.db.Exec(updateStmt, datadoc.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update dataset documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "Dataset documentation updated successfully"}, " ")
}


// db instance initialisation for byos documentation

type ByosDocHandler struct {
	db *sql.DB
}

func NewByosDocHandler(db *sql.DB) *ByosDocHandler {
	return &ByosDocHandler{
		db: db,
	}
}

//struct for byos document

type ByosDoc struct {
	DocumentName string `json:"document_name"`
	DocumentLink string `json:"document_link"`
}

// function for getting byos documentation link

func (h *ByosDocHandler) GetByosDocumentationLink(c echo.Context) error {
	var documentation ByosDoc
	err := h.db.QueryRow("SELECT document_link FROM document_table WHERE document_name = 'Byos Document'").Scan(&documentation.DocumentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "Byos documentation not found"}, " ")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user onboarding documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"byos_document": documentation.DocumentLink}, " ")
}

// function to update byos documentation link

func (h *ByosDocHandler) UpdateByosDocumentationLink(c echo.Context) error {
	byosdoc := new(ByosDoc)
	if err := c.Bind(byosdoc); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid onboarding data")
	}

	updateStmt := "UPDATE document_table SET document_link = $1 WHERE document_name = 'Byos Document'"
	_, err := h.db.Exec(updateStmt, byosdoc.DocumentLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update byos documentation")
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"msg": "Byos documentation updated successfully"}, " ")
}