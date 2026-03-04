import { useEffect, useState } from 'react'
import { getSystemState } from '@/api/system'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export default function SystemStatePage(){
  const [state,setState]=useState<any>(null)
  useEffect(()=>{(async()=>{const r=await getSystemState(); if(r?.code===0) setState(r.data||r)})()},[])
  return <Card><CardHeader><CardTitle>系统状态</CardTitle></CardHeader><CardContent><pre className='max-h-[70vh] overflow-auto rounded-md bg-muted p-3 text-xs'>{JSON.stringify(state||{},null,2)}</pre></CardContent></Card>
}
