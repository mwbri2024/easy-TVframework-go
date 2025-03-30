package update

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"easy-itv/config"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var logger = log.Default()

type VersionInfo struct {
	Version string            `json:"version"`
	Hash    map[string]string `json:"hash"`
}

// 获取最新版本信息
func GetLatestVersionInfo() (*VersionInfo, error) {
	resp, err := http.Get(config.VersionCheckURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info VersionInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// 获取当前平台的二进制文件名称
func GetPlatformBinaryName(version string) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// 构建平台特定的二进制文件名称
	return fmt.Sprintf("itv_%s_%s", goos, goarch)
}

// 下载新版本的二进制文件
func DownloadNewBinary(version string) (string, error) {

	// 根据当前平台获取正确的二进制文件名
	platformBinary := GetPlatformBinaryName(version)

	// 构建下载路径
	url := config.BinaryDownloadBaseURL + platformBinary + ".tar.gz"

	tempPath := filepath.Join(os.TempDir(), platformBinary)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(tempPath + ".tar.gz")
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return tempPath, err
}

// 校验下载的文件哈希值是否匹配 (使用 MD5 校验)
func VerifyFileHash(filePath, expectedHash string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 MD5 计算文件哈希
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	// 计算得到的哈希值
	computedHash := hex.EncodeToString(hash.Sum(nil))
	if computedHash != expectedHash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, computedHash)
	}
	return nil
}

// 删除临时文件
func CleanUpTempFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除临时文件失败: %v", err)
	}
	return nil
}

// 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// 解压 tar.gz 文件
func ExtractTarGz(tarGzPath string) (string, error) {
	// 获取 tar.gz 所在目录
	destDir := filepath.Dir(tarGzPath)

	file, err := os.Open(tarGzPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 解压 gzip
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	// 读取 tar
	tr := tar.NewReader(gzr)
	var extractedFilePath string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// 仅解压 itv_linux_* 文件
		baseName := filepath.Base(header.Name)
		if header.Typeflag == tar.TypeReg && strings.HasPrefix(baseName, "itv_linux_") {
			extractedFilePath = filepath.Join(destDir, baseName)
			outFile, err := os.Create(extractedFilePath)
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			if err != nil {
				return "", err
			}

			break
		}
	}

	if extractedFilePath == "" {
		return "", fmt.Errorf("未在 tar.gz 中找到 itv_linux_* 文件")
	}

	return extractedFilePath, nil
}

// 替换当前运行的二进制文件并重启
func ReplaceAndRestart(newBinaryPath, version, expectedHash string) error {
	// 解压 tar.gz
	newBinaryPath, err := ExtractTarGz(newBinaryPath + ".tar.gz")
	if err != nil {
		return fmt.Errorf("解压失败: %v", err)
	}

	// 校验哈希值，确保文件完整性
	if err := VerifyFileHash(newBinaryPath, expectedHash); err != nil {
		return err
	}

	// 获取当前二进制文件路径
	currentBinary, err := os.Executable()
	if err != nil {
		return err
	}

	// 备份当前版本
	backupPath := filepath.Join(filepath.Dir(currentBinary), fmt.Sprintf("itv_%s", version))

	// 尝试重命名当前文件作为备份
	if err := os.Rename(currentBinary, backupPath); err != nil {
		return fmt.Errorf("重命名当前文件失败: %v", err)
	}

	// 尝试将新版本复制到当前二进制文件
	if err := CopyFile(newBinaryPath, currentBinary); err != nil {
		// 回滚操作：如果替换失败，则恢复备份
		os.Rename(backupPath, currentBinary)
		return fmt.Errorf("替换新版本失败，回滚操作: %v", err)
	}

	// 尝试设置新文件的权限
	if err := os.Chmod(currentBinary, 0755); err != nil {
		// 回滚操作：如果权限设置失败，恢复备份
		os.Rename(backupPath, currentBinary)
		return fmt.Errorf("设置权限失败，回滚操作: %v", err)
	}

	// 删除临时文件
	if err := CleanUpTempFile(newBinaryPath); err != nil {
		return err
	}

	// 启动新版本 使用 syscall.Exec 替换当前进程，避免父进程退出
	if config.Debug {
		logger.Printf("更新成功 重启服务中...")
	}

	args := append([]string{currentBinary}, os.Args[1:]...)
	err = syscall.Exec(currentBinary, args, os.Environ())
	if err != nil {
		// 回滚操作：如果启动失败，恢复备份
		os.Rename(backupPath, currentBinary)
		return fmt.Errorf("启动新版本失败，回滚操作: %v", err)
	}

	return nil
}
