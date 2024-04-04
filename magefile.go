//go:build mage
// +build mage

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/joho/godotenv/autoload"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	cwd, _           = os.Getwd()
	binDir           = filepath.Join(cwd, "_bin")
	binReleasesDir   = filepath.Join(binDir, "releases")
	releasesDir      = filepath.Join(cwd, "releases")
	releaseUrlPrefix = "https://github.com/BossRighteous/MiSTer_Games_GUI/releases/latest/download"
	docsDir          = filepath.Join(cwd, "docs")
	upxBin           = os.Getenv("UPX_BIN")
	// docker arm build
	armBuild          = filepath.Join(cwd, "scripts", "armbuild")
	armBuildImageName = "mistergamesgui/armbuild"
	armBuildCache     = filepath.Join(os.TempDir(), "mistergamesgui-buildcache")
	armModCache       = filepath.Join(os.TempDir(), "mistergamesgui-modcache")
	// docker kernel build
	kernelBuild          = filepath.Join(cwd, "scripts", "kernelbuild")
	kernelBuildImageName = "mistergamesgui/kernelbuild"
	kernelRepoName       = "Linux-Kernel_MiSTer"
	kernelRepoPath       = filepath.Join(kernelBuild, "_build", kernelRepoName)
	kernelRepoUrl        = fmt.Sprintf("https://github.com/MiSTer-devel/%s.git", kernelRepoName)
)

type app struct {
	name         string
	path         string
	bin          string
	ldFlags      string
	releaseId    string
	reboot       bool
	inAll        bool
	releaseFiles []string
}

var apps = []app{
	{
		name: "mistergamesgui",
		path: filepath.Join(cwd, "cmd", "mistergamesgui"),
		bin:  "mistergamesgui",
	},
}

func getApp(name string) *app {
	for _, a := range apps {
		if a.name == name {
			return &a
		}
	}
	return nil
}

func cleanPlatform(name string) {
	_ = sh.Rm(filepath.Join(binDir, name))
}

func Clean() {
	_ = sh.Rm(binDir)
	_ = sh.Rm(armBuildCache)
	_ = sh.Rm(armModCache)
	_ = sh.Rm(kernelRepoPath)
}

func buildApp(a app, out string) {
	if a.ldFlags == "" {
		env := map[string]string{
			"GOPROXY": "https://goproxy.io,direct",
		}
		_ = sh.RunWithV(env, "go", "build", "-o", out, a.path)
	} else {
		staticEnv := map[string]string{
			"GOPROXY":     "https://goproxy.io,direct",
			"CGO_ENABLED": "1",
			"CGO_LDFLAGS": a.ldFlags,
		}
		_ = sh.RunWithV(staticEnv, "go", "build", "--ldflags", "-linkmode external -extldflags -static", "-o", out, a.path)
	}
}

func Build(appName string) {
	platform := runtime.GOOS + "_" + runtime.GOARCH
	if appName == "all" {
		mg.Deps(func() { cleanPlatform(platform) })
		for _, app := range apps {
			fmt.Println("Building", app.name)
			buildApp(app, filepath.Join(binDir, platform, app.bin))
		}
	} else {
		app := getApp(appName)
		if app == nil {
			fmt.Println("Unknown app", appName)
			os.Exit(1)
		}
		buildApp(*app, filepath.Join(binDir, platform, app.bin))
	}
}

func MakeArmImage() {
	if runtime.GOOS != "linux" {
		_ = sh.RunV("docker", "build", "--platform", "linux/arm/v7", "-t", armBuildImageName, armBuild)
	} else {
		_ = sh.RunV("sudo", "docker", "build", "--platform", "linux/arm/v7", "-t", armBuildImageName, armBuild)
	}
}

func Mister(appName string) {
	buildCache := fmt.Sprintf("%s:%s", armBuildCache, "/home/build/.cache/go-build")
	_ = os.Mkdir(armBuildCache, 0755)
	modCache := fmt.Sprintf("%s:%s", armModCache, "/home/build/go/pkg/mod")
	_ = os.Mkdir(armModCache, 0755)
	buildDir := fmt.Sprintf("%s:%s", cwd, "/build")
	if runtime.GOOS != "linux" {
		_ = sh.RunV("docker", "run", "--rm", "--platform", "linux/arm/v7", "-v", buildCache, "-v", modCache, "-v", buildDir, "--user", "1000:1000", armBuildImageName, "mage", "build", appName)
	} else {
		_ = sh.RunV("sudo", "docker", "run", "--rm", "--platform", "linux/arm/v7", "-v", buildCache, "-v", modCache, "-v", buildDir, "--user", "1000:1000", armBuildImageName, "mage", "build", appName)
	}
}

func getMd5Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	_ = file.Close()
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func getFileSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	stat, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return 0, err
	}

	size := stat.Size()
	_ = file.Close()

	return size, nil
}

func Release(name string) {
	a := getApp(name)
	if a == nil {
		fmt.Println("Unknown app", name)
		os.Exit(1)
	}

	Mister(name)

	rd := filepath.Join(releasesDir, a.name)
	_ = os.MkdirAll(rd, 0755)
	_ = os.MkdirAll(binReleasesDir, 0755)
	releaseBin := filepath.Join(binReleasesDir, a.bin)
	err := sh.Copy(releaseBin, filepath.Join(binDir, "linux_arm", a.bin))
	if err != nil {
		fmt.Println("Error copying binary", err)
		os.Exit(1)
	}

	for _, f := range a.releaseFiles {
		err := sh.Copy(filepath.Join(binReleasesDir, filepath.Base(f)), f)
		if err != nil {
			fmt.Println("Error copying release file", err)
			os.Exit(1)
		}
	}

	if upxBin == "" {
		fmt.Println("UPX is required for releases")
		os.Exit(1)
	} else {
		if runtime.GOOS != "windows" {
			err := os.Chmod(releaseBin, 0755)
			if err != nil {
				fmt.Println("Error chmod release bin", err)
				os.Exit(1)
			}
		}

		err := sh.RunV(upxBin, "-9", releaseBin)
		if err != nil {
			fmt.Println("Error compressing binary", err)
			os.Exit(1)
		}
	}
}

func PrepRelease() {
	_ = sh.Rm(binReleasesDir)
	_ = os.MkdirAll(binReleasesDir, 0755)
	cleanPlatform("linux_arm")
	for _, app := range apps {
		if app.releaseId != "" {
			fmt.Println("Preparing release:", app.name)
			Release(app.name)
		}
	}
}

func MakeKernelImage() {
	_ = sh.RunV("sudo", "docker", "build", "-t", kernelBuildImageName, kernelBuild)
}

func Kernel() {
	if _, err := os.Stat(kernelRepoPath); os.IsNotExist(err) {
		_ = sh.RunV("git", "clone", "--depth", "1", kernelRepoUrl, kernelRepoPath)
	}

	patches, _ := filepath.Glob(filepath.Join(kernelBuild, "*.patch"))
	for _, path := range patches {
		_ = sh.RunV("git", "-C", kernelRepoPath, "apply", path)
	}

	kCmd := sh.RunCmd("sudo", "docker", "run", "--rm", "-v", fmt.Sprintf("%s:%s", kernelRepoPath, "/build"), "--user", "1000:1000", kernelBuildImageName)
	_ = kCmd("make", "MiSTer_defconfig")
	_ = kCmd("make", "modules")
	_ = kCmd("make", "-j16", "zImage")
	_ = kCmd("make", "socfpga_cyclone5_de10_nano.dtb")

	zImage, _ := os.Open(filepath.Join(kernelRepoPath, "arch", "arm", "boot", "zImage"))
	dtb, _ := os.Open(filepath.Join(kernelRepoPath, "arch", "arm", "boot", "dts", "socfpga_cyclone5_de10_nano.dtb"))

	_ = os.MkdirAll(filepath.Join(binDir, "linux"), 0755)
	kernel, _ := os.Create(filepath.Join(binDir, "linux", "zImage_dtb"))

	_, _ = io.Copy(kernel, zImage)
	_, _ = io.Copy(kernel, dtb)

	_ = kernel.Close()
	_ = dtb.Close()
	_ = zImage.Close()
}

func MakeArmApp(name string) {
	buildScript := name + ".sh"
	if _, err := os.Stat(filepath.Join(armBuild, buildScript)); os.IsNotExist(err) {
		fmt.Println("No build script for", name)
		os.Exit(1)
	}

	buildDir := filepath.Join(armBuild, "_build")
	_ = os.MkdirAll(buildDir, 0755)

	err := sh.Copy(filepath.Join(buildDir, buildScript), filepath.Join(armBuild, buildScript))
	if err != nil {
		fmt.Println("Error copying build script", err)
		os.Exit(1)
	}

	if runtime.GOOS != "linux" {
		_ = sh.RunV("docker", "run", "--rm", "--platform", "linux/arm/v7", "-v", buildDir+":/build", "--user", "1000:1000", armBuildImageName, "bash", "./"+buildScript)
	} else {
		_ = sh.RunV("sudo", "docker", "run", "--rm", "--platform", "linux/arm/v7", "-v", buildDir+":/build", "--user", "1000:1000", armBuildImageName, "bash", "./"+buildScript)
	}
}
