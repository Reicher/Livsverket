package com.livsverket.model

data class Taxon(
    val id: String,
    val name: String,
    val level: String,
    val parentId: String?,
    val childCount: Int,
    val description: String? = null
)
