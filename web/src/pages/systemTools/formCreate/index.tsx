import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { toast } from 'sonner'

type FieldType = 'text' | 'number' | 'textarea' | 'select'
type Field = {
  key: string
  label: string
  type: FieldType
  required: boolean
  placeholder?: string
  options?: string[]
}

const emptyField: Field = {
  key: '',
  label: '',
  type: 'text',
  required: false,
  placeholder: ''
}

export default function FormCreatePage() {
  const [fields, setFields] = useState<Field[]>([])
  const [field, setField] = useState<Field>(emptyField)
  const [schemaText, setSchemaText] = useState('')
  const [formData, setFormData] = useState<Record<string, any>>({})

  const addField = () => {
    if (!field.key || !field.label) {
      toast.error('字段 key 和 label 必填')
      return
    }
    if (fields.some((f) => f.key === field.key)) {
      toast.error('字段 key 不能重复')
      return
    }
    setFields((prev) => [...prev, { ...field }])
    setField(emptyField)
  }

  const deleteField = (idx: number) => {
    const f = fields[idx]
    setFields((prev) => prev.filter((_, i) => i !== idx))
    setFormData((prev) => {
      const next = { ...prev }
      delete next[f.key]
      return next
    })
  }

  const exportSchema = () => {
    const text = JSON.stringify({ fields }, null, 2)
    setSchemaText(text)
    navigator.clipboard?.writeText(text)
    toast.success('Schema 已复制到剪贴板')
  }

  const importSchema = () => {
    try {
      const parsed = JSON.parse(schemaText)
      if (!Array.isArray(parsed.fields)) throw new Error('schema.fields 必须是数组')
      setFields(parsed.fields)
      toast.success('Schema 导入成功')
    } catch (e: any) {
      toast.error(`导入失败: ${e.message || 'JSON 格式错误'}`)
    }
  }

  const renderPreviewField = (f: Field) => {
    if (f.type === 'textarea') {
      return (
        <textarea
          className="min-h-20 w-full rounded-md border bg-background px-3 py-2 text-sm"
          value={formData[f.key] || ''}
          placeholder={f.placeholder}
          onChange={(e) => setFormData((d) => ({ ...d, [f.key]: e.target.value }))}
        />
      )
    }
    if (f.type === 'select') {
      return (
        <select
          className="h-10 w-full rounded-md border bg-background px-3 text-sm"
          value={formData[f.key] || ''}
          onChange={(e) => setFormData((d) => ({ ...d, [f.key]: e.target.value }))}
        >
          <option value="">请选择</option>
          {(f.options || []).map((opt) => (
            <option key={opt} value={opt}>
              {opt}
            </option>
          ))}
        </select>
      )
    }
    return (
      <Input
        type={f.type === 'number' ? 'number' : 'text'}
        value={formData[f.key] || ''}
        placeholder={f.placeholder}
        onChange={(e) => setFormData((d) => ({ ...d, [f.key]: f.type === 'number' ? Number(e.target.value) : e.target.value }))}
      />
    )
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>表单设计器（轻量版）</CardTitle>
          <CardDescription>支持字段配置、Schema 导入导出和实时预览。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="grid gap-2 rounded-md border p-3 md:grid-cols-6">
            <Input placeholder="字段 key" value={field.key} onChange={(e) => setField((f) => ({ ...f, key: e.target.value }))} />
            <Input placeholder="字段 label" value={field.label} onChange={(e) => setField((f) => ({ ...f, label: e.target.value }))} />
            <select className="h-10 rounded-md border bg-background px-3 text-sm" value={field.type} onChange={(e) => setField((f) => ({ ...f, type: e.target.value as FieldType }))}>
              <option value="text">text</option>
              <option value="number">number</option>
              <option value="textarea">textarea</option>
              <option value="select">select</option>
            </select>
            <Input placeholder="placeholder" value={field.placeholder || ''} onChange={(e) => setField((f) => ({ ...f, placeholder: e.target.value }))} />
            <Input
              placeholder="select 选项(逗号分隔)"
              value={(field.options || []).join(',')}
              onChange={(e) => setField((f) => ({ ...f, options: e.target.value.split(',').map((x) => x.trim()).filter(Boolean) }))}
            />
            <div className="flex items-center gap-2">
              <label className="flex items-center gap-1 text-sm">
                <input type="checkbox" checked={field.required} onChange={(e) => setField((f) => ({ ...f, required: e.target.checked }))} />
                必填
              </label>
              <Button type="button" onClick={addField}>
                添加字段
              </Button>
            </div>
          </div>

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>key</TableHead>
                <TableHead>label</TableHead>
                <TableHead>type</TableHead>
                <TableHead>required</TableHead>
                <TableHead>options</TableHead>
                <TableHead>操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {fields.map((f, idx) => (
                <TableRow key={`${f.key}-${idx}`}>
                  <TableCell>{f.key}</TableCell>
                  <TableCell>{f.label}</TableCell>
                  <TableCell>{f.type}</TableCell>
                  <TableCell>{f.required ? '是' : '否'}</TableCell>
                  <TableCell>{(f.options || []).join(', ') || '-'}</TableCell>
                  <TableCell>
                    <Button size="sm" variant="destructive" onClick={() => deleteField(idx)}>
                      删除
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Schema</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex gap-2">
            <Button type="button" variant="outline" onClick={exportSchema}>
              导出 Schema
            </Button>
            <Button type="button" onClick={importSchema}>
              导入 Schema
            </Button>
          </div>
          <textarea
            className="min-h-48 w-full rounded-md border bg-background px-3 py-2 text-xs"
            value={schemaText}
            onChange={(e) => setSchemaText(e.target.value)}
            placeholder="Schema JSON"
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>实时预览</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-2">
          {fields.map((f) => (
            <div key={f.key}>
              <div className="mb-1 text-sm">
                {f.label} {f.required ? '*' : ''}
              </div>
              {renderPreviewField(f)}
            </div>
          ))}
          {!fields.length && <div className="text-sm text-muted-foreground">请先添加字段</div>}
          <pre className="md:col-span-2 max-h-52 overflow-auto rounded-md bg-muted p-3 text-xs">{JSON.stringify(formData, null, 2)}</pre>
        </CardContent>
      </Card>
    </div>
  )
}
