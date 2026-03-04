import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
const icons=['menu','user','settings','search','shield','database','server','terminal','activity','bell']
export default function IconListPage(){return <Card><CardHeader><CardTitle>图标列表示例</CardTitle></CardHeader><CardContent><div className='grid grid-cols-2 gap-2 md:grid-cols-5'>{icons.map(i=><div key={i} className='rounded-md border p-3 text-sm'>{i}</div>)}</div></CardContent></Card>}
