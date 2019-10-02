package boots

// todo 日志改用微服务
func init() {
	//logs.SetLogger(logs.AdapterFile,`{"filename":"` + beego.AppConfig.String("LogPath") + `", "level":7, "maxdays":10}`)
	//logs.EnableFuncCallDepth(true)
}
