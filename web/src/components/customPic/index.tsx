import { useMemo } from 'react'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useAuthStore } from '@/stores/useAuthStore'

type CustomPicProps = {
  picType?: 'avatar' | 'img' | 'file'
  picSrc?: string
}

export default function CustomPic({ picType = 'avatar', picSrc = '' }: CustomPicProps) {
  const userInfo = useAuthStore((s) => s.userInfo)
  const base = `${import.meta.env.VITE_BASE_API || ''}/`

  const resolved = useMemo(() => {
    const source = picSrc || userInfo?.headerImg || ''
    if (!source) return ''
    return source.startsWith('http') ? source : `${base}${source}`
  }, [base, picSrc, userInfo?.headerImg])

  if (picType === 'file') {
    return <img src={resolved} alt="file" className="h-20 w-20 object-cover" />
  }

  if (picType === 'img') {
    return <img src={resolved} alt="avatar" className="h-10 w-10 rounded-full object-cover" />
  }

  return (
    <Avatar>
      <AvatarImage src={resolved} />
      <AvatarFallback>{(userInfo?.nickName || 'U').slice(0, 1)}</AvatarFallback>
    </Avatar>
  )
}
