import { AlertTriangle } from 'lucide-react'

type WarningBarProps = {
  title?: string
}

export default function WarningBar({ title = '' }: WarningBarProps) {
  return (
    <div className="mb-3 flex items-center rounded-sm bg-orange-50 px-3 py-2 text-sm text-orange-500">
      <AlertTriangle className="mr-2 h-4 w-4" />
      <span>{title}</span>
    </div>
  )
}
