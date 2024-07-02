
export default function getPageTitle(pageTitle:string) {
  const  title = import.meta.env.APP_APP_TITLE || 'NX design'
  if (pageTitle) {
    return `${pageTitle} - ${title}`
  }
  return `${title}`
}
