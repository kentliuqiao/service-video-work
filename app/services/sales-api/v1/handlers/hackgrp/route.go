package hackgrp

import "github.com/dimfeld/httptreemux/v5"

func Routes(mux *httptreemux.ContextMux) {
	mux.GET("/hack", Hack)
}
