package util

import (
	"crypto/sha256"
	"encoding/hex"
	"lightOA-end/src/entity"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 自定义token生成方法
func FormToken(username string) string {
	u := uuid.NewV4().String()
	hash := sha256.New()
	hash.Write([]byte(u + time.Now().Format(time.RFC1123) + username))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

func FormUserRole(role *entity.Role, resources []entity.ResourceRaw) *entity.Role {
	nodeMap := make(map[int]*entity.Resource)
	var roots []*entity.Resource
	// 将节点按照 ParentID 存储到 map 中
	for _, node := range resources {
		nodeMap[node.Id] = &entity.Resource{
			Id:       node.Id,
			Alias:    node.Alias,
			Name:     node.Name,
			Type:     node.Type,
			ParentId: node.ParentId,
		}
	}
	// 构建树结构
	for _, node := range resources {
		parentNode, parentOk := nodeMap[node.ParentId]
		node, ok := nodeMap[node.Id]
		if parentOk {
			parentNode.Children = append(parentNode.Children, node)
		} else if ok {
			roots = append(roots, node)
		}
	}
	role.Resources = roots
	return role
}

func FormResources(resources []entity.ResourceRaw) *entity.Resource {
	nodeMap := make(map[int]*entity.Resource)
	var roots []*entity.Resource
	// 将节点按照 ParentID 存储到 map 中
	for _, node := range resources {
		nodeMap[node.Id] = &entity.Resource{
			Id:       node.Id,
			Alias:    node.Alias,
			Name:     node.Name,
			Type:     node.Type,
			ParentId: node.ParentId,
		}
	}
	// 构建树结构
	for _, node := range resources {
		parentNode, parentOk := nodeMap[node.ParentId]
		node, ok := nodeMap[node.Id]
		if parentOk {
			parentNode.Children = append(parentNode.Children, node)
		} else if ok {
			roots = append(roots, node)
		}
	}
	return roots[0]
}

func Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}
