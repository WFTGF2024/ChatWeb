
import { httpCore } from './http'

export async function createPage({ url, title='', fetch=true }){
  const { data } = await httpCore.post('/web/items', { url, title, fetch })
  return data
}
export async function listPages(){
  const { data } = await httpCore.get('/web/items')
  return data?.items || []
}
export async function getPage(page_id){
  const { data } = await httpCore.get(`/web/page/${page_id}`)
  return data
}
export async function updatePage(page_id, payload){
  const { data } = await httpCore.put(`/web/items/${page_id}`, payload)
  return data
}
export async function deletePage(page_id){
  const { data } = await httpCore.delete(`/web/items/${page_id}`)
  return data
}
export async function webSearch({ q, urls=[], top_k=5 }){
  const body = urls?.length ? { urls, top_k } : { q, top_k }
  const { data } = await httpCore.post('/web/search', body)
  return data
}
export async function ingest(){ const { data } = await httpCore.post('/web/ingest'); return data }
export async function chunk(){ const { data } = await httpCore.post('/web/chunk'); return data }
