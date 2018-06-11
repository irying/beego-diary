package main

import (
	"net/http"
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
)

func hello() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:router,
		ReadTimeout:setting.ReadTimeout,
		WriteTimeout:setting.WriteTimeout,
		MaxHeaderBytes:1 << 20,
	}

	s.ListenAndServe();
}
