import { FormEvent, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { createTemp, getColumn, getDB, getTable, preview } from '@/api/autoCode'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Field = {
  fieldName: string
  fieldDesc: string
  fieldType: string
  dataType?: string
  fieldJson: string
  columnName?: string
  dataTypeLong?: string
  comment?: string
  fieldSearchType?: string
  dictType?: string
}

const emptyField: Field = { fieldName: '', fieldDesc: '', fieldType: 'string', fieldJson: '' }

export default function AutoCodePage() {
  const [dbOptions, setDbOptions] = useState<any[]>([])
  const [tableOptions, setTableOptions] = useState<any[]>([])
  const [dbName, setDbName] = useState('')
  const [tableName, setTableName] = useState('')
  const [previewCode, setPreviewCode] = useState('')

  const [form, setForm] = useState<any>({
    structName: '',
    tableName: '',
    packageName: '',
    abbreviation: '',
    description: '',
    autoCreateApiToSql: false,
    autoMoveFile: false,
    fields: [] as Field[]
  })

  useEffect(() => {
    const load = async () => {
      const res = await getDB()
      if (res?.code === 0) setDbOptions(res.data || [])
    }
    load()
  }, [])

  const loadTables = async (db: string) => {
    setDbName(db)
    setTableName('')
    const res = await getTable({ dbName: db })
    if (res?.code === 0) setTableOptions(res.data || [])
  }

  const useTable = async () => {
    if (!dbName || !tableName) return
    const res = await getColumn({ dbName, tableName })
    if (res?.code === 0) {
      setForm((f: any) => ({
        ...f,
        tableName,
        structName: tableName.replace(/(^|_)(\w)/g, (_: string, __: string, c: string) => c.toUpperCase()),
        packageName: tableName,
        abbreviation: tableName.slice(0, 3),
        description: tableName,
        fields: (res.data || []).map((c: any) => ({
          fieldName: c.fieldName || c.columnName || '',
          fieldDesc: c.comment || c.columnName || '',
          fieldType: c.fieldType || 'string',
          dataType: c.dataType,
          fieldJson: c.fieldJson || c.columnName || '',
          columnName: c.columnName,
          dataTypeLong: c.dataTypeLong,
          comment: c.comment
        }))
      }))
      toast.success('已按数据表加载字段')
    }
  }

  const updateField = (idx: number, patch: Partial<Field>) => {
    setForm((f: any) => {
      const next = [...f.fields]
      next[idx] = { ...next[idx], ...patch }
      return { ...f, fields: next }
    })
  }

  const addField = () => setForm((f: any) => ({ ...f, fields: [...f.fields, { ...emptyField }] }))
  const deleteField = (idx: number) => setForm((f: any) => ({ ...f, fields: f.fields.filter((_: any, i: number) => i !== idx) }))

  const submit = async (e: FormEvent, isPreview: boolean) => {
    e.preventDefault()
    if (!form.structName || !form.abbreviation || !form.description || !form.packageName || !form.fields.length) {
      toast.error('请完善基础信息和字段')
      return
    }
    if (isPreview) {
      const res = await preview(form)
      if (res?.code === 0) {
        setPreviewCode(JSON.stringify(res.data?.autoCode || {}, null, 2))
      }
      return
    }
    const res = await createTemp(form)
    if (res) toast.success('生成成功，已触发下载')
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>自动代码生成</CardTitle>
          <CardDescription>支持从数据库表拉取字段，预览并生成代码。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="grid gap-2 md:grid-cols-4">
            <select className="h-10 rounded-md border bg-background px-3 text-sm" value={dbName} onChange={(e) => loadTables(e.target.value)}>
              <option value="">选择数据库</option>
              {dbOptions.map((d: any) => <option key={d.database} value={d.database}>{d.database}</option>)}
            </select>
            <select className="h-10 rounded-md border bg-background px-3 text-sm" value={tableName} onChange={(e) => setTableName(e.target.value)}>
              <option value="">选择数据表</option>
              {tableOptions.map((t: any) => <option key={t.tableName} value={t.tableName}>{t.tableName}</option>)}
            </select>
            <Button variant="outline" onClick={useTable}>使用此表创建</Button>
          </div>

          <form className="grid gap-2 md:grid-cols-5" onSubmit={(e) => submit(e, false)}>
            <Input placeholder="Struct名称" value={form.structName} onChange={(e) => setForm((f: any) => ({ ...f, structName: e.target.value }))} />
            <Input placeholder="表名" value={form.tableName} onChange={(e) => setForm((f: any) => ({ ...f, tableName: e.target.value }))} />
            <Input placeholder="简称" value={form.abbreviation} onChange={(e) => setForm((f: any) => ({ ...f, abbreviation: e.target.value }))} />
            <Input placeholder="中文描述" value={form.description} onChange={(e) => setForm((f: any) => ({ ...f, description: e.target.value }))} />
            <Input placeholder="文件名" value={form.packageName} onChange={(e) => setForm((f: any) => ({ ...f, packageName: e.target.value }))} />
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!form.autoCreateApiToSql} onChange={(e) => setForm((f: any) => ({ ...f, autoCreateApiToSql: e.target.checked }))} />自动创建API</label>
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!form.autoMoveFile} onChange={(e) => setForm((f: any) => ({ ...f, autoMoveFile: e.target.checked }))} />自动移动文件</label>
            <div className="flex gap-2 md:col-span-3">
              <Button type="button" variant="outline" onClick={(e) => submit(e as any, true)}>预览代码</Button>
              <Button type="submit">生成代码</Button>
            </div>
          </form>

          <div className="flex justify-end"><Button variant="secondary" onClick={addField}>新增Field</Button></div>
          <Table>
            <TableHeader><TableRow><TableHead>Field名</TableHead><TableHead>中文名</TableHead><TableHead>JSON</TableHead><TableHead>类型</TableHead><TableHead>数据库列</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
            <TableBody>
              {form.fields.map((row: Field, idx: number) => (
                <TableRow key={`${row.fieldName}-${idx}`}>
                  <TableCell><Input value={row.fieldName} onChange={(e) => updateField(idx, { fieldName: e.target.value })} /></TableCell>
                  <TableCell><Input value={row.fieldDesc} onChange={(e) => updateField(idx, { fieldDesc: e.target.value })} /></TableCell>
                  <TableCell><Input value={row.fieldJson} onChange={(e) => updateField(idx, { fieldJson: e.target.value })} /></TableCell>
                  <TableCell><Input value={row.fieldType} onChange={(e) => updateField(idx, { fieldType: e.target.value })} /></TableCell>
                  <TableCell><Input value={row.columnName || ''} onChange={(e) => updateField(idx, { columnName: e.target.value })} /></TableCell>
                  <TableCell><Button size="sm" variant="destructive" onClick={() => deleteField(idx)}>删除</Button></TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {!!previewCode && <pre className="max-h-80 overflow-auto rounded-md bg-muted p-3 text-xs">{previewCode}</pre>}
        </CardContent>
      </Card>
    </div>
  )
}
