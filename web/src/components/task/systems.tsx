import { useEffect, useMemo, useState } from 'react'
import { getAdminSystems } from '@/api/cmdb'
import { Button } from '@/components/ui/button'

type SystemsProps = {
  keys?: number[]
  onChecked?: (rows: any[]) => void
}

export default function Systems({ keys = [], onChecked }: SystemsProps) {
  const [systems, setSystems] = useState<any[]>([])
  const [checkedIds, setCheckedIds] = useState<number[]>(keys)

  useEffect(() => {
    setCheckedIds(keys)
  }, [keys])

  useEffect(() => {
    const load = async () => {
      const res = await getAdminSystems({})
      setSystems(res?.data?.systems || [])
    }
    load()
  }, [])

  const checkedRows = useMemo(() => {
    const map = new Map<number, any>()
    systems.forEach((item) => map.set(Number(item.ID), item))
    return checkedIds.map((id) => map.get(Number(id))).filter(Boolean)
  }, [checkedIds, systems])

  return (
    <div className="space-y-3">
      <div className="flex justify-end">
        <Button size="sm" onClick={() => onChecked?.(checkedRows)}>
          确定
        </Button>
      </div>
      <div className="max-h-[70vh] overflow-auto rounded-md border p-2">
        {systems.map((item) => {
          const id = Number(item.ID)
          return (
            <label key={id} className="flex items-center gap-2 py-1 text-sm">
              <input
                type="checkbox"
                checked={checkedIds.includes(id)}
                onChange={(e) => {
                  setCheckedIds((prev) =>
                    e.target.checked ? Array.from(new Set([...prev, id])) : prev.filter((x) => x !== id)
                  )
                }}
              />
              <span>{item.name}</span>
            </label>
          )
        })}
      </div>
    </div>
  )
}
