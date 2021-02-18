/**
 * @Author: LOFTER
 * @Description:
 * @File:  fileinfo
 * @Date: 2021/2/5 1.txt:49 下午
 */
package model

type FileInfo struct {
	BaseInfo
	FileName string `json:"file_name,omitempty"`
	FileSize string `json:"file_size,omitempty"`
	Path     string `json:"path,omitempty"`
	Hash     string `json:"hash"`
}
