import { useEffect, useState } from 'react'

interface SetItem {
  id: string
  nameSv: string
  rank: string
  parentId: string
  childCount: number
}

interface Sighting {
  setId: string
  timestamp?: string
  note?: string
}

export default function App() {
  const [currentId, setCurrentId] = useState('root')
  const [children, setChildren] = useState<SetItem[]>([])
  const [breadcrumbs, setBreadcrumbs] = useState<SetItem[]>([])
  const [view, setView] = useState<'tree' | 'collection'>('tree')
  const [collection, setCollection] = useState<SetItem[]>([])

  useEffect(() => {
    if (view === 'tree') {
      fetch(`/sets/${currentId === 'root' ? 'root' : currentId}/children`)
        .then(r => r.json()).then(setChildren)
      if (currentId !== 'root') {
        fetch(`/sets/${currentId}/breadcrumbs`).then(r => r.json()).then(setBreadcrumbs)
      } else {
        setBreadcrumbs([])
      }
    } else {
      fetch('/collections').then(r => r.json()).then(async (ids: string[]) => {
        const sets = await Promise.all(ids.map(id => fetch(`/sets/${id}`).then(r => r.json())))
        setCollection(sets)
      })
    }
  }, [currentId, view])

  const explore = (id: string) => {
    setCurrentId(id)
    setView('tree')
  }

  const collect = (id: string) => {
    fetch('/collections', {method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({setId:id})})
  }

  const report = (id: string) => {
    const sighting: Sighting = { setId: id }
    fetch('/sightings', {method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify(sighting)})
  }

  return (
    <div className="space-y-4">
      <header className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Livsverket</h1>
        <div className="space-x-2">
          <button className="px-2 py-1 bg-blue-500 text-white" onClick={() => setView('tree')}>Tree</button>
          <button className="px-2 py-1 bg-green-500 text-white" onClick={() => setView('collection')}>My Collection</button>
        </div>
      </header>

      {view === 'tree' && (
        <div>
          {breadcrumbs.length > 0 && (
            <div className="mb-2">
              {breadcrumbs.map((b, idx) => (
                <span key={b.id}>
                  {idx > 0 && ' / '}
                  <button className="underline" onClick={() => explore(b.id)}>{b.nameSv}</button>
                </span>
              ))}
            </div>
          )}
          <ul className="space-y-1">
            {children.map(c => (
              <li key={c.id} className="border p-2 flex justify-between">
                <div>
                  <div className="font-semibold">{c.nameSv}</div>
                  <div className="text-sm text-gray-600">{c.rank}</div>
                </div>
                <div className="space-x-1">
                  <button className="px-1 bg-blue-200" onClick={() => explore(c.id)}>Explore</button>
                  <button className="px-1 bg-green-200" onClick={() => collect(c.id)}>Collect</button>
                  <button className="px-1 bg-yellow-200" onClick={() => report(c.id)}>Report</button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}

      {view === 'collection' && (
        <div>
          <ul className="space-y-1">
            {collection.map(c => (
              <li key={c.id} className="border p-2 flex justify-between">
                <div>
                  <div className="font-semibold">{c.nameSv}</div>
                  <div className="text-sm text-gray-600">{c.rank}</div>
                </div>
                <div>
                  <button className="underline" onClick={() => explore(c.id)}>Go</button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}
