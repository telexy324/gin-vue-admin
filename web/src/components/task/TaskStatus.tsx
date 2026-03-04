import { Badge } from '@/components/ui/badge'

type TaskStatusValue = 'waiting' | 'running' | 'success' | 'error' | 'stopping' | 'stopped' | string

type TaskStatusProps = {
  status?: TaskStatusValue
}

const labelMap: Record<string, string> = {
  waiting: 'Waiting',
  running: 'Running',
  success: 'Success',
  error: 'Failed',
  stopping: 'Stopping...',
  stopped: 'Stopped'
}

const classMap: Record<string, string> = {
  waiting: 'bg-sky-400 hover:bg-sky-400',
  running: 'bg-blue-500 hover:bg-blue-500',
  success: 'bg-emerald-500 hover:bg-emerald-500',
  error: 'bg-rose-500 hover:bg-rose-500',
  stopping: 'bg-orange-500 hover:bg-orange-500',
  stopped: 'bg-amber-500 hover:bg-amber-500'
}

export default function TaskStatus({ status }: TaskStatusProps) {
  if (!status) return null
  const label = labelMap[status] || status
  return <Badge className={classMap[status] || 'bg-slate-500 hover:bg-slate-500'}>{label}</Badge>
}
