import type { ReactNode } from 'react'

type NoDataProps = {
  title?: string
  children?: ReactNode
}

export default function NoData({ title = 'No Data', children }: NoDataProps) {
  return (
    <div className="my-24 text-center">
      <div className="text-lg text-slate-500">{title}</div>
      {children}
    </div>
  )
}
