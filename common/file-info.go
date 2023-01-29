package common

type FileInfo struct {
	FileName string `json:"fileName"`
	Hash     string `json:"hash"`
}

func (f *FileInfo) Eq(o *FileInfo) bool {
	return f.Hash == o.Hash
}

