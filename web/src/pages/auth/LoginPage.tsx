import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Eye, EyeOff, RefreshCw } from 'lucide-react'
import { toast } from 'sonner'
import { login, captcha } from '@/api/user'
import { checkDB } from '@/api/initdb'
import { useAuthStore } from '@/stores/useAuthStore'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import bg from '@/assets/login_background_dark.jpg'

export default function LoginPage() {
  const navigate = useNavigate()
  const token = useAuthStore((s) => s.token)
  const setToken = useAuthStore((s) => s.setToken)
  const setUserInfo = useAuthStore((s) => s.setUserInfo)

  const [showPwd, setShowPwd] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [captchaImg, setCaptchaImg] = useState('')
  const [form, setForm] = useState({
    username: '',
    password: '',
    captcha: '',
    captchaId: ''
  })

  useEffect(() => {
    if (token) {
      navigate('/layout', { replace: true })
    }
  }, [token, navigate])

  const fetchCaptcha = async () => {
    const res = await captcha({})
    if (res?.code === 0) {
      setCaptchaImg(res.data.picPath)
      setForm((prev) => ({ ...prev, captchaId: res.data.captchaId, captcha: '' }))
    }
  }

  useEffect(() => {
    fetchCaptcha()
  }, [])

  const onChange = (key, value) => setForm((prev) => ({ ...prev, [key]: value }))

  const goInit = async () => {
    const res = await checkDB()
    if (res?.code === 0 && res?.data?.needInit) {
      navigate('/init')
      return
    }
    toast.info('数据库已初始化，无需重复初始化')
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    if (form.username.length < 5 || form.password.length < 6 || form.captcha.length < 5) {
      toast.error('请填写正确的登录信息')
      fetchCaptcha()
      return
    }

    setSubmitting(true)
    try {
      const res = await login(form)
      if (res?.code === 0) {
        setToken(res.data.token)
        setUserInfo(res.data.user)
        navigate('/layout', { replace: true })
      } else {
        fetchCaptcha()
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-cover bg-center p-4" style={{ backgroundImage: `url(${bg})` }}>
      <Card className="w-full max-w-md border-white/30 bg-white/90 shadow-xl backdrop-blur">
        <CardHeader>
          <CardTitle className="text-2xl">Gin Admin</CardTitle>
          <CardDescription>React + shadcn/ui 登录入口</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={onSubmit}>
            <Input placeholder="用户名" value={form.username} onChange={(e) => onChange('username', e.target.value)} />
            <div className="relative">
              <Input
                placeholder="密码"
                type={showPwd ? 'text' : 'password'}
                value={form.password}
                onChange={(e) => onChange('password', e.target.value)}
              />
              <button
                type="button"
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground"
                onClick={() => setShowPwd((v) => !v)}
              >
                {showPwd ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>

            <div className="flex items-center gap-2">
              <Input placeholder="验证码" value={form.captcha} onChange={(e) => onChange('captcha', e.target.value)} />
              <button type="button" onClick={fetchCaptcha} className="h-10 overflow-hidden rounded-md border bg-white">
                {captchaImg ? <img src={captchaImg} alt="captcha" className="h-full w-[120px] object-cover" /> : <RefreshCw className="mx-10 h-4 w-4" />}
              </button>
            </div>

            <div className="flex gap-3">
              <Button type="button" variant="outline" className="flex-1" onClick={goInit}>
                前往初始化
              </Button>
              <Button type="submit" className="flex-1" disabled={submitting}>
                {submitting ? '登录中...' : '登录'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
