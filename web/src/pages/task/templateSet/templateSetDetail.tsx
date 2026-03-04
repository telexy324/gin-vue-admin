import { useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { processSetTask, getSetTaskById, redoSetTask, setTaskForceCorrect } from '@/api/template'
import { getTaskListBySetTaskId } from '@/api/task'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { toast } from 'sonner'

type InnerInput = {
  setTaskInnerSeq: number
  templateId: number
  commandVars: string[]
  targetIds: number[]
}

const calcStepStatus = (stepTasks: any[] = []) => {
  if (!stepTasks.length) return 'wait'
  if (stepTasks.every((t) => t.status === 'success')) return 'success'
  if (stepTasks.every((t) => ['success', 'waiting', 'running'].includes(t.status))) return 'process'
  return 'error'
}

export default function TemplateSetDetailPage() {
  const { setTaskId } = useParams()
  const [setTask, setSetTask] = useState<any>(null)
  const [tasks, setTasks] = useState<any[]>([])
  const [innerInputs, setInnerInputs] = useState<Record<number, InnerInput>>({})

  const load = async () => {
    if (!setTaskId) return
    const [sr, tr] = await Promise.all([
      getSetTaskById({ ID: Number(setTaskId) }),
      getTaskListBySetTaskId({ setTaskId: Number(setTaskId), page: 1, pageSize: 999 })
    ])
    if (sr?.code === 0) {
      const st = sr.data
      setSetTask(st)
      const currentStep = st?.currentStep ?? 0
      const currentTemplates = st?.templates?.[currentStep] || []
      const nextMap: Record<number, InnerInput> = {}
      currentTemplates.forEach((t: any) => {
        const seq = Number(t.seqInner ?? t.innerSeq ?? t.seq ?? 0)
        const varCount = Number(t.commandVarNumbers || 0)
        nextMap[seq] = {
          setTaskInnerSeq: seq,
          templateId: Number(t.ID ?? t.templateId ?? 0),
          commandVars: Array.from({ length: varCount }).map(() => ''),
          targetIds: (t.targetServerIds || []).map((x: any) => Number(x)).filter((x: number) => x > 0)
        }
      })
      setInnerInputs(nextMap)
    }
    if (tr?.code === 0) setTasks(tr.data?.list || [])
  }

  useEffect(() => {
    load()
  }, [setTaskId])

  const currentOuterSeq = useMemo(() => {
    const seq = setTask?.currentStep
    return typeof seq === 'number' ? seq : 0
  }, [setTask])

  const steps = useMemo(() => {
    const arr = setTask?.templates || []
    return arr.map((stepTemplates: any[], idx: number) => {
      const stepTasks = setTask?.tasks?.[idx] || []
      return {
        idx,
        templates: stepTemplates || [],
        status: calcStepStatus(stepTasks),
        taskCount: stepTasks.length
      }
    })
  }, [setTask])

  const currentTemplates = useMemo(() => {
    if (!setTask?.templates?.length) return []
    return setTask.templates[currentOuterSeq] || []
  }, [setTask, currentOuterSeq])

  const updateVars = (seq: number, idx: number, value: string) => {
    setInnerInputs((prev) => {
      const item = prev[seq]
      if (!item) return prev
      const commandVars = [...item.commandVars]
      commandVars[idx] = value
      return { ...prev, [seq]: { ...item, commandVars } }
    })
  }

  const updateTargets = (seq: number, value: string) => {
    setInnerInputs((prev) => {
      const item = prev[seq]
      if (!item) return prev
      return {
        ...prev,
        [seq]: {
          ...item,
          targetIds: value
            .split(',')
            .map((x) => Number(x.trim()))
            .filter((x) => Number.isFinite(x) && x > 0)
        }
      }
    })
  }

  const executeSelected = async () => {
    if (!setTask) return
    const setTasks = Object.values(innerInputs)
      .filter((x) => x.setTaskInnerSeq > 0)
      .map((x) => ({
        setTaskInnerSeq: x.setTaskInnerSeq,
        commandVars: x.commandVars,
        targetIds: x.targetIds
      }))

    if (!setTasks.length) {
      toast.error('当前步骤无可执行模板')
      return
    }
    try {
      const res = await processSetTask({ ID: setTask.ID, setTasks })
      if (res?.code === 0) {
        toast.success('已按模板逐项触发执行')
        load()
      } else {
        throw new Error('批量参数执行失败')
      }
    } catch (_) {
      const first = setTasks[0]
      const res = await processSetTask({
        ID: setTask.ID,
        commandVars: first?.commandVars || [],
        targetIds: first?.targetIds || []
      })
      if (res?.code === 0) {
        toast.success('已按兼容模式触发执行')
        load()
      }
    }
  }

  const forceCorrect = async () => {
    if (!setTask) return
    const res = await setTaskForceCorrect({ ID: setTask.ID })
    if (res?.code === 0) {
      toast.success('已强制纠正')
      load()
    }
  }

  const redo = async () => {
    if (!setTask) return
    const res = await redoSetTask({ ID: setTask.ID })
    if (res?.code === 0) {
      toast.success('已重做')
      load()
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>模板集执行详情</CardTitle>
        <CardDescription>
          SetTask ID: {setTaskId || '-'}，当前步骤：{currentOuterSeq} / {(steps.length || 1) - 1}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 md:grid-cols-4">
          <Badge variant="secondary">总步骤 {steps.length}</Badge>
          <Badge variant="secondary">总任务 {tasks.length}</Badge>
          <Badge variant="secondary">可重做 {setTask?.needRedo === 1 ? '是' : '否'}</Badge>
          <Badge variant="secondary">强制纠正 {setTask?.forceCorrect === 1 ? '已开启' : '未开启'}</Badge>
        </div>

        <div className="grid gap-2 rounded-md border p-3 md:grid-cols-2">
          {steps.map((step) => (
            <div key={step.idx} className="rounded-md border bg-muted/30 p-2 text-sm">
              <div className="flex items-center justify-between">
                <span>步骤 {step.idx}</span>
                <Badge variant={step.status === 'success' ? 'default' : step.status === 'error' ? 'outline' : 'secondary'}>
                  {step.status}
                </Badge>
              </div>
              <div className="mt-1 text-xs text-muted-foreground">模板数 {step.templates.length}，任务数 {step.taskCount}</div>
            </div>
          ))}
        </div>

        <div className="space-y-3 rounded-md border p-3">
          <div className="text-sm font-medium">当前步骤模板参数面板</div>
          {currentTemplates.length ? (
            currentTemplates.map((t: any, i: number) => {
              const seq = Number(t.seqInner ?? t.innerSeq ?? t.seq ?? i + 1)
              const input = innerInputs[seq]
              return (
                <div key={`${t.ID || i}-${seq}`} className="rounded-md border p-3">
                  <div className="mb-2 text-sm">
                    seqInner={seq} | templateId={t.ID ?? t.templateId ?? '-'} | vars={t.commandVarNumbers ?? 0}
                  </div>
                  <div className="grid gap-2 md:grid-cols-2">
                    {(input?.commandVars || []).map((v, idx) => (
                      <Input
                        key={`${seq}-${idx}`}
                        placeholder={`参数 ${idx}`}
                        value={v}
                        onChange={(e) => updateVars(seq, idx, e.target.value)}
                      />
                    ))}
                    {!input?.commandVars?.length && <div className="text-xs text-muted-foreground">无需参数</div>}
                    <Input
                      className="md:col-span-2"
                      placeholder="目标ID（逗号分隔）"
                      value={(input?.targetIds || []).join(',')}
                      onChange={(e) => updateTargets(seq, e.target.value)}
                    />
                  </div>
                </div>
              )
            })
          ) : (
            <div className="text-sm text-muted-foreground">无可执行模板</div>
          )}
          <div className="flex gap-2">
            <Button onClick={executeSelected}>按模板执行当前步骤</Button>
            <Button variant="secondary" onClick={redo}>
              重做
            </Button>
            <Button variant="destructive" onClick={forceCorrect}>
              强制执行
            </Button>
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>模板ID</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>开始</TableHead>
              <TableHead>结束</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {tasks.map((row) => (
              <TableRow key={row.ID}>
                <TableCell>#{row.ID}</TableCell>
                <TableCell>{row.templateId}</TableCell>
                <TableCell>{row.status}</TableCell>
                <TableCell>{row.beginTime?.Time || '-'}</TableCell>
                <TableCell>{row.endTime?.Time || '-'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <pre className="max-h-72 overflow-auto rounded-md bg-muted p-3 text-xs">{JSON.stringify(setTask || {}, null, 2)}</pre>
      </CardContent>
    </Card>
  )
}
