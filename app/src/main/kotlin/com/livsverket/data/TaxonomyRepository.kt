package com.livsverket.data

import android.content.Context
import com.livsverket.model.Taxon
import org.json.JSONArray

class TaxonomyRepository(private val context: Context) {
    private val taxons: List<Taxon> by lazy { loadTaxonomy() }

    private fun loadTaxonomy(): List<Taxon> {
        val jsonStr = context.assets.open("taxonomy.json").bufferedReader().use { it.readText() }
        val jsonArray = JSONArray(jsonStr)
        val list = mutableListOf<Taxon>()
        for (i in 0 until jsonArray.length()) {
            val obj = jsonArray.getJSONObject(i)
            list.add(
                Taxon(
                    id = obj.getString("id"),
                    name = obj.getString("name"),
                    level = obj.getString("level"),
                    parentId = if (obj.isNull("parentId")) null else obj.getString("parentId"),
                    childCount = obj.optInt("childCount", 0),
                    description = if (obj.has("description")) obj.getString("description") else null
                )
            )
        }
        return list
    }

    fun getChildren(parentId: String?): List<Taxon> =
        taxons.filter { it.parentId == parentId }

    fun getTaxon(id: String?): Taxon? =
        taxons.find { it.id == id }
}
