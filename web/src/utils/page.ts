const appName = 'Gin Admin React'

export default function getPageTitle(pageTitle?: string) {
  if (pageTitle) return `${pageTitle} - ${appName}`
  return appName
}

