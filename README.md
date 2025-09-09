# Go Ip Country Resolver

## Objectif du package

1. **Stocker** les ranges IP (IPv4 et IPv6) de chaque pays dans une base **bbolt**.
2. **Géolocaliser rapidement** une IP vers son pays (code ISO).
3. **Lister** toutes les ranges d’un pays.
4. **Mettre à jour** automatiquement la DB en comparant des checksums locaux/distants.
5. **Accélérer les recherches** grâce à un **cache interne** et un **index optimisé**.

---

## Structure des données

* Chaque pays → représenté par un **code ISO** (`FR`, `US`, `DE`, etc.).
* Les ranges IP sont stockés **par intervalles** (`Start`, `End`).
* Les fichiers `.zone` (un par pays) listent des ranges en **CIDR** ou **IP range**.

---

## Organisation de la base bbolt

3 buckets principaux :

1. **`ipv4`** → Ranges IPv4, triés par adresse de début (`Start`).
2. **`ipv6`** → Ranges IPv6, triés par `Start`.
3. **`metadata`** → Checksum IPv4/IPv6, date de mise à jour, version…

---

## Optimisation de la recherche (`LookupIP`)

### Problème

* Si les ranges ne sont pas indexés, la recherche serait **linéaire O(n)** → trop lente avec des millions d’entrées.

### Solution (Index par `Start`)

1. Conversion IP → entier (32 bits pour IPv4, 128 bits pour IPv6).
2. Tri des ranges par `Start`.
3. **Recherche dichotomique** (binary search) :

    * Trouver le dernier range dont `Start <= IP`.
    * Vérifier si `IP <= End`.
    * Si oui → IP appartient à ce range.

### Implémentation possible

* **En mémoire** :
  Charger en RAM les ranges triés → recherche avec `sort.Search`.
* **Directement en bbolt** :
  Stocker la clé du bucket = `Start` (valeur entière).
  Utiliser `Cursor.Seek(IP)` pour trouver le range immédiatement inférieur ou égal → complexité logarithmique.

✅ Résultat : Recherche en **O(log n)** au lieu de O(n).

---

## Cache IP

* Stocke les résultats (IP → pays) des requêtes les plus fréquentes.
* Taille **configurable**, avec **limite max** (stratégie LRU).
* Évite des accès répétés à la DB pour les IP populaires.
* Cache **non persistant**, reconstruit à chaque redémarrage.

---

## Processus de mise à jour

1. Comparer **checksums distants** et **locaux** (IPv4 & IPv6).
2. Si différents :

    * Télécharger, déchiffrer et décompresser les archives dans un dossier temporaire.
    * Copier la DB actuelle en sauvegarde temporaire.
    * Vider la DB principale (sauf `metadata`).
    * Charger les nouveaux ranges (triés par `Start`).
    * Tester cohérence (lookup aléatoires, comptage).
    * Mettre à jour les métadata sur la nouvelle DB.
    * Si OK → supprimer la sauvegarde.
    * Sinon → rollback depuis la sauvegarde.
3. Le tout se fait dans un **dossier temporaire** (`os.TempDir()`).

---

## Organisation du code (fichiers proposés)

* **`resolver.go`** → Interface publique (API).
* **`db.go`** → Gestion bbolt (open, read, write, buckets).
* **`iprange.go`** → Structures + conversion + comparaison IP/range.
* **`search.go`** → Recherche optimisée (`binary search`, `Cursor.Seek`).
* **`loader.go`** → Importation depuis `.zone`.
* **`updater.go`** → Vérification checksum, mise à jour, rollback.
* **`cache.go`** → Cache mémoire LRU pour IP → pays.
* **`utils.go`** → Outils divers (checksum, IP→entier, fichiers temporaires).
