import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'

type ScriptViewProps = {
  script?: string
  open?: boolean
  onClose?: () => void
}

export default function ScriptView({ script = '', open = true, onClose }: ScriptViewProps) {
  const [visible, setVisible] = useState(open)

  useEffect(() => {
    setVisible(open)
  }, [open])

  if (!visible) return null

  return (
    <div className="fixed inset-0 z-50 bg-black/40">
      <div className="mx-auto mt-[10vh] w-[min(960px,92vw)] rounded-lg border bg-white shadow-xl">
        <div className="flex items-center justify-between border-b px-4 py-3">
          <h3 className="text-sm font-semibold">Script</h3>
          <Button
            size="sm"
            variant="outline"
            onClick={() => {
              setVisible(false)
              onClose?.()
            }}
          >
            关闭
          </Button>
        </div>
        <pre className="max-h-[70vh] overflow-auto bg-black p-4 text-xs text-white">{script}</pre>
      </div>
    </div>
  )
}
