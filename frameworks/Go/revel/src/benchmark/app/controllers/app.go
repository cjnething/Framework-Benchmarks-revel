package controllers

import (
	"database/sql"
	"math/rand"
	"net/http"
	"runtime"
	"sort"

	"github.com/agtorre/gocolorize"
	"github.com/revel/modules/db/app"
	"github.com/revel/revel"
)

// Revel test constants
const (
	WorldSelect        = "SELECT id,randomNumber FROM World where id=?"
	FortuneSelect      = "SELECT id,message FROM Fortune"
	WorldUpdate        = "UPDATE World SET randomNumber = ? where id = ?"
	WorldRowCount      = 10000
	MaxConnectionCount = 256

	// Added this based on
	// http://frameworkbenchmarks.readthedocs.org/en/latest/Project-Information/Framework-Tests/#specific-test-requirements
	jsonContentType  = "application/json"
	plainContentType = "text/plain"
)

var (
	worldStatement   *sql.Stmt
	fortuneStatement *sql.Stmt
	updateStatement  *sql.Stmt
)

type World struct {
	ID           uint16 `json:"id"`
	RandomNumber uint16 `json:"randomNumber"`
}

type Fortune struct {
	ID      uint16 `json:"id"`
	Message string `json:"message"`
}

type Fortunes []*Fortune

func (s Fortunes) Len() int      { return len(s) }
func (s Fortunes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByMessage struct{ Fortunes }

func (s ByMessage) Less(i, j int) bool { return s.Fortunes[i].Message < s.Fortunes[j].Message }

type App struct {
	*revel.Controller
}

// JSON method is Test 1 - JSON Serialization
func (c App) JSON() revel.Result {
	c.Response.ContentType = jsonContentType
	return c.RenderJson(map[string]string{"message": "Hello, World!"})
}

// Db method is for
// Test 2 - Single database query
func (c App) Db() revel.Result {
	var w World
	c.Response.ContentType = jsonContentType
	err := worldStatement.QueryRow(rand.Intn(WorldRowCount)+1).
		Scan(&w.ID, &w.RandomNumber)
	if err != nil {
		revel.ERROR.Println("Db query:", err)
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(map[string]string{"error": err.Error()})
	}
	return c.RenderJson(w)
}

// Dbs method is for
// Test 3 - Mulitple database query
func (c App) Dbs(queries int) revel.Result {
	c.Response.ContentType = jsonContentType
	queries = sanitizeQueryValue(queries)
	ww := make([]World, queries)
	for i := 0; i < queries; i++ {
		err := worldStatement.QueryRow(rand.Intn(WorldRowCount)+1).
			Scan(&ww[i].ID, &ww[i].RandomNumber)
		if err != nil {
			revel.ERROR.Println("Db queries:", err)
			c.Response.Status = http.StatusInternalServerError
			return c.RenderJson(map[string]string{"error": err.Error()})
		}
	}
	return c.RenderJson(ww)
}

// Fortune method is Test 4 - Fortunes
func (c App) Fortune() revel.Result {
	fortunes := make([]*Fortune, 0, 16)

	rows, err := fortuneStatement.Query()
	if err != nil {
		revel.ERROR.Println("Fortune:", err)
		return c.RenderError(err)
	}

	var fortune *Fortune
	for rows.Next() {
		fortune = new(Fortune)
		if err = rows.Scan(&fortune.ID, &fortune.Message); err != nil {
			revel.ERROR.Println("Fortune:", err)
			return c.RenderError(err)
		}
		fortunes = append(fortunes, fortune)
	}
	fortunes = append(fortunes, &Fortune{Message: "Additional fortune added at request time."})

	sort.Sort(ByMessage{fortunes})
	return c.Render(fortunes)
}

// Update method is for Test 5 - database updates
func (c App) Update(queries int) revel.Result {
	c.Response.ContentType = jsonContentType
	queries = sanitizeQueryValue(queries)
	ww := make([]World, queries)
	for i := 0; i < queries; i++ {
		err := worldStatement.QueryRow(rand.Intn(WorldRowCount)+1).
			Scan(&ww[i].ID, &ww[i].RandomNumber)
		if err != nil {
			revel.ERROR.Println("Update:", err)
			c.Response.Status = http.StatusInternalServerError
			return c.RenderJson(map[string]string{"error": err.Error()})
		}
		ww[i].RandomNumber = uint16(rand.Intn(WorldRowCount) + 1)
		updateStatement.Exec(ww[i].RandomNumber, ww[i].ID)
	}
	return c.RenderJson(ww)
}

// Plaintext method is Test 6 - plaintext
func (c App) Plaintext() revel.Result {
	c.Response.ContentType = plainContentType
	return c.RenderText("Hello, World!")
}

func sanitizeQueryValue(v int) int {
	// up to 500
	if v < 1 {
		v = 1
	} else if v > 500 {
		v = 500
	}
	return v
}

// ServerHeaderFilter adds 'Server' into response
var ServerHeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("Server", "Revel")

	fc[0](c, fc[1:])
}

func init() {
	revel.Filters = []revel.Filter{
		revel.RouterFilter,
		revel.ParamsFilter,
		ServerHeaderFilter,
		revel.ActionInvoker,
	}

	revel.OnAppStart(func() {
		var err error
		runtime.GOMAXPROCS(runtime.NumCPU())
		db.Init()
		db.Db.SetMaxIdleConns(MaxConnectionCount)
		if worldStatement, err = db.Db.Prepare(WorldSelect); err != nil {
			revel.ERROR.Fatalln(err)
		}
		if fortuneStatement, err = db.Db.Prepare(FortuneSelect); err != nil {
			revel.ERROR.Fatalln(err)
		}
		if updateStatement, err = db.Db.Prepare(WorldUpdate); err != nil {
			revel.ERROR.Fatalln(err)
		}

		gocolorize.SetPlain(true)
	})
}
