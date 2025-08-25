package fileutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/h2non/filetype"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"go-chat/internal/model"
)

type ffprobeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

// ParseFile 解析上传文件，返回统一的 model.File 信息，包含音视频(需要安装ffmpeg https://ffmpeg.org/index.html)时长等
func ParseFile(fileHeader *multipart.FileHeader) (*model.File, error) {
	f := &model.File{
		Name: fileHeader.Filename,
		Size: uint64(fileHeader.Size),
	}

	// 获取扩展名，去掉点
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if len(ext) > 0 {
		f.Ext = ext[1:]
	}

	// 打开文件，读取头部字节检测类型
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	head := make([]byte, 261)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取文件头失败: %w", err)
	}
	head = head[:n]

	kind, err := filetype.Match(head)
	if err != nil {
		return nil, fmt.Errorf("检测文件类型失败: %w", err)
	}

	if kind != filetype.Unknown {
		f.Mime = kind.MIME.Value
		f.Ext = kind.Extension
	} else {
		f.Mime = "application/octet-stream"
	}

	// 重新打开文件提取图片宽高
	if strings.HasPrefix(f.Mime, "image/") {
		if _, err := file.Seek(0, io.SeekStart); err == nil {
			cfg, _, err := image.DecodeConfig(file)
			if err == nil {
				w := uint(cfg.Width)
				h := uint(cfg.Height)
				f.Width = &w
				f.Height = &h
			}
		}
	}

	// 处理音视频时长
	if strings.HasPrefix(f.Mime, "audio/") || strings.HasPrefix(f.Mime, "video/") {
		duration, err := getMediaDuration(fileHeader)
		if err == nil {
			f.Duration = &duration
		} else {
			fmt.Printf("getMediaDuration error: %v\n", err)
		}
	}

	// 根据 MIME 和扩展名分类文件类型
	f.Type = classifyByMime(f.Mime, f.Ext)

	return f, nil
}

// getMediaDuration 使用 ffprobe 获取音视频时长
func getMediaDuration(fileHeader *multipart.FileHeader) (float64, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(fileHeader.Filename))
	if err != nil {
		return 0, err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	// 拷贝上传文件内容到临时文件
	if _, err := io.Copy(tmpFile, file); err != nil {
		return 0, err
	}

	if err := tmpFile.Sync(); err != nil {
		return 0, err
	}

	cmd := exec.Command("ffprobe",
		"-hide_banner",
		"-loglevel", "error",
		"-show_entries", "format=duration",
		"-of", "json",
		tmpFile.Name(),
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("ffprobe 执行失败: %w, 输出: %s", err, out.String())
	}

	var probeOutput ffprobeOutput
	if err := json.Unmarshal(out.Bytes(), &probeOutput); err != nil {
		return 0, fmt.Errorf("ffprobe 输出解析失败: %w", err)
	}

	duration, err := strconv.ParseFloat(probeOutput.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("时长转换失败: %w", err)
	}

	fmt.Printf("ffprobe duration: %f seconds\n", duration) // 调试用

	return duration, nil
}

// classifyByMime 按 MIME 和扩展名分类文件类型
func classifyByMime(mime string, ext string) string {
	if strings.HasPrefix(mime, "image/") {
		return "image"
	}
	if strings.HasPrefix(mime, "audio/") {
		return "audio"
	}
	if strings.HasPrefix(mime, "video/") {
		return "video"
	}
	if strings.HasPrefix(mime, "application/") || strings.HasPrefix(mime, "text/") {
		switch ext {
		case "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "rtf":
			return "document"
		case "zip", "rar", "7z", "tar", "gz", "bz2":
			return "archive"
		default:
			if isCodeFile(ext) {
				return "code"
			}
			return "file"
		}
	}
	return "file"
}

// isCodeFile 判断是否源码文件
func isCodeFile(ext string) bool {
	codeExts := []string{"go", "js", "ts", "java", "py", "c", "cpp", "cs", "rb", "php", "html", "css", "json", "xml", "sh", "bat", "md", "yml"}
	for _, e := range codeExts {
		if ext == e {
			return true
		}
	}
	return false
}
