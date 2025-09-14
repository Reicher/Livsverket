# 📱 MVP Specification: Livsverket (Web)

## 🔑 Main Idea

A small **React web app** where the user can **explore the tree of life** (biological taxonomy à la Linnaeus) and choose to **collect any “set”** – any taxonomic level from *domain* down to *species*. A simple **Go backend** serves data from JSON files and stores user actions (collections & sightings) in JSON as well.

---

## 🧬 Core Features

### 🌳 Tree of Life Navigation

* Start view shows the three top-level **sets**: *Bacteria*, *Archaea*, *Eukaryotes*
* Each set is a **taxonomic group**: domain, kingdom, phylum, class, order, family, genus, species.
* The user can:

    * **Explore** – view child sets
    * **Collect** – add it to collection
    * **Report Sighting** – mark as seen
* Navigation includes:

    * Going downward via **Explore**
    * Going upward via a **breadcrumb** (title shows parent chain)

### 🎯 Interacting with a Set

When the user clicks a set in the tree view, a small menu appears:

```
🦉 Eurasian Eagle-Owl (species)
– [Explore]
– [Collect]
– [Report]
```

### 📦 Data & Storage

* **Data source**: `sets.json` with the taxonomy tree
* **User data**: `collections.json` and `sightings.json`
* Each set contains:

    * `id`
    * `nameSv`
    * `rank`
    * `parentId`
    * `childCount`
    * `description` (optional)

### 📋 My Collection

* A list of all collected sets.
* Each entry shows:

    * Name, rank
    * Stats: total children vs collected children
* Clicking an entry navigates to that set in the tree.

---

## 🧱 System Architecture

* **Frontend**: React + TypeScript + Vite, Tailwind CSS
* **Backend**: Go HTTP server that:

    * Serves JSON endpoints from `sets.json`
    * Updates and reads `collections.json` and `sightings.json`
* **API**: Simple REST JSON

---

## 🌐 REST API

### Sets

* `GET /sets/root`
* `GET /sets/{id}`
* `GET /sets/{id}/children`
* `GET /sets/{id}/breadcrumbs`

### Collections

* `GET /collections`
* `POST /collections { setId }`
* `DELETE /collections/{setId}`

### Sightings

* `GET /sightings`
* `POST /sightings { setId, timestamp?, note? }`

---

## 🖥️ Frontend UX

* **Home / Tree**: breadcrumb + child sets
* **My Collection**: list with progress
* Context menu: Explore / Collect / Report

---

## 🛠️ Tech Stack

* **Frontend**: React 18, TypeScript, Vite, Tailwind
* **Backend**: Go 1.22+, net/http, file read/write JSON

---

## 📄 README Requirements

README must describe how to:

* Clone the repo
* Install Node.js & Go
* Run backend: `go run main.go` (reads/writes JSON files)
* Run frontend: `pnpm i && pnpm dev`
* Access app in browser
* Troubleshoot (e.g., port conflicts)
* License: MIT
