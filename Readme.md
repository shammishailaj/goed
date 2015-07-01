## GoEd 
Terminal based code/text editor

It's currently in the "prototyping" state and barely usable, there is a lot of 
prototype code that is meant to be replaced and many features yet to be developed.

The end goal is to be powerful yet easy to use.

One of the main source of inspiration is Acme, however it will be more 
configurable and less mouse reliant.

That is not to say it will be "keyboard only" ant it will likely not rely on 
"modes" as Vi and Emacs do.

Early screenshot (6/2/2015): 
![Screenshot](https://raw.github.com/tcolar/goed/master/screenshot.png)

## Eventing design

Make it possible to send but also intercept events.

    [Event]
      InstanceId
      ViewId
      //Workdir
      EventName
      EventAliases
      Triggers
      Shortcuts (kb)
      []Args
      
Events (View):
  - Append(string)
  - Bounds() (y,x,y2,x2)  (ie: moveView)
  - BufferLoc() string
  - Cursor() (line, col)
  - Dirty() bool
  - Insert(line, col, text)
  - LineCount() int
  - Reload(<newpath> string)
  - Remove(ln1, col2, ln2, col2)
  - Render() ?? 
  - Save(<newloc> string)
  - Selections() ([] y,x,y2,x2)
  - Slice() []string
  - SrcLoc() string
  - Title() string 
  - Wipe()
  - Workdir() string
    
Events (Editor) :
  - NewView(<col>?, file ?) -> viewId -> only smart or allow specify column, ratio etc .. ?
  - CloseView(viewId)
  - CurView()
  - GoedHome() string
  - Open(file, viewid, string) ?? -> view ?
  - Render() ?
  - SetStatus(string)
  - SetStatusErr(string)
  - Version / ApiVersion
  - ViewByLoc / FindView() / FindOrCreateView() ?
  - Views()
  - Sed(string) ?

Triggers:
  - DirtiedView (changed)
  - SavedView
  - New view / Opened file
  - Closed view
  - Started goed instance
  - Closed goed instance
  
Sam? : http://doc.cat-v.org/bell_labs/sam_lang_tutorial/   

####case study 1 : gofmt
    Trigger: SavedView
    Match: "*.go"
    Action: "gofmt.sh"
    Type: "before | after | instead" ?

Gofmt.sh:
    #!/bin/bash
    gofmt -w $GOEDBUFFER_LOC
    if [ $? ne 0 ] ; then
      exit $?
    fi
    goed reload $GOED_VIEW_ID

####case study 2 : fortune on startup
    Trigger: Startup
    Match: ""
    Action: "fortune.sh"

Fortune.sh
    #!/bin/bash
    viewid=`goed newview` # col, size ??
    fortune > $GOED_HOME/buffers/$viewid
    goed reload $GOED_VIEW_ID
    
case study 3 : structs/funcs viewer
    Trigger: SavedView
    Match: "*.go"
    Action: "structure.sh"
    Type: "after"

structure.sh
    #!/bin/bash
    viewid=`goed findview "structure"` # find or create ?
    grep -rniH "struct" $GOED_BUFFER_LOC > $GOED_HOME/buffers/$viewid
    goed reload $viewid
    
case study 4 : date script
    #!/bin/bash 
    # called from cron every whatever
    viewid=`goed findview "structure"`
    date > $GOED_HOME/buffers/$viewid
    goed reload $viewid 
    
