package libraries

import (
	"database/sql"
	"log"
	"os"

	"./../books"
)

const (
	getLibrariesQuery = "SELECT libraries.id, name, permission, usr FROM libraries JOIN permissions ON libraries.id=permissions.libraryid join library_members on libraries.ownerid=library_members.id WHERE permissions.permission & 1 and permissions.userid=(SELECT id from library_members join usersession on library_members.id=usersession.userid WHERE sessionkey=?)"
	getBreaks         = "SELECT BreakType, ValueType, Value FROM breaks WHERE libraryid=? AND (ValueType=? OR ValueType='ID')"
)

var logger = log.New(os.Stderr, "log: ", log.LstdFlags|log.Lshortfile)

//Break is a case break
type Break struct {
	LibraryID int    `json:"libaryid"`
	ValueType string `json:"valuetype"`
	Value     string `json:"value"`
	BreakType string `json:"breaktype"`
}

//EditedCases is a set of edited cases and ids to remove
type EditedCases struct {
	EditedCases []EditedCase `json:"editedcases"`
	ToRemoveIds []int64      `json:"toremoveids"`
	LibraryID   int64        `json:"libraryid"`
}

//EditedCase is an edited case
type EditedCase struct {
	ID              int64 `json:"id"`
	SpacerHeight    int64 `json:"spacerheight"`
	PaddingLeft     int64 `json:"paddingleft"`
	PaddingRight    int64 `json:"paddingright"`
	Width           int64 `json:"width"`
	ShelfHeight     int64 `json:"shelfheight"`
	NumberOfShelves int64 `json:"numberofshelves"`
	CaseNumber      int64 `json:"casenumber"`
}

//Bookcase is a bookcase
type Bookcase struct {
	ID                int64         `json:"id"`
	SpacerHeight      int64         `json:"spacerheight"`
	PaddingLeft       int64         `json:"paddingleft"`
	PaddingRight      int64         `json:"paddingright"`
	BookMargin        int64         `json:"bookmargin"`
	Width             int64         `json:"width"`
	Shelves           []Bookshelf   `json:"shelves"`
	AverageBookHeight float64       `json:"averagebookheight"`
	AverageBookWidth  float64       `json:"averagebookwidth"`
	Library           books.Library `json:"library"`
	CaseNumber        int64         `json:"casenumber"`
}

//Bookshelf is a shelf on a bookcase
type Bookshelf struct {
	ID     int64        `json:"id"`
	Height int64        `json:"height"`
	Books  []books.Book `json:"books"`
}

//GetCases gets cases
func GetCases(db *sql.DB, libraryid, sortMethod, session string) ([]Bookcase, error) {
	books, _, err := GetBooks("dewey", "both", "both", "yes", "both", "both", "both", "", "1", "-1", "0", "FIC", libraryid, session)
	breaks, err := GetBreaks(libraryid, sortMethod)
	if err != nil {
		logger.Printf("Error: %+v", err)
		return nil, err
	}
	query := "SELECT CaseId, Width, SpacerHeight, PaddingLeft, PaddingRight, BookMargin, CaseNumber, NumberOfShelves, ShelfHeight FROM bookcases WHERE libraryid=? ORDER BY CaseNumber"
	rows, err := db.Query(query, libraryid)
	if err != nil {
		logger.Printf("Error: %+v", err)
		return nil, err
	}
	dim, err := GetDimensions(libraryid)
	if err != nil {
		logger.Printf("Error: %+v", err)
		return nil, err
	}
	var cases []Bookcase
	for rows.Next() {
		var id, width, spacerHeight, paddingLeft, paddingRight, bookMargin, caseNumber, numberOfShelves, shelfHeight int64
		err = rows.Scan(&id, &width, &spacerHeight, &paddingLeft, &paddingRight, &bookMargin, &caseNumber, &numberOfShelves, &shelfHeight)
		if err != nil {
			logger.Printf("Error: %+v", err)
			return nil, err
		}
		bookcase := Bookcase{
			ID:                id,
			Width:             width,
			SpacerHeight:      spacerHeight,
			PaddingLeft:       paddingLeft,
			PaddingRight:      paddingRight,
			BookMargin:        bookMargin,
			AverageBookWidth:  dim["averagewidth"],
			AverageBookHeight: dim["averageheight"],
			CaseNumber:        caseNumber,
		}
		// 		shelfquery := "SELECT id, Height FROM shelves WHERE CaseId=? ORDER BY ShelfNumber"
		// 		shelfrows, err := db.Query(shelfquery, id)
		// 		if err != nil {
		// logger.Printf("Error: %+v", err)
		// return nil, err
		// 		}
		// 		for shelfrows.Next() {
		// 			var shelfid, height int64
		// 			shelfrows.Scan(&shelfid, &height)
		// 			bookcase.Shelves = append(bookcase.Shelves, Bookshelf{
		// 				ID:     shelfid,
		// 				Height: height,
		// 			})
		// 		}
		for i := 0; i < int(numberOfShelves); i++ {
			bookcase.Shelves = append(bookcase.Shelves, Bookshelf{
				ID:     0,
				Height: shelfHeight,
			})
		}
		cases = append(cases, bookcase)
	}
	index := 0
	x := 0
	for c, bookcase := range cases {
		breakcase := false
		for s := range bookcase.Shelves {
			if breakcase {
				break
			}
			x = int(bookcase.PaddingLeft)
			useWidth := int(dim["averagewidth"])
			if index < len(books) && books[index].Width != 0 {
				useWidth = int(books[index].Width)
			}
			for index < len(books) && useWidth+x <= int(bookcase.Width)-int(bookcase.PaddingRight) {
				cases[c].Shelves[s].Books = append(cases[c].Shelves[s].Books, books[index])
				x += useWidth
				index++
				useWidth = int(dim["averagewidth"])
				if index < len(books) && books[index].Width != 0 {
					useWidth = int(books[index].Width)
				}
				if breaktype, present := breaks[0][books[index-1].ID]; present {
					if breaktype == "CASE" {
						breakcase = true
					}
					break
				}
				if len(books) > index {
					if sortMethod == "DEWEY" {
						if breaktype, present := breaks[1][books[index-1].Dewey]; present && books[index-1].Dewey != books[index].Dewey {
							if breaktype == "CASE" {
								breakcase = true
							}
							break
						}
					}
					if sortMethod == "SERIES" {
						if breaktype, present := breaks[1][books[index-1].Series]; present && books[index-1].Series != books[index].Series {
							if breaktype == "CASE" {
								breakcase = true
							}
							break
						}
					}
				}
			}
		}
	}
	return cases, nil
}

//GetLibraries gets the libraries available to a user
func GetLibraries(db *sql.DB, session string) ([]Library, error) {
	var libraries []Library
	rows, err := db.Query(getLibrariesQuery, session)
	if err != nil {
		logger.Printf("%+v", err)
		return nil, err
	}
	for rows.Next() {
		var l Library
		if err := rows.Scan(&l.ID, &l.Name, &l.Permissions, &l.Owner); err != nil {
			logger.Printf("Error: %+v", err)
			return nil, err
		}
		libraries = append(libraries, l)
	}
	return libraries, nil
}

//SaveCases saves book cases
func SaveCases(db *sql.DB, cases EditedCases) error {
	libraryid := cases.LibraryID
	for _, c := range cases.EditedCases {
		id := c.ID
		caseNumber := c.CaseNumber
		width := c.Width
		spacerHeight := c.SpacerHeight
		paddingLeft := c.PaddingLeft
		paddingRight := c.PaddingRight
		shelfHeight := c.ShelfHeight
		numberOfShelves := c.NumberOfShelves
		if id == 0 {
			query := "INSERT INTO bookcases (casenumber, width, spacerheight, paddingleft, paddingright, libraryid, numberofshelves, shelfheight) VALUES (?,?,?,?,?,?,?,?)"
			_, err := db.Exec(query, caseNumber, width, spacerHeight, paddingLeft, paddingRight, libraryid, numberOfShelves, shelfHeight)
			if err != nil {
				logger.Printf("Error: %+v", err)
				return err
			}
		} else {
			query := "UPDATE bookcases SET casenumber=?, width=?, spacerheight=?, paddingleft=?, paddingright=?, libraryid=?, numberOfShelves=?, shelfheight=? WHERE caseid=?"
			_, err := db.Exec(query, caseNumber, width, spacerHeight, paddingLeft, paddingRight, libraryid, numberOfShelves, shelfHeight, id)
			if err != nil {
				logger.Printf("Error: %+v", err)
				return err
			}
		}
	}
	var removedCases []string
	for _, id := range cases.ToRemoveIds {
		i := strconv.FormatInt(id, 10)
		removedCases = append(removedCases, i)
	}
	if len(removedCases) > 0 {
		query := "DELETE FROM bookcases WHERE CaseId IN (" + strings.Join(removedCases, ",") + ")"
		_, err := db.Exec(query)
		if err != nil {
			logger.Printf("Error: %+v", err)
			return err
		}
	}
	return nil
}

//AddBreak adds a shelf break
func AddBreak(db *sql.DB, libraryid int, valuetype, value, breaktype string) error {
	query := "REPLACE INTO breaks (libaryid, valuetype, value, breaktype) VALUES (?,?,?,?)"
	_, err := db.Exec(query, libraryid, valuetype, value, breaktype)
	return err
}

//GetBreaks gets shelf breaks
func GetBreaks(db *sql.DB, libraryid, valuetype string) ([]map[string]string, error) {
	idBreaks := make(map[string]string)
	customBreaks := make(map[string]string)
	rows, err := db.Query(getBreaks, libraryid, valuetype)
	if err != nil {
		logger.Printf("Error: %+v", err)
		return nil, err
	}
	for rows.Next() {
		var breaktype string
		var valuetype string
		var value string
		rows.Scan(&breaktype, &valuetype, &value)
		if valuetype == "ID" {
			idBreaks[value] = breaktype
		} else {
			customBreaks[value] = breaktype
		}
	}
	return []map[string]string{idBreaks, customBreaks}, nil
}
