# shrtn
a url shortener service with rate limiter written in go. it provides a rest API to shorten long urls and redirects short urls to their original destinations.
    
### API

#### 1. shorten url

- **endpoint:** `post /`
    
- **request body:**
    
```json
{
  "url": "https://www.example.com/very/long/url"
}
```

- **response body:**
    
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

#### redirect
- **endpoint:** `get /{shortid}`
- redirects to the original url.
 
---

### prerequisites

- go installed
- postgresql running and accessible
- goose the database migration tool
      
### setup

1. clone the repository:
    
```bash
git clone <repo_url>
cd shrtn
```
    
2. create a `.env` file in the project root with the following content:
    
```bash
conn_str="postgres://postgres:postgres@localhost:5432/shrtn?sslmode=disable"
port="8080"
```

- adjust the connection string according to your `postgresql` credentials.
   
- ensure the database `shrtn` exists and then run goose:

```bash
goose -dir sql/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/shrtn" up    
```
       
3. build and run the application:
    
```bash
go build -o shrtn
./shrtn
```

the server will start and listen on the specified port (default: 8080).
    
### example

#### shorten a url

use `curl` or any http client to shorten a url:

```bash
curl -x post \
  -h "content-type: application/json" \
  -d '{"url":"https://github.com/"}' \
  http://localhost:8080/
```

**response:**

```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

#### redirect to original url

open the returned `short_url` in your browser or use curl as the following:

```bash
curl -l http://localhost:8080/abc123
```
