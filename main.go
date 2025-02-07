package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag" // pour les options dans la commande
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Variable globale pour activer/désactiver le mode verbose
var verbose bool
var quiet bool

// Affiche les messages seulement si le mode verbose est activé
func logVerbose(format string, args ...interface{}) {
	if !quiet && verbose {
		log.Printf(format, args...)
	}
}

// calcule de la signature SHA256 d'un fichier donné et retourne sa valeur hexadécimale.
func filehash(filename string) (string, error) {
	logVerbose("Scan du fichier %s...", filename)

	hasher := sha256.New()
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("échec d'ouverture du fichier : %w", err)
	}
	defer file.Close()

	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", fmt.Errorf("échec de calcul du hash : %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Copie le fichier source vers sa destination.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("échec d'ouverture du fichier source %s : %w", src, err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("échec de création du fichier destination %s : %w", dst, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("échec de copie de %s vers %s : %w", src, dst, err)
	}

	logVerbose("Fichier %s copié vers %s\n", src, dst)
	return nil
}

// synchronise la structure des dossiers de src vers dst.
func syncFolder(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erreur lors de la navigation dans %s : %w", path, err)
		}

		rpath, err := filepath.Rel(src, path)
		if err != nil || rpath == "." {
			return err
		}

		dstPath := filepath.Join(dst, rpath)

		if info.IsDir() {
			if _, err := os.Stat(dstPath); os.IsNotExist(err) {
				err := os.MkdirAll(dstPath, info.Mode())
				if err != nil {
					return fmt.Errorf("échec de création du dossier %s : %w", dstPath, err)
				}
				logVerbose("Dossier %s synchronisé\n", rpath)
			}
		}

		return nil
	})
}

// synchronise les fichiers de src vers dst en vérifiant les différences de contenu.
func sync(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erreur lors de la navigation dans %s : %w", path, err)
		}

		rpath, err := filepath.Rel(src, path)
		if err != nil || rpath == "." {
			return err
		}

		dstPath := filepath.Join(dst, rpath)

		if !info.IsDir() {
			if _, err := os.Stat(dstPath); os.IsNotExist(err) {
				err := copyFile(path, dstPath)
				if err != nil {
					return err
				}
				logVerbose("Fichier %s copié\n", rpath)
			} else {
				srcHash, err := filehash(path)
				if err != nil {
					return err
				}
				dstHash, err := filehash(dstPath)
				if err != nil {
					return err
				}

				if srcHash != dstHash {
					err := copyFile(path, dstPath)
					if err != nil {
						return err
					}
					logVerbose("Fichier %s mis à jour\n", rpath)
				}
			}
		}

		return nil
	})
}

// Classic Entry Point.
func main() {
	// Définition de l'option verbose
	flag.BoolVar(&verbose, "verbose", false, "Afficher des détails sur les fichiers transférés")
	flag.BoolVar(&verbose, "v", false, "Alias pour --verbose")
	flag.BoolVar(&quiet, "quiet", false, "Ne rien afficher")
	flag.BoolVar(&quiet, "q", false, "Alias pour --quiet")
	flag.Parse()

	// Vérifie que suffisamment d'arguments sont fournis
	if flag.NArg() < 2 {
		fmt.Println("Usage : ./gosync [options] <SRC1> [SRC2] ... <DST>")
		fmt.Println("Options :")
		fmt.Println("  --verbose, -v    Afficher des détails sur les fichiers transférés")
		fmt.Println("  --quiet, -q      Mode silencieux")
		fmt.Println("SRC : dossier(s) à synchroniser")
		fmt.Println("DST : dossier de destination")
		return
	}

	dst := flag.Arg(flag.NArg() - 1)

	for i := 0; i < flag.NArg()-1; i++ {
		src := flag.Arg(i)

		if !quiet {
			log.Printf("Synchronisation de %s vers %s\n", src, dst)
		}

		err := syncFolder(src, dst)
		if err != nil && !quiet {
			log.Fatalf("échec de synchronisation des dossiers : %v\n", err)
		}
		err = sync(src, dst)
		if err != nil && !quiet {
			log.Fatalf("échec de synchronisation des fichiers : %v\n", err)
		}
	}

	if !quiet {
		log.Println("Synchronisation terminée avec succès.")
	}
}
