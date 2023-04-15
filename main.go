package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/adamlahbib/gitaz/cmd/create"
	"github.com/adamlahbib/gitaz/cmd/imaging"

	"github.com/joho/godotenv"
)

var githubAccessToken string

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	// Set return type JSON
	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, string(prettyJSON.Bytes()))

	// create github client
	log.Println(githubAccessToken)
	client := create.NewClient(githubAccessToken)

	// create deployment
	deployment, err := create.CreateDeployment(client, "asauce0972", "tp_net", "main", "production", "test")
	if err != nil {
		log.Panic(err)
	}

	deployment_id := deployment.GetID()

	// update deployment status to success
	deploymentStatus, err := create.CreateDeploymentStatus(client, "asauce0972", "tp_net", deployment_id, "success", "test", "http://localhost:8080/")
	if err != nil {
		log.Panic(err)
	}

	// set website
	_, err = create.SetWebsite(client, "asauce0972", "tp_net", "http://localhost:8080/")
	if err != nil {
		log.Panic(err)
	}

	// set checks

	_, err = create.SetChecks(client, "asauce0972", "tp_net", false)
	if err != nil {
		log.Panic(err)
	}

	// check if hook exists

	exists, err := create.HookExists(client, "asauce0972", "tp_net", "https://eovxryzicqvvnn7.m.pipedream.net")
	if err != nil {
		log.Panic(err)
	}

	if !exists {
		// create hook
		_, err = create.CreateHook(client, "asauce0972", "tp_net", "https://eovxryzicqvvnn7.m.pipedream.net", []string{"push"})
		if err != nil {
			log.Panic(err)
		}
	}

	// clone repo locally
	err = imaging.CloneRepo("https://github.com/asauce0972/tp_net.git", githubAccessToken, "repos")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// build image and push
	err = imaging.Build("repos", "tpnet", "latest")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// update checks status

	_, err = create.UpdateChecks(client, "asauce0972", "tp_net", true)
	if err != nil {
		log.Panic(err)
	}

	log.Println(deploymentStatus)

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github/">LOGIN</a>`)
}

func getGithubClientID() string {

	githubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientID
}

func getGithubClientSecret() string {

	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientSecret
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the environment variable
	githubClientID := getGithubClientID()

	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"http://localhost:3000/login/github/callback",
	)

	// add scopes X-OAuth-Scopes: repo, user
	redirectURL = fmt.Sprintf("%s&scope=%s", redirectURL, "repo,user")

	http.Redirect(w, r, redirectURL, 301)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken = getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	loggedinHandler(w, r, githubData)
}

func getGithubAccessToken(code string) string {

	clientID := getGithubClientID()
	clientSecret := getGithubClientSecret()

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Response body converted to stringified JSON
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user/repos",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	log.Println("Access token: ", accessToken)
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}

func main() {

	// Simply returns a link to the login route
	http.HandleFunc("/", rootHandler)

	// Login route
	http.HandleFunc("/login/github/", githubLoginHandler)

	// Github callback
	http.HandleFunc("/login/github/callback", githubCallbackHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		loggedinHandler(w, r, "")
	})

	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}
