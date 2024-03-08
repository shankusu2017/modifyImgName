package main

const (
	DEBUG_MODEL    = false
	FUNC_DEL_SAME  = true
	ROUTINUSCNT    = 256 // 限制并发量，避免对文件系统造成过大压力
	EXIF_TRY_TIMES = 3   // exif 有时候调用不会返回正确的值，这里设置尝试调用的次数
)
