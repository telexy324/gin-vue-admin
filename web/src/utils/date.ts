type DateWithFormat = Date & {
  Format: (fmt: string) => string
}

declare global {
  interface Date {
    Format: (fmt: string) => string
  }
}

Date.prototype.Format = function (fmt: string) {
  const date = this as Date
  const o: Record<string, number> = {
    'M+': date.getMonth() + 1,
    'd+': date.getDate(),
    'h+': date.getHours(),
    'm+': date.getMinutes(),
    's+': date.getSeconds(),
    'q+': Math.floor((date.getMonth() + 3) / 3),
    S: date.getMilliseconds()
  }
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substring(4 - RegExp.$1.length))
  }
  for (const k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(RegExp.$1, RegExp.$1.length === 1 ? String(o[k]) : ('00' + o[k]).substring(String(o[k]).length))
    }
  }
  return fmt
}

export function formatTimeToStr(times: string | number | Date, pattern?: string) {
  const d = new Date(times) as DateWithFormat
  return pattern ? d.Format(pattern) : d.Format('yyyy-MM-dd hh:mm:ss')
}
