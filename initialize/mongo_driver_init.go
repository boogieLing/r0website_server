// Package initialize
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: mongo操作驱动
 * @File:  mongo_driver
 * @Version: 1.0.0
 * @Date: 2022/7/3 23:52
 */
package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"r0Website-server/config"
	"r0Website-server/dao"
	"r0Website-server/global"
	"time"
)

// InitDB 获取一个mongo的db操作实例
// 一定要记得defer关闭client
func InitDB() (context.Context, *mongo.Client) {
	// FOLLOW:
	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/#std-label-golang-authentication-mechanisms
	credential := options.Credential{
		AuthSource: global.Config.Mongo.DB,
		Username:   global.Config.Mongo.SuperAdminName,
		Password:   global.Config.Mongo.SuperAdminPswd,
	}
	dsn := fmt.Sprintf("mongodb://%s:%s/?authSource=admin",
		global.Config.Mongo.Address,
		global.Config.Mongo.Part,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn).SetAuth(credential))
	if err != nil {
		global.Logger.Error(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		global.Logger.Error(err)
	}
	// db := client.Database(global.Config.Mongo.DB)
	global.Logger.Infof("=== Connected MongoDB:%s ===", global.Config.Mongo.DB)
	return ctx, client
}

// MongoConstructor Mongo构造
func MongoConstructor(cfg *config.Mongo) *dao.BasicDaoMongo {
	credential := options.Credential{
		AuthSource: cfg.DB,
		Username:   cfg.SuperAdminName,
		Password:   cfg.SuperAdminPswd,
	}
	dsn := fmt.Sprintf("mongodb://%s:%s/?authSource=admin",
		cfg.Address,
		cfg.Part,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, connectErr := mongo.Connect(ctx, options.Client().ApplyURI(dsn).SetAuth(credential))
	var errs []error
	if connectErr != nil {
		global.Logger.Error(connectErr)
		errs = append(errs, connectErr)
	}
	if pingErr := client.Ping(ctx, readpref.Primary()); pingErr != nil {
		global.Logger.Error(pingErr)
		errs = append(errs, pingErr)
	}
	// db := client.Database(cfg.DB)
	global.Logger.Infof("=== Connected MongoDB:%s ===", cfg.DB)
	if len(errs) == 0 {
		global.Logger.Infof("=== Connected MongoDB:%s ===", cfg.DB)
		fmt.Printf("=== Connected MongoDB:%s ===\n", cfg.DB)
	} else {
		fmt.Println(errs)
	}
	db := client.Database(cfg.DB)
	// ✅ 初始化图集项目的索引（幂等，极速）
	if err := InitPicMongoIndexes(db); err != nil {
		fmt.Printf("⚠ InitPicMongoIndexes failed: %v\n", err)
	} else {
		fmt.Printf("✅ Pic-related Mongo indexes ensured\n")
	}
	return &dao.BasicDaoMongo{Mc: client, Mdb: client.Database(cfg.DB)}
}

// PicIndexSpec 定义一个图集项目所需索引的结构
type PicIndexSpec struct {
	CollectionName string
	Model          mongo.IndexModel
}

// InitPicMongoIndexes 初始化图集项目相关的索引，幂等并极致压缩时间
func InitPicMongoIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. 定义所有需要的索引
	requiredIndexes := []PicIndexSpec{
		// images 索引
		{"images", mongo.IndexModel{
			Keys:    bson.D{{Key: "cos_url", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_cos_url_unique"),
		}},
		{"images", mongo.IndexModel{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("idx_tags"),
		}},
		{"images", mongo.IndexModel{
			Keys:    bson.D{{Key: "uploaded_at", Value: -1}},
			Options: options.Index().SetName("idx_uploaded_at_desc"),
		}},
		// albums 索引
		{"albums", mongo.IndexModel{
			Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
			Options: options.Index().SetName("idx_text_title_description"),
		}},
		{"albums", mongo.IndexModel{
			Keys:    bson.D{{Key: "author", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_author_createdat"),
		}},
		{"albums", mongo.IndexModel{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("idx_album_tags"),
		}},
		{"albums", mongo.IndexModel{
			Keys:    bson.D{{Key: "visibility", Value: 1}},
			Options: options.Index().SetName("idx_visibility"),
		}},
		{"albums", mongo.IndexModel{
			Keys:    bson.D{{Key: "cover_image", Value: 1}},
			Options: options.Index().SetName("idx_cover_image"),
		}},
		{"albums", mongo.IndexModel{
			Keys:    bson.D{{Key: "image_refs.image_id", Value: 1}},
			Options: options.Index().SetName("idx_image_refs_image_id"),
		}},
	}

	// 2. 扫描所有已存在索引（每个集合只查一次）
	existingIndexNames := make(map[string]map[string]bool)

	for _, spec := range requiredIndexes {
		coll := spec.CollectionName

		if _, scanned := existingIndexNames[coll]; !scanned {
			cursor, err := db.Collection(coll).Indexes().List(ctx)
			if err != nil {
				return fmt.Errorf("failed to list indexes on %s: %w", coll, err)
			}
			indexNameSet := map[string]bool{}
			for cursor.Next(ctx) {
				var entry bson.M
				if err := cursor.Decode(&entry); err == nil {
					if name, ok := entry["name"].(string); ok {
						indexNameSet[name] = true
					}
				}
			}
			existingIndexNames[coll] = indexNameSet
		}
	}

	// 3. 判断缺失并创建
	for _, spec := range requiredIndexes {
		collName := spec.CollectionName
		indexName := ""
		if spec.Model.Options != nil && spec.Model.Options.Name != nil {
			indexName = *spec.Model.Options.Name
		}
		if indexName == "" {
			continue // 无名索引不建议继续执行（默认略过）
		}
		if !existingIndexNames[collName][indexName] {
			_, err := db.Collection(collName).Indexes().CreateOne(ctx, spec.Model)
			if err != nil {
				global.Logger.Warnf("❌ Failed to create index %s on %s: %v", indexName, collName, err)
			} else {
				global.Logger.Infof("✅ Created index %s on %s", indexName, collName)
			}
		} else {
			global.Logger.Infof("⏩ Skipped existing index %s on %s", indexName, collName)
		}
	}

	return nil
}
