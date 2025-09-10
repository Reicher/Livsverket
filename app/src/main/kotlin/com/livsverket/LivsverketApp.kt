package com.livsverket

import android.widget.Toast
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.List
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import com.livsverket.data.CollectionStore
import com.livsverket.data.TaxonomyRepository
import com.livsverket.model.Taxon

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LivsverketApp() {
    val context = LocalContext.current
    val repo = remember { TaxonomyRepository(context) }
    val collectionStore = remember { CollectionStore(context) }
    val collected = remember { mutableStateListOf<String>().apply { addAll(collectionStore.getCollection()) } }
    var currentId by remember { mutableStateOf<String?>(null) }
    var menuTaxon by remember { mutableStateOf<Taxon?>(null) }
    var showCollection by remember { mutableStateOf(false) }

    val title = repo.getTaxon(currentId)?.name ?: "Tree of Life"

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        title,
                        modifier = Modifier.clickable(enabled = currentId != null) {
                            currentId = repo.getTaxon(currentId)?.parentId
                        }
                    )
                },
                actions = {
                    IconButton(onClick = { showCollection = true }) {
                        Icon(Icons.Filled.List, contentDescription = "Collection")
                    }
                }
            )
        }
    ) { padding ->
        Box(modifier = Modifier.padding(padding)) {
            val children = repo.getChildren(currentId)
            LazyColumn {
                items(children) { taxon ->
                    ListItem(
                        headlineText = { Text(taxon.name) },
                        supportingText = { Text(taxon.level) },
                        modifier = Modifier.clickable { menuTaxon = taxon }
                    )
                }
            }
            if (menuTaxon != null) {
                TaxonMenu(
                    taxon = menuTaxon!!,
                    onDismiss = { menuTaxon = null },
                    onExplore = {
                        currentId = menuTaxon!!.id
                        menuTaxon = null
                    },
                    onCollect = {
                        if (!collected.contains(menuTaxon!!.id)) {
                            collected.add(menuTaxon!!.id)
                            collectionStore.saveCollection(collected)
                        }
                        menuTaxon = null
                    },
                    onReport = {
                        Toast.makeText(context, "Reported sighting of ${menuTaxon!!.name}", Toast.LENGTH_SHORT).show()
                        menuTaxon = null
                    }
                )
            }
            if (showCollection) {
                CollectionView(
                    repo = repo,
                    collectedIds = collected,
                    onClose = { showCollection = false },
                    onSelectTaxon = {
                        currentId = it.id
                        showCollection = false
                    }
                )
            }
        }
    }
}

@Composable
fun TaxonMenu(
    taxon: Taxon,
    onDismiss: () -> Unit,
    onExplore: () -> Unit,
    onCollect: () -> Unit,
    onReport: () -> Unit
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("${taxon.name} (${taxon.level})") },
        text = {
            Column {
                TextButton(onClick = onExplore) { Text("Explore") }
                TextButton(onClick = onCollect) { Text("Collect") }
                TextButton(onClick = onReport) { Text("Report") }
            }
        },
        confirmButton = {},
        dismissButton = {}
    )
}

@Composable
fun CollectionView(
    repo: TaxonomyRepository,
    collectedIds: List<String>,
    onClose: () -> Unit,
    onSelectTaxon: (Taxon) -> Unit
) {
    AlertDialog(
        onDismissRequest = onClose,
        title = { Text("My Collection") },
        text = {
            LazyColumn {
                items(collectedIds) { id ->
                    val taxon = repo.getTaxon(id) ?: return@items
                    val total = repo.getChildren(id).size
                    val collected = repo.getChildren(id).count { collectedIds.contains(it.id) }
                    ListItem(
                        headlineText = { Text(taxon.name) },
                        supportingText = { Text("${taxon.level} – $collected/$total") },
                        modifier = Modifier.clickable { onSelectTaxon(taxon) }
                    )
                }
            }
        },
        confirmButton = {
            TextButton(onClick = onClose) { Text("Close") }
        }
    )
}
