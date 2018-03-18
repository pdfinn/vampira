package main

import (
	"fmt"
	"image"
	"os"
	"unicode"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func minu(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func region(a, b int) int {
	if a < b {
		return -1
	}
	if a == b {
		return 0
	}
	return 1
}

func acmeerror(s string, err error) {
	fmt.Fprintf(os.Stderr, "acme: %s: %v\n", s, err)
	// panic(fmt.Sprintf(os.Stderr, "acme: %s: %v\n", s, err))
}

var (
	prevmouse image.Point
	mousew    *Window
)

func clearmouse() {
	mousew = nil
}

func savemouse(w *Window) {
	prevmouse = mouse.Point
	mousew = w
}

func restoremouse(w *Window) bool {
	defer func() { mousew = nil }()
	if mousew != nil && mousew == w {
		display.MoveTo(prevmouse)
		return true
	}
	return false
}

func isalnum(c rune) bool {
	return unicode.IsNumber(c) || unicode.IsLetter(c)
}

func runestrchr(s []rune, r rune) int {
	for ret, sr := range s {
		if sr == r {
			return ret
		}
	}
	return -1
}

func utfrune(s []rune, r int) int {
	for i, c := range s {
		if c == rune(r) {
			return i
		}
	}
	return -1
}

func errorwin1(dir string, incl []string) *Window {
	var Lpluserrors = "+Errors"

	r := dir + string("/") + Lpluserrors
	w := lookfile(r)
	if w == nil {
		if len(row.col) == 0 {
			if row.Add(nil, -1) == nil {
				acmeerror("can't create column to make error window", nil)
			}
		}
		w = row.col[len(row.col)-1].Add(nil, nil, -1)
		w.filemenu = false
		w.SetName(r)
		xfidlog(w, "new")
	}
	for _, in := range incl {
		w.AddIncl(in)
	}
	w.autoindent = globalautoindent
	return w
}

/* make new window, if necessary; return with it locked */
func errorwin(md *MntDir, owner int) *Window {
	var w *Window

	for {
		if md == nil {
			w = errorwin1("", nil)
		} else {
			w = errorwin1(md.dir, md.incl)
		}
		w.Lock(owner)
		if w.col != nil {
			break
		}
		/* window was deleted too fast */
		w.Unlock()
	}
	return w
}

/*
 * Incoming window should be locked.
 * It will be unlocked and returned window
 * will be locked in its place.
 */
func errorwinforwin(w *Window) *Window {
	var (
		owner int
		incl  []string
		dir   string
		t     *Text
	)

	t = &w.body
	dir = t.DirName()
	if dir == "." { /* sigh */
		dir = ""
	}
	incl = []string{}
	for _, in := range w.incl {
		incl = append(incl, in)
	}
	owner = w.owner
	w.Unlock()
	for {
		w = errorwin1(dir, incl)
		w.Lock(owner)
		if w.col != nil {
			break
		}
		/* window deleted too fast */
		w.Unlock()
	}
	return w
}

/*
 * Heuristic city.
 */
func makenewwindow(t *Text) *Window {
	var (
		c               *Column
		w, bigw, emptyw *Window
		emptyb          *Text
		i, y, el        int
	)
	switch {
	case activecol != nil:
		c = activecol
	case seltext != nil && seltext.col != nil:
		c = seltext.col
	case t != nil && t.col != nil:
		c = t.col
	default:
		if len(row.col) == 0 && row.Add(nil, -1) == nil {
			acmeerror("can't make column", nil)
		}
		c = row.col[len(row.col)-1]
	}
	activecol = c
	if t == nil || t.w == nil || len(c.w) == 0 {
		return c.Add(nil, nil, -1)
	}

	/* find biggest window and biggest blank spot */
	emptyw = c.w[0]
	bigw = emptyw
	for i = 1; i < len(c.w); i++ {
		w = c.w[i]
		/* use >= to choose one near bottom of screen */
		if w.body.fr.MaxLines >= bigw.body.fr.MaxLines {
			bigw = w
		}
		if w.body.fr.MaxLines-w.body.fr.NLines >= emptyw.body.fr.MaxLines-emptyw.body.fr.NLines {
			emptyw = w
		}
	}
	emptyb = &emptyw.body
	el = emptyb.fr.MaxLines - emptyb.fr.NLines
	/* if empty space is big, use it */
	if el > 15 || (el > 3 && el > (bigw.body.fr.MaxLines-1)/2) {
		y = emptyb.fr.Rect.Min.Y + emptyb.fr.NLines*tagfont.Height
	} else {
		/* if this window is in column and isn't much smaller, split it */
		if t.col == c && t.w.r.Dy() > 2*bigw.r.Dy()/3 {
			bigw = t.w
		}
		y = (bigw.r.Min.Y + bigw.r.Max.Y) / 2
	}
	w = c.Add(nil, nil, y)
	if w.body.fr.MaxLines < 2 {
		w.col.Grow(w, 1)
	}
	return w
}

func mousescrollsize(nl int) int {
	// Unimpl()
	return 1
}

type Warning struct{
	md *MntDir
	buf Buffer
};

var warnings = []Warning{}

func flushwarnings() {
var (
	warn Warning
	w *Window
	t *Text
	owner, nr, q0, n int
)
	for _, warn=range warnings  {
		w = errorwin(warn.md, 'E');
		t = &w.body;
		owner = w.owner;
		if owner == 0 {
			w.owner = 'E';
		}
		w.Commit(t);
		/*
		 * Most commands don't generate much output. For instance,
		 * Edit ,>cat goes through /dev/cons and is already in blocks
		 * because of the i/o system, but a few can.  Edit ,p will
		 * put the entire result into a single hunk.  So it's worth doing
		 * this in blocks (and putting the text in a buffer in the first
		 * place), to avoid a big memory footprint.
		 */
		q0 = t.file.b.nc();
		for n = 0; n < warn.buf.nc(); n += nr {
			nr = warn.buf.nc() - n;
			if nr > RBUFSIZE {
				nr = RBUFSIZE;
			}
			r := warn.buf.Read(n, nr);
			_, nr = t.BsInsert(t.file.b.nc(), r, true);
		}
		t.Show(q0, t.file.b.nc(), true);
		t.w.SetTag();
		t.ScrDraw();
		w.owner = owner;
		w.dirty = false;
		w.Unlock();
		warn.buf.Close();
		if warn.md != nil {
			fsysdelid(warn.md);
		}
	}
	warnings = warnings[0:0]
}


func warning(md *MntDir, s string, args... interface{})() {
	r := []rune(fmt.Sprintf(s, args...))
	addwarningtext(md, r);
}

func addwarningtext(md *MntDir, r []rune) {
	for _, warn := range warnings {
		if(warn.md == md){
			warn.buf.Insert(warn.buf.nc(), r);
			return;
		}
	}
	warn := Warning{}
	warn.md = md;
	if(md!=nil) {
		fsysincid(md);
	}
	warnings = append(warnings, warn)
	warn.buf.Insert(0, r);
	cwarn <- 0
}