package meta

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]*FileMeta

func init() {
	fileMetas = make(map[string]*FileMeta)
}

// 修改文件的信息
func UpdateFileMeta(fileMeta *FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

// 通过sha1获取文件的元信息对象
func GetFileMeta(fileSha1 string) *FileMeta {
	return fileMetas[fileSha1]
}

func RemoveFileMeta(fileHash string) {
	delete(fileMetas, fileHash)
}
