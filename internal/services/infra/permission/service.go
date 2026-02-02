package permission

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/pelletier/go-toml/v2"
)

var permissions []string

type PermissionService interface {
	LoadPermissionsFromDef(fileFs embed.FS)
}

type PermissionServiceImpl struct {
	client *ent.Client
}

func NewPermissionServiceImpl(client *ent.Client) *PermissionServiceImpl {
	return &PermissionServiceImpl{client: client}
}

func (s *PermissionServiceImpl) LoadPermissionsFromDef(fileFs embed.FS) {
	dir, _ := fs.ReadDir(fileFs, "assets/moduleDefs")
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		filePath := fmt.Sprintf("assets/moduleDefs/%s", entry.Name())

		content, err := fileFs.ReadFile(filePath)
		if err != nil {
			log.Printf("读取文件 %s 失败: %v", filePath, err)
			continue
		}

		var moduleDef model.ModuleDef
		err = toml.Unmarshal(content, &moduleDef)
		if err != nil {
			log.Printf("解析文件 %s 失败: %v", filePath, err)
			continue
		}
		var scopes []string
		// 将 permission.scope复制到 scopes
		for _, permission := range moduleDef.Meta.Permissions {
			scopes = append(scopes, permission.Scope)
		}

		permissions = append(permissions, scopes...)
	}
	// log.Printf("%v", permissions)
}
