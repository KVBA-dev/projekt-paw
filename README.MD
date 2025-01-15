# guess who
guess who is a simple game of guess who, written in Go.

## Installation and setup
1. Install Go from https://go.dev/doc/install
2. Install templ and air
```sh
$ go install github.com/a-h/templ/cmd/templ@latest
$ go install github.com/air-verse/air@latest
```
3. Clone this repository
```sh
git clone https://github.com/KVBA-dev/projekt-paw.git 
```
4. Create a file "env/schema.sql" and put following there
```sql
CREATE TABLE IF NOT EXISTS users (
 	id INTEGER PRIMARY KEY AUTOINCREMENT,
 	name TEXT NOT NULL, 
	games_played INTEGER DEFAULT 0,
 	games_won INTEGER DEFAULT 0,
 	passhash TEXT NOT NULL 
);
```
5. Get Go modules
```sh
$ go mod download
$ go mod tidy
```
6. Run the application
```sh
$ air
```

## Usage
By default this app will run on port 8080, this is configurable in cmd/main.go

Open the app with a browser, get some friends, and play along!

## Future improvements
- cleaner UI
- different flags for more balanced gameplay
  - randomised per game?
- progress tracking
- alternative game modes
