import { forwardRef, useImperativeHandle, useMemo, useState } from 'react'
import { getFileList } from '@/api/fileUploadAndDownload'
import { Button } from '@/components/ui/button'

type ChooseImgProps = {
  target?: Record<string, any> | null
  targetKey?: string
  onEnterImg?: (url: string) => void
}

export type ChooseImgRef = {
  open: () => Promise<void>
}

const ChooseImg = forwardRef<ChooseImgRef, ChooseImgProps>(function ChooseImg(
  { target = null, targetKey = '', onEnterImg },
  ref
) {
  const [visible, setVisible] = useState(false)
  const [picList, setPicList] = useState<any[]>([])
  const base = import.meta.env.VITE_BASE_API || ''

  useImperativeHandle(ref, () => ({
    open: async () => {
      const res = await getFileList({ page: 1, pageSize: 9999 })
      setPicList(res?.data?.list || [])
      setVisible(true)
    }
  }))

  const list = useMemo(
    () =>
      picList.map((item) => {
        const src = item?.url || ''
        return {
          ...item,
          fullUrl: src && !String(src).startsWith('http') ? `${base}${src}` : src
        }
      }),
    [base, picList]
  )

  if (!visible) return null

  return (
    <div className="fixed inset-0 z-50 bg-black/30">
      <div className="absolute right-0 top-0 h-full w-[min(700px,96vw)] bg-white shadow-xl">
        <div className="flex items-center justify-between border-b px-4 py-3">
          <h3 className="text-sm font-semibold">媒体库</h3>
          <Button size="sm" variant="outline" onClick={() => setVisible(false)}>
            关闭
          </Button>
        </div>
        <div className="grid max-h-[calc(100vh-56px)] grid-cols-3 gap-3 overflow-auto p-4">
          {list.map((item) => (
            <button
              key={item.ID || item.url}
              className="aspect-square overflow-hidden rounded-md border"
              onClick={() => {
                if (target && targetKey) target[targetKey] = item.url
                onEnterImg?.(item.url)
                setVisible(false)
              }}
            >
              <img src={item.fullUrl} alt="" className="h-full w-full object-cover" />
            </button>
          ))}
        </div>
      </div>
    </div>
  )
})

export default ChooseImg
