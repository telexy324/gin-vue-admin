const dictionaryCache: Record<string, any[]> = {}

export const setDict = (type: string, values: any[]) => {
  dictionaryCache[type] = values || []
}

export const getDict = async (type: string) => {
  return dictionaryCache[type] || []
}

