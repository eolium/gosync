# GoSync


Gosync est un programme écrit en go, donc disponible pour Linux, MacOS et Windows, qui permet de copier un dossier, puis de mettre à jour seulement les fichiers modifiés. Exemple d'utilisation :

## Linux

```bash
./gosync source backup 
```

```bash
./gosync source1 source2 source3 backup
```

## Windows

```cmd
gosync.exe source backup
```

```bat
gosync.exe source1 source2 source3 backup
```

Un usage conseillé est de sauvegarder l'exécutable dans un dossier quelconque, puis de créer un script batch pour exécuter la commande et éviter de reécrire le même chemin.