package modpack_service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	file_service "github.com/brutalzinn/boberto-modpack-api/infra/services/file"
	manifest_models "github.com/brutalzinn/boberto-modpack-api/infra/services/manifest/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/infra/services/modpack/models"
)

func GetModPackFiles(modPack modpack_models.MinecraftModPack,
	environment modpack_models.MinecraftEnvironment) []manifest_models.ManifestFile {
	normalizedName := common.NormalizeString(modPack.Name)
	cfg := config.GetConfig()
	modpackPath := filepath.Join(
		cfg.ModPacks.PublicPath,
		normalizedName,
		environment.GetFolderName())

	modpackFiles := []manifest_models.ManifestFile{}
	err := filepath.Walk(modpackPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			relativePath := strings.ReplaceAll(path, modpackPath+string(os.PathSeparator), "")
			fileUrl := modPack.FileUrl + relativePath
			checksum, _ := file_service.GetChecksum(path)
			fileType := GetType(relativePath)
			modpackFile := manifest_models.ManifestFile{
				Name:        info.Name(),
				Size:        info.Size(),
				Path:        relativePath,
				Url:         fileUrl,
				Checksum:    checksum,
				Environment: environment,
				Type:        fileType,
			}
			modpackFiles = append(modpackFiles, modpackFile)
			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return modpackFiles
}

func GetType(file string) modpack_models.MinecraftFileType {
	parts := strings.Split(file, string(os.PathSeparator))
	switch parts[0] {
	case "mods":
		return modpack_models.Mod
	case "data":
		return modpack_models.Data
	case "saves":
		return modpack_models.World
	case "config":
		return modpack_models.Config
	case "natives":
		return modpack_models.Library
	case "shaderpacks":
		return modpack_models.Texture
	default:
		return modpack_models.Other
	}
}
