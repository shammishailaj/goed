package actions

import "github.com/tcolar/goed/core"

func EdActivateView(viewId int64, y, x int) {
	d(edActivateView{viewId: viewId, y: y, x: x})
}

func EdDelColCheck(colIndex int) {
	d(edDelColCheck{colIndex: colIndex})
}

func EdDelViewCheck(viewId int64) {
	d(edDelViewCheck{viewId: viewId})
}

func EdOpen(loc string, viewId int64, rel string, create bool) {
	d(edOpen{loc: loc, viewId: viewId, rel: rel, create: create})
}

// Retuns whether the editor can be quit.
func EdQuitCheck() bool {
	answer := make(chan (bool), 1)
	d(edQuitCheck{answer: answer})
	return <-answer
}

func EdRender() {
	d(edRender{})
}

func EdResize(h, w int) {
	d(edResize{h: h, w: w})
}

func EdSetStatus(status string) {
	d(edSetStatus{status: status, err: false})
}

func EdSetStatusErr(status string) {
	d(edSetStatus{status: status, err: true})
}

func EdSwapViews(view1Id, view2Id int64) {
	d(edSwapViews{view1Id: view1Id, view2Id: view2Id})
}

func EdTermFlush() {
	d(edTermFlush{})
}

func EdViewMove(viewId int64, y1, x1, y2, x2 int) {
	d(edViewMove{viewId: viewId, y1: y1, x1: x1, y2: y2, x2: x2})
}

func EdViewNavigate(mvmt core.CursorMvmt) {
	d(edViewNavigate{mvmt})
}

// ########  Impl ......

type edActivateView struct {
	viewId int64
	y, x   int
}

func (a edActivateView) Run() error {
	core.Ed.ViewActivate(a.viewId, a.y, a.x)
	return nil
}

type edDelColCheck struct {
	colIndex int
}

func (a edDelColCheck) Run() error {
	core.Ed.DelColCheckByIndex(a.colIndex)
	return nil
}

type edDelViewCheck struct {
	viewId int64
}

func (a edDelViewCheck) Run() error {
	core.Ed.DelViewCheck(a.viewId)
	return nil
}

type edOpen struct {
	loc, rel string
	viewId   int64
	create   bool
}

func (a edOpen) Run() error {
	_, err := core.Ed.Open(a.loc, a.viewId, a.rel, a.create)
	return err
}

type edQuitCheck struct {
	answer chan (bool)
}

func (a edQuitCheck) Run() error {
	a.answer <- core.Ed.QuitCheck()
	return nil
}

type edRender struct{}

func (a edRender) Run() error {
	core.Ed.Render()
	return nil
}

type edResize struct {
	h, w int
}

func (a edResize) Run() error {
	core.Ed.Resize(a.h, a.w)
	return nil
}

type edSetStatus struct {
	status string
	err    bool
}

func (a edSetStatus) Run() error {
	if a.err {
		core.Ed.SetStatusErr(a.status)
	} else {
		core.Ed.SetStatus(a.status)
	}
	return nil
}

type edSwapViews struct {
	view1Id, view2Id int64
}

func (a edSwapViews) Run() error {
	core.Ed.SwapViews(a.view1Id, a.view2Id)
	return nil
}

type edTermFlush struct{}

func (a edTermFlush) Run() error {
	core.Ed.TermFlush()
	return nil
}

type edViewMove struct {
	viewId         int64
	y1, x1, y2, x2 int
}

func (a edViewMove) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ViewMove(a.y1, a.x1, a.y2, a.x2)
	return nil
}

type edViewNavigate struct {
	mvmt core.CursorMvmt
}

func (a edViewNavigate) Run() error {
	core.Ed.ViewNavigate(a.mvmt)
	return nil
}
