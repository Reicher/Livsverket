package com.livsverket.data

import android.content.Context

class CollectionStore(context: Context) {
    private val prefs = context.getSharedPreferences("collection", Context.MODE_PRIVATE)

    fun getCollection(): MutableSet<String> =
        prefs.getStringSet("collected", emptySet())?.toMutableSet() ?: mutableSetOf()

    fun saveCollection(ids: Collection<String>) {
        prefs.edit().putStringSet("collected", ids.toSet()).apply()
    }
}
