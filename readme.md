# GoSync

GoSync est un programme écrit en Go, disponible pour Linux et Windows, qui permet de copier un dossier, puis de mettre à jour seulement les fichiers modifiés.

## Compilation

Pour compiler GoSync à partir des sources, suivez les étapes ci-dessous :

1. Assurez-vous d'avoir Go installé sur votre machine. Vous pouvez télécharger Go depuis [le site officiel](https://golang.org/dl/).
2. Clonez le dépôt GoSync :
   ```bash
   git clone https://github.com/eolium/gosync.git
   ```
3. Accédez au répertoire du projet :
   ```bash
   cd gosync
   ```
4. Compilez le programme :
   ```bash
   go build main.go
   ```

## Utilisation

### Linux

```bash
./main <source> <backup>
```

```bash
./main <source1> <source2> <source3> <backup>
```

### Windows


```cmd
main.exe <source> <backup>
```

```bat
main.exe <source1> <source2> <source3> <backup>
```


## Options
| Option          | Description                                                                 |
|-----------------|-----------------------------------------------------------------------------|
| `-v, --verbose` | Affiche des informations détaillées pendant l'exécution.                    |

## Conseils d'utilisation

Il est conseillé de sauvegarder l'exécutable dans un dossier quelconque, puis de créer un script batch pour exécuter la commande et éviter de réécrire le même chemin.
