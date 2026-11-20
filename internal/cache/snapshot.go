package cache

import "github.com/attmonk/attmonk-plugin-dispatcharr/internal/model"

func SnapshotFromCatalog(catalog model.CatalogState) Snapshot {
	return Snapshot{Catalog: catalog, Health: catalog.Health}
}
