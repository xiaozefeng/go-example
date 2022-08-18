package main

import "fmt"

func main() {
	nginxServer := NewNginx(&Application{}, 10)
	appStatusURL := "/app/status"
	createUserURL := "/create/user"
	httpCode, body := nginxServer.handleRequest(appStatusURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.handleRequest(appStatusURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.handleRequest(appStatusURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.handleRequest(createUserURL, "POST")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.handleRequest(createUserURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)
}

type server interface {
	handleRequest(url, method string) (int, string)
}
type Nginx struct {
	app               *Application
	maxAllowedRequest int
	rateLimiter       map[string]int
}

func (n *Nginx) handleRequest(url, method string) (int, string) {
	allowed := n.checkRateLimiter(url, method)
	if !allowed {
		return 403, "Forbidden"
	}
	return n.app.handleRequest(url, method)
}

func (n *Nginx) checkRateLimiter(url string, method string) bool {
	key := fmt.Sprintf("%s:%s", url, method)
	if val, ok := n.rateLimiter[key]; !ok {
		n.rateLimiter[key] = 1
	} else {
		n.rateLimiter[key] = val + 1
	}
	return n.rateLimiter[key] < n.maxAllowedRequest
}

func NewNginx(app *Application, maxAllowedRequest int) *Nginx {
	return &Nginx{app: app, maxAllowedRequest: maxAllowedRequest, rateLimiter: make(map[string]int)}
}

type Application struct {
}

func (a *Application) handleRequest(url, method string) (int, string) {
	if url == "/app/status" && method == "GET" {
		return 200, "OK"
	}

	if url == "/create/user" && method == "POST" {
		return 201, "User Created"
	}
	return 404, "Not Found"
}
