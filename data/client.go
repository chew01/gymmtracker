package data

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
	"gymmtracker/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Options struct {
	Stdout   bool
	LogPath  string
	LogFlags int
}

type Client struct {
	Options *Options
	utils.Logger
}

func NewClient(options *Options) *Client {
	return &Client{
		Options: options,
		Logger:  utils.NewLogger(options.Stdout, options.LogPath, options.LogFlags),
	}
}

func (c *Client) Record() {
	resp, err := FetchData()
	if err != nil {
		c.Warnln(err.Error())
	}

	count, err := InsertData(resp)
	if err != nil {
		c.Warnln(err.Error())
	}

	c.Printf("Inserted %d values into database", count)
}

func FetchData() (*FormattedResponse, error) {
	resp, err := http.Get("https://smartentry.org/status/api/metrics/gymmboxx")
	if err != nil {
		return nil, errors.New("failed GET request: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body to bytes: " + err.Error())
	}

	var formatted FormattedResponse
	err = json.Unmarshal(body, &formatted)
	if err != nil {
		return nil, errors.New("error unmarshalling response to struct: " + err.Error())
	}

	return &formatted, nil
}

func InsertData(resp *FormattedResponse) (int64, error) {
	// Connect to database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return 0, errors.New("error connecting to database: " + err.Error())
	}

	// Defer to end of function: close database connection
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			err = errors.New("error closing database connection: " + err.Error())
		}
	}(conn, context.Background())

	// Initialize array of outlet data to be added
	currentTime := time.Now()
	var rows [][]interface{}

	for _, outlet := range resp.Outlets {
		row := []interface{}{outlet.Name, outlet.Occupancy, outlet.Limit, outlet.Queue, currentTime}
		rows = append(rows, row)
	}

	// Copy prepared array into database
	copyCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"outlet_data"},
		[]string{"name", "occupancy", "occupancy_limit", "queue", "timestamp"},
		pgx.CopyFromRows(rows))
	if err != nil {
		return 0, errors.New("error copying into database: " + err.Error())
	}

	return copyCount, err
}
