import { useEffect, useMemo, useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
import TaskStatus from '@/components/task/TaskStatus'
import { getTaskById, getTaskOutputs, stopTask } from '@/api/task'
import { getUserById } from '@/api/user'
import socket from '@/socket'

type TaskLogViewProps = {
  itemId?: number
  open?: boolean
  onClose?: () => void
}

export default function TaskLogView({ itemId = 0, open = true, onClose }: TaskLogViewProps) {
  const [visible, setVisible] = useState(open)
  const [item, setItem] = useState<any>({})
  const [output, setOutput] = useState<any[]>([])
  const [userName, setUserName] = useState('')
  const outputRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    setVisible(open)
  }, [open])

  useEffect(() => {
    if (!itemId || !visible) return
    const load = async () => {
      const taskRes = await getTaskById({ ID: itemId })
      const outRes = await getTaskOutputs({ taskId: itemId })
      const nextItem = taskRes?.data?.task || {}
      setItem(nextItem)
      setOutput(outRes?.data?.taskOutputs || [])
      if (nextItem?.systemUserId === 999999) {
        setUserName('定时任务')
      } else if (nextItem?.systemUserId) {
        const userRes = await getUserById({ ID: nextItem.systemUserId })
        setUserName(userRes?.data?.user?.userName || '')
      }
    }
    load()
  }, [itemId, visible])

  useEffect(() => {
    if (!outputRef.current) return
    outputRef.current.scrollTop = outputRef.current.scrollHeight
  }, [output])

  useEffect(() => {
    if (!visible || !itemId) return
    if (!socket.isRunning()) {
      socket.start()
    }
    const listenerId = socket.addListener((data: any) => {
      if (!data || Number(data.taskId) !== Number(itemId)) return
      if (data.type === 'update') {
        setItem((prev: any) => ({ ...prev, ...data }))
      } else if (data.type === 'log') {
        setOutput((prev) => [...prev, data])
      }
    })
    return () => {
      socket.removeListener(listenerId)
    }
  }, [itemId, visible])

  const title = useMemo(() => {
    const id = item?.ID ? `Task #${item.ID}` : 'Task'
    return `${id} ${userName ? `| ${userName}` : ''}`
  }, [item, userName])

  if (!visible) return null

  return (
    <div className="fixed inset-0 z-50 bg-black/40">
      <div className="mx-auto mt-[6vh] w-[min(1200px,94vw)] rounded-lg border bg-white shadow-xl">
        <div className="flex items-center justify-between border-b px-4 py-3">
          <div className="flex items-center gap-3">
            <span className="text-sm font-semibold">{title}</span>
            <TaskStatus status={item?.status} />
          </div>
          <div className="flex items-center gap-2">
            {(item?.status === 'running' || item?.status === 'waiting') && (
              <Button size="sm" onClick={() => stopTask({ ID: item.ID })} className="bg-rose-500 hover:bg-rose-500">
                Stop
              </Button>
            )}
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
        </div>
        <div ref={outputRef} className="max-h-[70vh] overflow-auto bg-black p-3 font-mono text-xs text-white">
          {output.map((record, index) => (
            <div key={record.ID || index} className="flex gap-3 py-0.5">
              <span className="w-40 shrink-0 text-slate-300">
                {record.recordTime ? new Date(record.recordTime).toLocaleString() : ''}
              </span>
              <span className="w-28 shrink-0 text-slate-300">{record.manageIp}</span>
              <span className="whitespace-pre-wrap break-all">{record.output}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
