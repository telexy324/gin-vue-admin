export const toUpperCase = (str: string) => (str?.[0] ? str.replace(str[0], str[0].toUpperCase()) : '')

export const toLowerCase = (str: string) => (str?.[0] ? str.replace(str[0], str[0].toLowerCase()) : '')

export const toSQLLine = (str: string) => {
  if (str === 'ID') return 'ID'
  return str.replace(/([A-Z])/g, '_$1').toLowerCase()
}

export const toHump = (name: string) => name.replace(/\_(\w)/g, (_, letter) => letter.toUpperCase())

