package main

import (
	_ "SuperEsbAdminWeb/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var RedirectHttp = func(ctx *context.Context) {
	if !ctx.Input.IsSecure() {
		// no need focToa+VI8Fovi/ECVJRXcpIgldlicRploflkqJz3NqeUeb4l3mH+AtN9Xha+y/R9Br an additional '/' between domain and uri
		url := "https://" + ctx.Input.Domain() + ":" + beego.AppConfig.String("HttpsPort") + ctx.Input.URI()
		ctx.Redirect(302, url)
	}
}
var AddCorsHeaders = func(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Headers", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func main() {
	if beego.AppConfig.String("EnableHTTPS") == "true" {
		beego.InsertFilter("/", beego.BeforeRouter, RedirectHttp) // for http://mysite
		beego.InsertFilter("*", beego.BeforeRouter, RedirectHttp) // for http://mysite/*
	}

	// Add the CORS headers filter
	beego.InsertFilter("*", beego.BeforeRouter, AddCorsHeaders)

	beego.SetStaticPath("/static", "frontend/dist")
	beego.Run()
}
