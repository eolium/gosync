package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

/*
Retourne le hash sha256 d'un fichier, fourni par son nom
*/
func filehash(filename string) (string, error) {
	// On crée un buffer sha256, que l'on remplie avec la lecture du fichier
	hasher := sha256.New()
	s, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	_, err = hasher.Write(s)
	if err != nil {
		return "", err
	}

	// On renvoie la somme, décodée depuis l'hexadécimal
	h := hex.EncodeToString(hasher.Sum(nil))

	return h, nil
}

/*
Synchronise la création des dossiers de src vers dst
*/
func syncFolder(src string, dst string) error {
	// On exécute pour chaque fichier/dossier récursivement
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// On récupère le chemin
		rpath, err := filepath.Rel(src, path)

		// Esquive les erreurs et la "recréation" dossier parent
		if err != nil || rpath == "." {
			return err
		}

		dst_path := filepath.Join(dst, rpath)

		if info.IsDir() {
			if _, err = os.Stat(dst_path); err != nil {
				// Si le dossier n'existe pas, on le crée

				err = os.MkdirAll(filepath.Join(dst, rpath), info.Mode())
				if err != nil {
					return err
				}

				fmt.Println(rpath, "mis à jour")
			}
		}

		return err
	})

	return err
}

func sync(src string, dst string) error {
	// On exécute pour chaque fichier/dossier récursivement
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Pas d'erreurs précédentes

		// On récupère le chemin
		rpath, err := filepath.Rel(src, path)

		// Esquive les erreurs et la "recréation" dossier parent
		if err != nil || rpath == "." {
			return err
		}

		dst_path := filepath.Join(dst, rpath)

		if !info.IsDir() {
			if _, err = os.Stat(dst_path); err != nil {
				// Le fichier n'existe pas, on le copie
				srcFile, _ := os.Open(path)
				destFile, _ := os.Create(filepath.Join(dst, rpath))

				_, err = io.Copy(destFile, srcFile)
				if err != nil {
					return err
				}

				fmt.Println(rpath, "copié")
			} else {
				// Le fichier existe
				a, _ := filehash(path)
				b, _ := filehash(dst_path)

				if a != b {
					// Les hash sont différents, On maj le fichier
					srcFile, _ := os.Open(path)
					destFile, _ := os.Create(filepath.Join(dst, rpath))

					_, err = io.Copy(destFile, srcFile)
					if err != nil {
						return err
					}

					fmt.Println(rpath, "mis à jour")
				}
			}
		}

		return err
	})

	return err
}

func main() {
	if len(os.Args) < 3 {
		// Pas assez d'arguments (compter ./gosync, SRC, DST)
		fmt.Println(
			"Usage : ./gosync [SRC1] [SRC2] ... [DST]\n",
			"SRC : dossier à copier\n",
			"DST : dossier destination",
		)
		return
	}

	dst := os.Args[len(os.Args)-1]

	for i := 1; i < len(os.Args)-1; i++ {
		src := os.Args[i]

		err := syncFolder(src, dst)
		if err != nil {
			panic(err)
		}
		err = sync(src, dst)
		if err != nil {
			panic(err)
		}
	}
}
