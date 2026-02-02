package util

import (
	"github.com/hashicorp/go-version"
)

// ValidateSemanticVersion 校验是否为合法的语义化版本号（支持v前缀，如v1.0.0）
func ValidateSemanticVersion(versionStr string) bool {
	_, err := version.NewVersion(versionStr)
	return err == nil
}

// CompareVersion 比较两个版本号大小
// return 1：v1 > v2；0：v1 == v2；-1：v1 < v2
func CompareVersion(v1, v2 string) int {
	ver1, err1 := version.NewVersion(v1)
	ver2, err2 := version.NewVersion(v2)
	if err1 != nil || err2 != nil {
		return 0
	}
	return ver1.Compare(ver2)
}
