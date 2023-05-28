package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/adamlahbib/gitaz/cmd/create"
	"github.com/adamlahbib/gitaz/cmd/imaging"
	"github.com/adamlahbib/gitaz/cmd/msg"
	"github.com/adamlahbib/gitaz/controllers"
	"github.com/adamlahbib/gitaz/initializers"
	"github.com/adamlahbib/gitaz/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"

	"github.com/joho/godotenv"
)

var githubAccessToken string

var mq msg.MQ

func main() {

	config := msg.MQConfig{
		Host: "localhost",
		Port: "5672",
		User: "guest",
		Pass: "guest",
	}

	mq = msg.MQ{}
	mq.Init(config)

	queue := mq.Queue("test")

	text := "Hello World!"

	body := []byte(text)
	queue.Publish(body)

	// Doing: converting http to Gin and implementing controllers

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"access-control-allow-origin, access-control-allow-headers, authorization, origin, content-type, accept"},
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Simply returns a link to the login route
	//http.HandleFunc("/", rootHandler)
	// //just a link no need for it, shows at the root of the server

	// Login route
	//http.HandleFunc("/login/github/", githubLoginHandler)
	r.GET("/login/github/", githubLoginHandler)

	// Github callback
	//http.HandleFunc("/login/github/callback", githubCallbackHandler)
	r.GET("/login/github/callback", githubCallbackHandler)

	// Route where the authenticated user is redirected to
	// http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
	// 	loggedinHandler(w, r, "")
	// })

	r.GET("/loggedin", func(c *gin.Context) {
		loggedinHandler(c.Writer, c.Request, "")
	})

	r.POST("/hook", controllers.GithubHooks)

	r.GET("/:username", controllers.GetUser)

	r.GET("/:username/update", func(c *gin.Context) {

	})

	r.GET("/:username/repos", controllers.FetchReposByUser)

	r.GET("/:username/:name", controllers.GetRepo)

	// r.GET("/user/:username/repos/refresh", ) to fetch user repos all over again maybe when user clicks refresh or sth
	// docker must be running and connected to a remote registry

	r.POST("/:username/:name/deploy", func(c *gin.Context) {
		deploymentHandler(c)
	})

	r.Run()

	// fmt.Println("[ UP ON PORT 3000 ]")
	// log.Panic(
	// 	http.ListenAndServe(":3000", nil),
	// )
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	initializers.ConnectDB()
	initializers.SyncDB()
}

func updateUserRepos(client *github.Client, user *github.User) {
	// Fetch repositories of the user
	// Fetch all repositories of the user with pagination
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(context.Background(), "", opt)
		if err != nil {
			log.Panic(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	log.Println(allRepos)

	var owner models.User

	initializers.DB.Where("username = ?", *user.Login).First(&owner)

	//saving user repos to DB
	for _, repo := range allRepos {
		controllers.AddUserRepository(models.Repo{
			OwnerID: owner.ID, // Set the ID of the user in the database
			Name:    *repo.Name,
			Branch:  *repo.DefaultBranch,
		})
	}
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {

	ctx := context.Background()

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

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Panic(err)
	}

	if !controllers.UserExists(*user.Login) {

		// save user data to db
		controllers.AddUser(
			models.User{
				Username: *user.Login,
				Avatar:   *user.AvatarURL,
				Token:    githubAccessToken,
			},
		)

		// fetch user repos, TODO call the function when user hits the update endpoint or sth as well.

		updateUserRepos(client, user)

	} else {
		// update user token
		controllers.UpdateUserToken(*user.Login, githubAccessToken)
	}
}

func deploymentHandler(c *gin.Context) {

	log.Println(c)

	lock := false // to lock heavy operations of the code, but they're functional

	//fetch user token
	var user models.User
	initializers.DB.Where("username = ?", c.Param("username")).First(&user)
	githubAccessToken := user.Token

	// create github client
	client := create.NewClient(githubAccessToken)

	// create deployment using controller function and specify the repo froeign key

	var repoTBD models.Repo

	initializers.DB.Where("name = ?", c.Param("name")).First(&repoTBD)

	controllers.AddDeployment(
		models.Deployment{
			RepoID:               repoTBD.ID,
			Stack:                c.Query("stack"),
			Version:              c.Query("version"),
			RepositoryURL:        c.Query("repository_url"),
			GithubToken:          c.Query("github_token"),
			ApplicationName:      c.Query("application_name"),
			RunCommand:           c.Query("run_command"),
			BuildCommand:         c.Query("build_command"),
			InstallCommand:       c.Query("install_command"),
			DependenciesFiles:    c.QueryArray("dependencies_files"),
			IsStatic:             c.Query("is_static") == "true",
			OutputDirectory:      c.Query("output_directory"),
			EnvironmentVariables: c.Query("environment_variables"),
			Port:                 c.Query("port"),
		},
	)

	jBody := map[string]interface{}{
		"stack":                 c.Query("stack"),
		"version":               c.Query("version"),
		"repository_url":        c.Query("repository_url"),
		"github_token":          c.Query("github_token"),
		"application_name":      c.Query("application_name"),
		"run_command":           c.Query("run_command"),
		"build_command":         c.Query("build_command"),
		"install_command":       c.Query("install_command"),
		"dependencies_files":    c.QueryArray("dependencies_files"),
		"is_static":             c.Query("is_static") == "true",
		"output_directory":      c.Query("output_directory"),
		"environment_variables": c.Query("environment_variables"),
		"port":                  c.Query("port"),
	}

	jsonString, err := json.Marshal(jBody)
	if err != nil {
		log.Panic(err)
	}

	queue := mq.Queue("test")
	queue.Publish(jsonString)

	// create cloudflare client
	clientCloudflare := create.NewCloudflareClient(os.Getenv("CLOUDFLARE_TOKEN"), os.Getenv("CLOUDFLARE_EMAIL"))

	log.Println(clientCloudflare)

	if !lock {
		// create deployment
		deployment, err := create.CreateDeployment(client, c.Param("username"), repoTBD.Name, repoTBD.Branch, "production", "test")
		if err != nil {
			log.Panic(err)
		}

		deployment_id := deployment.GetID()

		// update deployment status to success
		deploymentStatus, err := create.CreateDeploymentStatus(client, c.Param("username"), repoTBD.Name, deployment_id, "success", "test", "http://localhost:8080/")
		if err != nil {
			log.Panic(err)
		}

		if deploymentStatus.GetState() == "success" {
			// create status check
			_, err := create.CreateStatusCheck(client, c.Param("username"), repoTBD.Name, repoTBD.Branch, "pending", "test", "http://localhost:8080/", "test")
			if err != nil {
				log.Panic(err)
			}

			// enable depandabot alerts for the repo

		}

		// set website
		_, err = create.SetWebsite(client, c.Param("username"), repoTBD.Name, "https://"+c.Query("subdomain")+os.Getenv("KLI8NT_DOMAIN"))
		if err != nil {
			log.Panic(err)
		}

		exists, err := create.HookExists(client, c.Param("username"), repoTBD.Name, os.Getenv("HOOK_URL"))
		if err != nil {
			log.Panic(err)
		}

		if !exists {
			// create hook
			_, err = create.CreateHook(client, c.Param("username"), repoTBD.Name, os.Getenv("HOOK_URL"), []string{"push"})
			if err != nil {
				log.Panic(err)
			}
		}

		// clone repo locally
		err = imaging.CloneRepo("https://github.com/"+c.Param("username")+"/"+repoTBD.Name+".git", githubAccessToken, repoTBD.Name)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// build image and push
		err = imaging.Build(repoTBD.Name, repoTBD.Name, "latest")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// create status check success
		if deploymentStatus.GetState() == "success" {
			_, err = create.CreateStatusCheck(client, c.Param("username"), repoTBD.Name, repoTBD.Branch, "success", "test", "http://localhost:8080/", "test")
			if err != nil {
				log.Panic(err)
			}
		}

		// create a record in cloudflare
		err = create.CreateCNAME(clientCloudflare, os.Getenv("CLOUDFLARE_ZONEID"), c.Query("subdomain"), repoTBD.Name+"."+c.Param("username")+".repl.co", false)
		if err != nil {
			log.Panic(err)
		}

		log.Println(deploymentStatus)
	}
}

// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, `<a href="/login/github/">LOGIN</a>`)
// }

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

// func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get the environment variable
// 	githubClientID := getGithubClientID()

// 	// Create the dynamic redirect URL for login
// 	redirectURL := fmt.Sprintf(
// 		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
// 		githubClientID,
// 		"http://localhost:3000/login/github/callback",
// 	)

// 	// add scopes X-OAuth-Scopes: repo, user
// 	redirectURL = fmt.Sprintf("%s&scope=%s", redirectURL, "repo,user")

// 	http.Redirect(w, r, redirectURL, 301)
// }

func githubLoginHandler(c *gin.Context) {
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

	// Redirect to the dynamic URL
	c.Redirect(http.StatusMovedPermanently, redirectURL)
}

// func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")

// 	githubAccessToken = getGithubAccessToken(code)

// 	githubData := getGithubData(githubAccessToken)

// 	loggedinHandler(w, r, githubData)
// }

func githubCallbackHandler(c *gin.Context) {
	code := c.Query("code")

	githubAccessToken = getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	loggedinHandler(c.Writer, c.Request, githubData)
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
