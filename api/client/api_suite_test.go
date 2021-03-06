package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
	. "gopkg.in/check.v1"
)

func vidStr(vid int64) string {
	return fmt.Sprintf("%d", vid)
}

func Test(t *testing.T) { TestingT(t) }

type ApiSuite struct {
	id      int64 // instance
	ftext   []string
	dirView int64
}

var _ = Suite(&ApiSuite{})

const refFile string = "../../test_data/file1.txt"

func (s *ApiSuite) SetUpSuite(c *C) {
	// reference file
	b, _ := ioutil.ReadFile(refFile)
	for _, str := range bytes.Split(b, []byte{'\n'}) {
		s.ftext = append(s.ftext, string(str))
	}
	if len(s.ftext[len(s.ftext)-1]) == 0 {
		s.ftext = s.ftext[:len(s.ftext)-1]
	}
	// start mock editor
	s.id = time.Now().Unix()
	core.Testing = true
	core.InitHome(s.id)
	core.Ed = ui.NewMockEditor()
	core.Bus = actions.NewActionBus()
	actions.RegisterActions()
	apiServer := api.Api{}
	apiServer.Start()
	core.Ed.Start([]string{})
	s.dirView = core.Ed.Views()[0]
}

func (s *ApiSuite) SetUpTest(c *C) {
	actions.Ar.EdActionBusFlush()
	// Put the editor back into known state (only dir view open)
	for _, v := range actions.Ar.EdViews() {
		if v != s.dirView {
			actions.Ar.EdDelView(v, false)
		}
	}
	if len(actions.Ar.EdViews()) > 1 {
		debugViews()
	}
	assert.Eq(c, len(actions.Ar.EdViews()), 1)
	actions.Ar.ViewClearSelections(s.dirView)
	actions.Ar.ViewSetScrollPos(s.dirView, 1, 1)
	actions.Ar.ViewSetCursorPos(s.dirView, 1, 1)
}

func (s *ApiSuite) TestNoSuchAction(c *C) {
	res, err := Action(s.id, []string{"foobar"})
	assert.NotNil(c, err)
	assert.Eq(c, len(res), 0)
}

func (s *ApiSuite) TestEdit(c *C) {
	done := false
	completed := make(chan struct{})
	go func() {
		err := Edit(s.id, "test_data", "fooedit")
		done = true
		assert.Nil(c, err)
		close(completed)
	}()
	vid := int64(-1)
	// view should open up and stay open until the view is closed
	// at which time the open action should be completed
	loc, _ := filepath.Abs("test_data/fooedit")
	for vid == -1 {
		vids := actions.Ar.EdViewsByLoc(loc)
		if len(vids) > 0 {
			vid = vids[0]
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.False(c, done)
	actions.Ar.EdDelView(vid, true)
	select {
	case <-time.After(5 * time.Second):
		c.Fatal("timeout waiting for edit to complete.")
	case <-completed: // good
	}
}

func (s *ApiSuite) TestOpen(c *C) {
	err := Open(s.id, "test_data", "empty.txt")
	assert.Nil(c, err)
	loc, _ := filepath.Abs("./test_data/empty.txt")
	vid := actions.Ar.EdViewsByLoc(loc)[0]
	assert.NotEq(c, vid, "-1")
	actions.Ar.EdDelView(vid, true)
}

func (s *ApiSuite) openFile1(c *C) int64 {
	vid := actions.Ar.EdOpen(refFile, -1, "", false)
	assert.NotEq(c, vid, int64(-1))
	return vid
}
