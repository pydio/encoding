package main

import (
	"github.com/mkideal/cli"
	"github.com/spf13/afero"
	"errors"
	"os"
	"path/filepath"
	"golang.org/x/text/unicode/norm"
	"bytes"
	"runtime"
)

type argT struct {
	cli.Helper
	Folder 	string `cli:"*f,folder" usage:"Root folder to start scanning for NFD encoded filenames"`
	DryRun  bool `cli:"d,dry-run" usage:"Use dry run to see all changes without applying them"`
}

func main() {

	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		fs := afero.NewBasePathFs(afero.NewOsFs(), argv.Folder)
		if runtime.GOOS == "darwin" && !argv.DryRun {
			return errors.New("Cannot convert filenames to NFC on MacOSX! Use -d to execute in dry-run mode.")
		}
		if exists, _ := afero.Exists(fs, "/"); !exists {
			return errors.New("Cannot find root path " + argv.Folder + "!")
		}
		ctx.String("Running Scan on folder=%s, Dry-run =%v\n", argv.Folder, argv.DryRun)
		ScanFolder(ctx, fs, argv.DryRun)
		return nil
	})
}

func ScanFolder(ctx *cli.Context, fs afero.Fs, dryRun bool) {

	toRename := make(map[string]string)
	var keys []string

	afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error{

		//ctx.String(path + "\n")
		dir, base := filepath.Split(path)
		stringBytes := []byte(base)
		normed := norm.NFC.Bytes(stringBytes)
		if bytes.Compare(normed, stringBytes) != 0 {
			fileType := "file"
			if info.IsDir() {
				fileType = "folder"
			}
			// Parent folder will already have been renamed at that point
			oldName := filepath.Join(string(norm.NFC.Bytes([]byte(dir))), base)
			renamed := string(norm.NFC.Bytes([]byte(path)))

			// Create alternative version
			if dryRun {
				ctx.String("=> Should update "+ fileType +" name " + path + " to " + renamed +"\n")
			}
			toRename[oldName] = renamed
			keys = append(keys, oldName)
		}

		return nil
	});

	if len(toRename) > 0 {
		ctx.String("=============================\n")
		ctx.String("   FOUND %d ITEMS TO RENAME  \n", len(toRename))
		ctx.String("=============================\n")
		if !dryRun {
			for _, oldName := range keys {

				newName := toRename[oldName]
				ctx.String("+ Renaming " + oldName + " to " + newName +"\n")
				e := fs.Rename(oldName, newName)
				if e != nil {
					ctx.String("  => ERROR while renaming:" + e.Error())
				}

			}
			ctx.String("====================================================\n")
			ctx.String(" Please make sure to now relaunch an indexation for \n" +
				"each workspaces that may have been modified \n")
			ctx.String("====================================================\n")
		}

	}else {
		ctx.String("===================================\n")
		ctx.String("   NO NFD ENCODED FILENAMES FOUND  \n")
		ctx.String("===================================\n")
	}


}