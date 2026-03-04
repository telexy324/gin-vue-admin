import { useRef } from 'react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/useAuthStore'

type UploadImageProps = {
  imageUrl?: string
  fileSize?: number
  onChange?: (url: string) => void
  onSuccess?: () => void
}

export default function UploadImage({ fileSize = 2048, onChange, onSuccess }: UploadImageProps) {
  const inputRef = useRef<HTMLInputElement | null>(null)
  const token = useAuthStore((s) => s.token)
  const base = import.meta.env.VITE_BASE_API || ''

  const handlePick = async (file?: File) => {
    if (!file) return
    if (!['image/jpeg', 'image/png'].includes(file.type)) return
    if (file.size / 1024 > fileSize) return

    const formData = new FormData()
    formData.append('file', file)

    const res = await fetch(`${base}/fileUploadAndDownload/upload`, {
      method: 'POST',
      headers: { 'x-token': token || '' },
      body: formData
    })
    const json = await res.json()
    const url = json?.data?.file?.url
    if (url) {
      onChange?.(url)
      onSuccess?.()
    }
  }

  return (
    <div>
      <input
        ref={inputRef}
        type="file"
        className="hidden"
        accept="image/jpeg,image/png"
        onChange={(e) => handlePick(e.target.files?.[0])}
      />
      <Button size="sm" onClick={() => inputRef.current?.click()}>
        压缩上传
      </Button>
    </div>
  )
}
